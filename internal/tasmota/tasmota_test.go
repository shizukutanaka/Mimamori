package tasmota

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

const (
	latestAdminPage = `<!doctype html><form id=but5 style="display:block;" action='up' method='get'></form><a href='https://github.com/arendst/Tasmota'>Tasmota 14.0.0 by Theo Arends</a>`
	v9AdminPage     = `<p><form action='up' method='get'><a href='https://bit.ly/tasmota'>Tasmota 9.1.0 by Theo Arends</a>`
	germanAdminPage = `<form action='up' method='get'><a href='https://github.com/arendst/Tasmota'>Tasmota 14.0.0 von Theo Arends</a>`
	userModePage    = `<a href='https://github.com/arendst/Tasmota'>Tasmota 14.0.0 by Theo Arends</a>`
)

func TestParsePage(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "latest admin page", body: latestAdminPage, want: true},
		{name: "version 9 admin page", body: v9AdminPage, want: true},
		{name: "localized German footer", body: germanAdminPage, want: true},
		{name: "user mode page", body: userModePage},
		{name: "missing author", body: `>Tasmota 14.0.0 action='up'`},
		{name: "missing Tasmota footer", body: `Tasmota 14.0.0 Theo Arends action='up'`},
		{name: "missing upgrade form", body: `>Tasmota 14.0.0 by Theo Arends`},
		{name: "empty page", body: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParsePage([]byte(test.body)); got != test.want {
				t.Fatalf("ParsePage() = %t, want %t", got, test.want)
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
		{name: "latest admin page", status: http.StatusOK, body: latestAdminPage, want: true},
		{name: "version 9 admin page", status: http.StatusOK, body: v9AdminPage, want: true},
		{name: "German footer", status: http.StatusOK, body: germanAdminPage, want: true},
		{name: "user mode page", status: http.StatusOK, body: userModePage},
		{name: "unauthorized", status: http.StatusUnauthorized, body: latestAdminPage},
		{name: "unrelated HTML", status: http.StatusOK, body: `<html><body>Welcome</body></html>`},
		{name: "redirect", status: http.StatusFound, body: latestAdminPage, location: "/redirect"},
		{name: "oversized page", status: http.StatusOK, body: latestAdminPage, oversize: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests, redirectRequests int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/redirect" {
					redirectRequests++
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, latestAdminPage)
					return
				}
				requests++
				if request.Method != http.MethodGet || request.URL.Path != "/" {
					t.Errorf("request = %s %s, want GET /", request.Method, request.URL.Path)
				}
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				w.WriteHeader(test.status)
				body := test.body
				if test.oversize {
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
		{name: "nil context", ctx: nil, ip: ipv4, port: 80, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 80, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), ip: netip.Addr{}, port: 80, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, port: 0, timeout: time.Second},
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

func TestRequestURL(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.5")
	tests := []struct {
		name string
		port int
		want string
	}{
		{name: "default HTTP port omitted", port: 80, want: "http://192.168.1.5/"},
		{name: "non-default port included", port: 8080, want: "http://192.168.1.5:8080/"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := requestURL(ip, test.port); got != test.want {
				t.Fatalf("requestURL() = %q, want %q", got, test.want)
			}
		})
	}
}
