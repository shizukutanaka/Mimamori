package codeserver

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

func TestParseHealth(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "alive", status: http.StatusOK, body: `{"status":"alive","lastHeartbeat":123}`, want: true},
		{name: "expired", status: http.StatusOK, body: `{"status":"expired","lastHeartbeat":123.5}`, want: true},
		{name: "heartbeat is a string", status: http.StatusOK, body: `{"status":"alive","lastHeartbeat":"123"}`},
		{name: "missing status", status: http.StatusOK, body: `{"lastHeartbeat":123}`},
		{name: "unknown status", status: http.StatusOK, body: `{"status":"ok","lastHeartbeat":123}`},
		{name: "not JSON", status: http.StatusOK, body: `alive`},
		{name: "forbidden", status: http.StatusForbidden, body: `{"status":"alive","lastHeartbeat":123}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseHealth(test.status, []byte(test.body)); got != test.want {
				t.Fatalf("ParseHealth() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestLoginRedirectsHome(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		location string
		want     bool
	}{
		{name: "dot slash", status: http.StatusFound, location: "./", want: true},
		{name: "dot", status: http.StatusFound, location: ".", want: true},
		{name: "root", status: http.StatusFound, location: "/", want: true},
		{name: "query", status: http.StatusFound, location: "./?x=1", want: true},
		{name: "empty location", status: http.StatusFound},
		{name: "login path", status: http.StatusFound, location: "./login"},
		{name: "external host", status: http.StatusFound, location: "http://other/"},
		{name: "ok", status: http.StatusOK, location: "/"},
		{name: "moved permanently", status: http.StatusMovedPermanently, location: "/"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := LoginRedirectsHome(test.status, test.location); got != test.want {
				t.Fatalf("LoginRedirectsHome(%d, %q) = %t, want %t", test.status, test.location, got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name           string
		healthStatus   int
		healthBody     string
		loginStatus    int
		loginLocation  string
		loginBody      string
		oversizeHealth bool
		oversizeLogin  bool
		want           bool
		wantRequests   int32
	}{
		{name: "password login page", healthStatus: http.StatusOK, healthBody: `{"status":"alive","lastHeartbeat":123}`, loginStatus: http.StatusOK, loginBody: "login", wantRequests: 2},
		{name: "not code-server", healthStatus: http.StatusOK, healthBody: `{"status":"ok","lastHeartbeat":123}`, wantRequests: 1},
		{name: "no-auth login redirects home", healthStatus: http.StatusOK, healthBody: `{"status":"alive","lastHeartbeat":123}`, loginStatus: http.StatusFound, loginLocation: "./", loginBody: "redirect", want: true, wantRequests: 2},
		{name: "oversized health", healthStatus: http.StatusOK, healthBody: `{"status":"alive","lastHeartbeat":123}`, oversizeHealth: true, wantRequests: 1},
		{name: "oversized login", healthStatus: http.StatusOK, healthBody: `{"status":"expired","lastHeartbeat":123}`, loginStatus: http.StatusFound, loginLocation: ".", oversizeLogin: true, wantRequests: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidMethod atomic.Bool
			var authHeaders atomic.Bool
			var unexpectedPath atomic.Bool
			paths := make(chan string, 4)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				paths <- request.Method + " " + request.URL.Path
				if request.Method != http.MethodGet {
					invalidMethod.Store(true)
				}
				if request.Header.Get("Authorization") != "" || request.Header.Get("Cookie") != "" {
					authHeaders.Store(true)
				}
				var body string
				switch request.URL.Path {
				case "/healthz":
					w.WriteHeader(test.healthStatus)
					body = test.healthBody
					if test.oversizeHealth {
						body += strings.Repeat("x", maxResponseBytes+1-len(body))
					}
				case "/login":
					if test.loginLocation != "" {
						w.Header().Set("Location", test.loginLocation)
					}
					w.WriteHeader(test.loginStatus)
					body = test.loginBody
					if test.oversizeLogin {
						body += strings.Repeat("x", maxResponseBytes+1-len(body))
					}
				default:
					unexpectedPath.Store(true)
					http.NotFound(w, request)
					return
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
			if count := requests.Load(); count != test.wantRequests {
				t.Errorf("received %d requests, want %d", count, test.wantRequests)
			}
			if invalidMethod.Load() {
				t.Error("probe sent a non-GET request")
			}
			if authHeaders.Load() {
				t.Error("probe sent an Authorization or Cookie header")
			}
			if unexpectedPath.Load() {
				t.Error("probe followed a redirect or requested an unexpected path")
			}
			wantPaths := []string{"GET /healthz"}
			if test.wantRequests == 2 {
				wantPaths = append(wantPaths, "GET /login")
			}
			for _, want := range wantPaths {
				select {
				case gotPath := <-paths:
					if gotPath != want {
						t.Errorf("request path = %q, want %q", gotPath, want)
					}
				case <-time.After(time.Second):
					t.Errorf("did not receive %q", want)
				}
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
		{name: "nil context", ip: ipv4, port: 8443, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 8443, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 8443, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 8443},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 8443, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
