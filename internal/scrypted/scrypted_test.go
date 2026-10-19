package scrypted

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const validLogin = `{"error":"Not logged in.","hasLogin":false,"externalAddresses":[],"addresses":["https://192.168.1.5:10443"],"hostname":"scrypted"}`

func TestParseLogin(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "real login shape", body: validLogin, want: true},
		{name: "empty addresses", body: `{"hasLogin":false,"addresses":[],"hostname":"scrypted"}`, want: true},
		{name: "account exists", body: `{"hasLogin":true,"addresses":[],"hostname":"scrypted"}`},
		{name: "hasLogin missing", body: `{"addresses":[],"hostname":"scrypted"}`},
		{name: "hasLogin null", body: `{"hasLogin":null,"addresses":[],"hostname":"scrypted"}`},
		{name: "hasLogin is a string", body: `{"hasLogin":"false","addresses":[],"hostname":"scrypted"}`},
		{name: "hasLogin is a number", body: `{"hasLogin":0,"addresses":[],"hostname":"scrypted"}`},
		{name: "addresses missing", body: `{"hasLogin":false,"hostname":"scrypted"}`},
		{name: "addresses is null", body: `{"hasLogin":false,"addresses":null,"hostname":"scrypted"}`},
		{name: "addresses is an object", body: `{"hasLogin":false,"addresses":{},"hostname":"scrypted"}`},
		{name: "addresses is a string", body: `{"hasLogin":false,"addresses":"[]","hostname":"scrypted"}`},
		{name: "hostname missing", body: `{"hasLogin":false,"addresses":[]}`},
		{name: "hostname is null", body: `{"hasLogin":false,"addresses":[],"hostname":null}`},
		{name: "hostname is a number", body: `{"hasLogin":false,"addresses":[],"hostname":1}`},
		{name: "malformed JSON", body: `{"hasLogin":false`},
		{name: "top-level array", body: `[{"hasLogin":false,"addresses":[],"hostname":"scrypted"}]`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseLogin([]byte(test.body)); got != test.want {
				t.Fatalf("ParseLogin() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnclaimed(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		body         string
		location     string
		oversizeBody bool
		want         bool
	}{
		{name: "unclaimed", status: http.StatusOK, body: validLogin, want: true},
		{name: "account exists", status: http.StatusOK, body: `{"hasLogin":true,"addresses":[],"hostname":"scrypted"}`},
		{name: "unauthorized", status: http.StatusUnauthorized, body: `{"error":"Unauthorized"}`},
		{name: "not found", status: http.StatusNotFound, body: `{"error":"Not Found"}`},
		{name: "redirect is not followed", status: http.StatusFound, location: "/other"},
		{name: "oversized login response", status: http.StatusOK, body: validLogin, oversizeBody: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidRequest atomic.Bool
			var redirectFollowed atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.Method != http.MethodGet || request.URL.Path != "/login" {
					invalidRequest.Store(true)
				}
				if request.URL.Path != "/login" {
					redirectFollowed.Store(true)
				}
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				w.WriteHeader(test.status)
				body := test.body
				if test.oversizeBody {
					body += strings.Repeat(" ", maxResponseBytes+1-len(body))
				}
				_, _ = io.WriteString(w, body)
			}))
			defer server.Close()
			port := server.Listener.Addr().(*net.TCPAddr).Port
			got := Unclaimed(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second)
			if got != test.want {
				t.Fatalf("Unclaimed() = %t, want %t", got, test.want)
			}
			if count := requests.Load(); count != 1 {
				t.Errorf("received %d requests, want exactly one", count)
			}
			if invalidRequest.Load() {
				t.Error("request was not GET /login")
			}
			if redirectFollowed.Load() {
				t.Error("followed redirect")
			}
		})
	}
}

func TestUnclaimedGuards(t *testing.T) {
	ipv4 := netip.MustParseAddr("127.0.0.1")
	tests := []struct {
		name    string
		ctx     context.Context
		ip      netip.Addr
		port    int
		timeout time.Duration
	}{
		{name: "nil context", ip: ipv4, port: 11080, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 11080, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 11080, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 11080},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 11080, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unclaimed(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unclaimed() = true, want false")
			}
		})
	}
}
