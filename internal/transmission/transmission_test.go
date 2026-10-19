package transmission

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

const (
	testSessionID     = "5m7fT0v6BXzoA6ZK"
	validConflictBody = `<html><body><h1>409: Conflict</h1><p>Your request has been rejected for reasons listed below:</p><ul><li><code>X-Transmission-Session-Id: 5m7fT0v6BXzoA6ZK</code></li></ul></body></html>`
)

func TestParseConflict(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		headerID string
		body     string
		want     bool
	}{
		{name: "real 4.x conflict response", status: http.StatusConflict, headerID: testSessionID, body: validConflictBody, want: true},
		{name: "success status", status: http.StatusOK, headerID: testSessionID, body: validConflictBody},
		{name: "unauthorized status", status: http.StatusUnauthorized, headerID: testSessionID, body: validConflictBody},
		{name: "empty header", status: http.StatusConflict, body: validConflictBody},
		{name: "whitespace header", status: http.StatusConflict, headerID: " \t\n", body: validConflictBody},
		{name: "header and body mismatch", status: http.StatusConflict, headerID: "another-id", body: validConflictBody},
		{name: "body is missing session header line", status: http.StatusConflict, headerID: testSessionID, body: "<html><body>Conflict</body></html>"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseConflict(test.status, test.headerID, []byte(test.body)); got != test.want {
				t.Fatalf("ParseConflict() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		headerID string
		location string
		oversize bool
		want     bool
	}{
		{name: "unauthenticated RPC", status: http.StatusConflict, body: validConflictBody, headerID: testSessionID, want: true},
		{name: "unauthorized", status: http.StatusUnauthorized, body: validConflictBody, headerID: testSessionID},
		{name: "forbidden", status: http.StatusForbidden, body: validConflictBody, headerID: testSessionID},
		{name: "redirect not followed", status: http.StatusMovedPermanently, body: validConflictBody, headerID: testSessionID, location: "/redirect"},
		{name: "oversized body", status: http.StatusConflict, body: validConflictBody, headerID: testSessionID, oversize: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidRequest atomic.Bool
			var forbiddenHeader atomic.Bool
			var redirectRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.URL.Path == "/redirect" {
					redirectRequests.Add(1)
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, validConflictBody)
					return
				}
				if request.Method != http.MethodGet || request.URL.Path != "/transmission/rpc" {
					invalidRequest.Store(true)
				}
				for _, name := range []string{"X-Transmission-Session-Id", "Authorization", "Origin"} {
					if request.Header.Get(name) != "" {
						forbiddenHeader.Store(true)
					}
				}
				if test.status == http.StatusUnauthorized {
					w.Header().Set("WWW-Authenticate", `Basic realm="Transmission"`)
				}
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				if test.headerID != "" {
					w.Header().Set("X-Transmission-Session-Id", test.headerID)
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
				t.Error("request was not GET /transmission/rpc")
			}
			if forbiddenHeader.Load() {
				t.Error("request included a session ID, Authorization, or Origin header")
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
		{name: "nil context", ip: ipv4, port: 9091, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 9091, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 9091, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 9091},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 9091, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
