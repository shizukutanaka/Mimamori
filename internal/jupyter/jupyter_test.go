package jupyter

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

const validStatusBody = `{"connections": 0, "kernels": 1, "last_activity": "2026-10-18T10:00:00.000000Z", "started": "2026-10-18T09:00:00.000000Z"}`

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "real status response", status: http.StatusOK, body: validStatusBody, want: true},
		{name: "forbidden status", status: http.StatusForbidden, body: validStatusBody},
		{name: "missing kernels", status: http.StatusOK, body: `{"connections":0,"last_activity":"2026-10-18T10:00:00Z","started":"2026-10-18T09:00:00Z"}`},
		{name: "kernels is a string", status: http.StatusOK, body: `{"connections":0,"kernels":"1","last_activity":"2026-10-18T10:00:00Z","started":"2026-10-18T09:00:00Z"}`},
		{name: "missing connections", status: http.StatusOK, body: `{"kernels":1,"last_activity":"2026-10-18T10:00:00Z","started":"2026-10-18T09:00:00Z"}`},
		{name: "connections is a string", status: http.StatusOK, body: `{"connections":"0","kernels":1,"last_activity":"2026-10-18T10:00:00Z","started":"2026-10-18T09:00:00Z"}`},
		{name: "non-JSON", status: http.StatusOK, body: "not JSON"},
		{name: "public version response", status: http.StatusOK, body: `{"version":"2.14.0"}`},
		{name: "empty timestamps", status: http.StatusOK, body: `{"connections":0,"kernels":1,"last_activity":"","started":""}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseStatus(test.status, []byte(test.body)); got != test.want {
				t.Fatalf("ParseStatus() = %t, want %t", got, test.want)
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
		{name: "auth disabled", status: http.StatusOK, body: validStatusBody, want: true},
		{name: "forbidden", status: http.StatusForbidden, body: validStatusBody},
		{name: "redirect not followed", status: http.StatusFound, body: validStatusBody, location: "/redirect"},
		{name: "oversized body", status: http.StatusOK, body: validStatusBody, oversize: true},
		{name: "public version endpoint", status: http.StatusOK, body: `{"version":"2.14.0"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidRequest atomic.Bool
			var forbiddenCredential atomic.Bool
			var redirectRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.URL.Path == "/redirect" {
					redirectRequests.Add(1)
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, validStatusBody)
					return
				}
				if request.Method != http.MethodGet || request.URL.Path != "/api/status" {
					invalidRequest.Store(true)
				}
				for _, name := range []string{"Authorization", "Cookie", "X-XSRFToken"} {
					if request.Header.Get(name) != "" {
						forbiddenCredential.Store(true)
					}
				}
				if request.URL.Query().Has("token") {
					forbiddenCredential.Store(true)
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
				t.Error("request was not GET /api/status")
			}
			if forbiddenCredential.Load() {
				t.Error("request included an Authorization, Cookie, X-XSRFToken header, or token query")
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
		{name: "nil context", ip: ipv4, port: 8888, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 8888, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 8888, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 8888},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 8888, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
