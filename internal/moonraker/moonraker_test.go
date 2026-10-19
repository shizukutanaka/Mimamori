package moonraker

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
	"time"
)

const validInfo = `{"result":{"moonraker_version":"v0.9.3-12-gabcdef","api_version":[1,5,0],"api_version_string":"1.5.0"}}`

func TestParseInfo(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "valid info", body: validInfo, want: true},
		{name: "extra result fields", body: `{"result":{"moonraker_version":"v0.8.0","api_version_string":"1.4.0","api_version":[1,4,0]}}`, want: true},
		{name: "missing result", body: `{"moonraker_version":"v0.9.0","api_version_string":"1.5.0"}`},
		{name: "missing moonraker version", body: `{"result":{"api_version_string":"1.5.0"}}`},
		{name: "empty moonraker version", body: `{"result":{"moonraker_version":"","api_version_string":"1.5.0"}}`},
		{name: "non-string moonraker version", body: `{"result":{"moonraker_version":1,"api_version_string":"1.5.0"}}`},
		{name: "missing api version string", body: `{"result":{"moonraker_version":"v0.9.0"}}`},
		{name: "empty api version string", body: `{"result":{"moonraker_version":"v0.9.0","api_version_string":""}}`},
		{name: "non-string api version string", body: `{"result":{"moonraker_version":"v0.9.0","api_version_string":1}}`},
		{name: "result is not an object", body: `{"result":[]}`},
		{name: "top-level array", body: `[{"result":{"moonraker_version":"v0.9.0","api_version_string":"1.5.0"}}]`},
		{name: "null", body: `null`},
		{name: "non-JSON", body: `<html>not Moonraker</html>`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseInfo([]byte(test.body)); got != test.want {
				t.Fatalf("ParseInfo() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		body      string
		location  string
		oversized bool
		want      bool
	}{
		{name: "trusted info", status: http.StatusOK, body: validInfo, want: true},
		{name: "authentication required", status: http.StatusUnauthorized, body: `{"error":{"code":401,"message":"Unauthorized"}}`},
		{name: "unrelated JSON", status: http.StatusOK, body: `{"result":{"version":"1.0"}}`},
		{name: "missing api version string", status: http.StatusOK, body: `{"result":{"moonraker_version":"v0.9.0"}}`},
		{name: "non-string moonraker version", status: http.StatusOK, body: `{"result":{"moonraker_version":1,"api_version_string":"1.5.0"}}`},
		{name: "redirect", status: http.StatusFound, location: "/redirect", body: validInfo},
		{name: "oversized body", status: http.StatusOK, body: validInfo, oversized: true},
		{name: "HTML", status: http.StatusOK, body: `<html>not Moonraker</html>`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests, redirectRequests int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/redirect" {
					redirectRequests++
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, validInfo)
					return
				}
				requests++
				if request.Method != http.MethodGet || request.URL.Path != "/server/info" {
					t.Errorf("request = %s %s, want GET /server/info", request.Method, request.URL.Path)
				}
				if got := request.Header.Get("Accept"); got != "application/json" {
					t.Errorf("Accept = %q, want application/json", got)
				}
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				w.WriteHeader(test.status)
				body := test.body
				if test.oversized {
					body += strings.Repeat(" ", maxResponseBytes-len(body)+1)
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
			if requests != 1 {
				t.Errorf("request count = %d, want exactly one", requests)
			}
			if redirectRequests != 0 {
				t.Errorf("followed redirect %d times", redirectRequests)
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
		{name: "nil context", ctx: nil, ip: ipv4, port: 7125, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 7125, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), ip: netip.Addr{}, port: 7125, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, port: 0, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 7125},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 7125, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
