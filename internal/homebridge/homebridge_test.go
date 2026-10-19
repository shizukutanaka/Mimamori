package homebridge

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

const validSettings = `{"env":{"packageName":"homebridge-config-ui-x","version":"1.0.0"},"formAuth":false}`

func TestParseSettings(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "latest settings shape", body: validSettings, want: true},
		{name: "form auth enabled", body: `{"env":{"packageName":"homebridge-config-ui-x"},"formAuth":true}`},
		{name: "form auth missing", body: `{"env":{"packageName":"homebridge-config-ui-x"}}`},
		{name: "form auth is a string", body: `{"env":{"packageName":"homebridge-config-ui-x"},"formAuth":"false"}`},
		{name: "form auth is null", body: `{"env":{"packageName":"homebridge-config-ui-x"},"formAuth":null}`},
		{name: "form auth is a number", body: `{"env":{"packageName":"homebridge-config-ui-x"},"formAuth":0}`},
		{name: "wrong package name", body: `{"env":{"packageName":"other"},"formAuth":false}`},
		{name: "missing package name", body: `{"env":{},"formAuth":false}`},
		{name: "missing environment", body: `{"formAuth":false}`},
		{name: "environment is an array", body: `{"env":[],"formAuth":false}`},
		{name: "environment is null", body: `{"env":null,"formAuth":false}`},
		{name: "malformed JSON", body: `{"env":`},
		{name: "array", body: `[{"env":{"packageName":"homebridge-config-ui-x"},"formAuth":false}]`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseSettings([]byte(test.body)); got != test.want {
				t.Fatalf("ParseSettings() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		body         string
		location     string
		oversizeBody bool
		want         bool
	}{
		{name: "no authentication", status: http.StatusOK, body: validSettings, want: true},
		{name: "form authentication enabled", status: http.StatusOK, body: `{"env":{"packageName":"homebridge-config-ui-x"},"formAuth":true}`},
		{name: "unauthorized", status: http.StatusUnauthorized, body: `{"error":"Unauthorized"}`},
		{name: "not found", status: http.StatusNotFound, body: `{"error":"Not Found"}`},
		{name: "redirect is not followed", status: http.StatusFound, location: "/redirect-target"},
		{name: "oversized settings", status: http.StatusOK, body: validSettings, oversizeBody: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidRequest atomic.Bool
			var redirectFollowed atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.Method != http.MethodGet || request.URL.Path != "/api/auth/settings" {
					invalidRequest.Store(true)
				}
				if request.URL.Path != "/api/auth/settings" {
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
			got := Unauthenticated(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second)
			if got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			if count := requests.Load(); count != 1 {
				t.Errorf("received %d requests, want exactly one", count)
			}
			if invalidRequest.Load() {
				t.Error("request was not GET /api/auth/settings")
			}
			if redirectFollowed.Load() {
				t.Error("followed redirect")
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
		{name: "nil context", ip: ipv4, port: 8581, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 8581, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 8581, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 8581},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 8581, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
