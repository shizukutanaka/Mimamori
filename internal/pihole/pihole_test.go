package pihole

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const validAuthResponse = `{"session":{"valid":true,"totp":false,"sid":null,"validity":-1,"message":"no password set"},"took":0.001}`

func TestParseAuth(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "no password response", body: validAuthResponse, want: true},
		{name: "validity uses equivalent JSON number", body: `{"session":{"valid":true,"sid":null,"validity":-10e-1,"message":"no password set"}}`, want: true},
		{name: "top-level array", body: `[]`},
		{name: "malformed JSON", body: `{"session":`},
		{name: "session missing", body: `{"took":0.001}`},
		{name: "session is null", body: `{"session":null}`},
		{name: "session is not an object", body: `{"session":[]}`},
		{name: "valid missing", body: `{"session":{"sid":null,"validity":-1,"message":"no password set"}}`},
		{name: "valid false", body: `{"session":{"valid":false,"sid":null,"validity":-1,"message":"no password set"}}`},
		{name: "valid null", body: `{"session":{"valid":null,"sid":null,"validity":-1,"message":"no password set"}}`},
		{name: "valid string", body: `{"session":{"valid":"true","sid":null,"validity":-1,"message":"no password set"}}`},
		{name: "sid missing", body: `{"session":{"valid":true,"validity":-1,"message":"no password set"}}`},
		{name: "sid is string", body: `{"session":{"valid":true,"sid":"session-id","validity":-1,"message":"no password set"}}`},
		{name: "validity missing", body: `{"session":{"valid":true,"sid":null,"message":"no password set"}}`},
		{name: "validity is string", body: `{"session":{"valid":true,"sid":null,"validity":"-1","message":"no password set"}}`},
		{name: "validity is boolean", body: `{"session":{"valid":true,"sid":null,"validity":true,"message":"no password set"}}`},
		{name: "validity differs from negative one", body: `{"session":{"valid":true,"sid":null,"validity":-1.1,"message":"no password set"}}`},
		{name: "message missing", body: `{"session":{"valid":true,"sid":null,"validity":-1}}`},
		{name: "message is null", body: `{"session":{"valid":true,"sid":null,"validity":-1,"message":null}}`},
		{name: "message differs", body: `{"session":{"valid":true,"sid":null,"validity":-1,"message":"logged in"}}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseAuth([]byte(test.body)); got != test.want {
				t.Fatalf("ParseAuth() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		location string
		oversize bool
		want     bool
	}{
		{name: "no password set", status: http.StatusOK, body: validAuthResponse, want: true},
		{name: "password set", status: http.StatusOK, body: `{"session":{"valid":false,"totp":false,"sid":null,"validity":0,"message":"Invalid"}}`},
		{name: "unauthorized", status: http.StatusUnauthorized, body: validAuthResponse},
		{name: "redirect rejected", status: http.StatusFound, body: validAuthResponse, location: "/redirect"},
		{name: "oversized response", status: http.StatusOK, body: validAuthResponse, oversize: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidRequest atomic.Bool
			var hasOriginOrReferer atomic.Bool
			var redirectRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.URL.Path == "/redirect" {
					redirectRequests.Add(1)
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, validAuthResponse)
					return
				}
				if request.Method != http.MethodGet || request.URL.Path != "/api/auth" {
					invalidRequest.Store(true)
				}
				if request.Header.Get("Origin") != "" || request.Header.Get("Referer") != "" {
					hasOriginOrReferer.Store(true)
				}
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				w.WriteHeader(test.status)
				body := test.body
				if test.oversize {
					body += strings.Repeat(" ", maxResponseBytes+1-len(body))
				}
				_, _ = io.WriteString(w, body)
			}))
			defer server.Close()
			address, err := netip.ParseAddrPort(server.Listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			got := Unauthenticated(context.Background(), address.Addr(), int(address.Port()), time.Second)
			if got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			if count := requests.Load(); count != 1 {
				t.Errorf("received %d requests, want exactly one", count)
			}
			if invalidRequest.Load() {
				t.Error("request was not GET /api/auth")
			}
			if hasOriginOrReferer.Load() {
				t.Error("request unexpectedly included Origin or Referer")
			}
			if count := redirectRequests.Load(); count != 0 {
				t.Errorf("followed redirect %d times", count)
			}
		})
	}
}

func TestUnauthenticatedGuards(t *testing.T) {
	ipv4 := netip.MustParseAddr("127.0.0.1")
	tests := []struct {
		name    string
		ctx     context.Context
		ip      netip.Addr
		port    int
		timeout time.Duration
	}{
		{name: "nil context", ip: ipv4, port: 80, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 80, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 80, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 80},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 80, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
