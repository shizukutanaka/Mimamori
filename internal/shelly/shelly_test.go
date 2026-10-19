package shelly

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestParseInfo(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{
			name: "Gen1 unauthenticated",
			body: `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":false,"fw":"1.0"}`,
			want: true,
		},
		{
			name: "Gen1 authentication enabled",
			body: `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":true}`,
		},
		{
			name: "Gen2 unauthenticated",
			body: `{"id":"shellyplus1-123","mac":"5ECF7F1632E8","model":"SPSW-004PE16EU","gen":2,"auth_en":false}`,
			want: true,
		},
		{
			name: "Gen2 authentication enabled",
			body: `{"id":"shellyplus1-123","mac":"5ECF7F1632E8","model":"SPSW-004PE16EU","gen":2,"auth_en":true}`,
		},
		{
			name: "generation is a string",
			body: `{"model":"SPSW-004PE16EU","gen":"2","auth_en":false}`,
		},
		{
			name: "missing Gen1 authentication field",
			body: `{"type":"SHSW-21","mac":"5ECF7F1632E8"}`,
		},
		{
			name: "authentication is a string",
			body: `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":"false"}`,
		},
		{
			name: "null authentication",
			body: `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":null}`,
		},
		{
			name: "missing Gen2 authentication field",
			body: `{"model":"SPSW-004PE16EU","gen":2}`,
		},
		{
			name: "Gen2 generation below two",
			body: `{"model":"SPSW-004PE16EU","gen":1,"auth_en":false}`,
		},
		{
			name: "Gen2 model is empty",
			body: `{"model":"","gen":2,"auth_en":false}`,
		},
		{
			name: "unrelated JSON",
			body: `{"online":true}`,
		},
		{
			name: "JSON array",
			body: `[]`,
		},
		{
			name: "non-JSON",
			body: `<html>Not a Shelly</html>`,
		},
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
		name     string
		status   int
		body     string
		location string
		oversize bool
		want     bool
	}{
		{
			name:   "Gen1 unauthenticated",
			status: http.StatusOK,
			body:   `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":false}`,
			want:   true,
		},
		{
			name:   "Gen2 unauthenticated",
			status: http.StatusOK,
			body:   `{"id":"shellyplus1-123","model":"SPSW-004PE16EU","gen":2,"auth_en":false}`,
			want:   true,
		},
		{
			name:   "unauthorized",
			status: http.StatusUnauthorized,
			body:   `{"error":"unauthorized"}`,
		},
		{
			name:     "redirect",
			status:   http.StatusFound,
			body:     `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":false}`,
			location: "/redirect",
		},
		{
			name:     "oversized response",
			status:   http.StatusOK,
			body:     `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":false}`,
			oversize: true,
		},
		{
			name:   "unrelated JSON",
			status: http.StatusOK,
			body:   `{"online":true}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests, redirectRequests int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/redirect" {
					redirectRequests++
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":false}`)
					return
				}
				requests++
				if request.Method != http.MethodGet || request.URL.Path != "/shelly" {
					t.Errorf("request = %s %s, want GET /shelly", request.Method, request.URL.Path)
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

func TestUnauthenticatedHTTPSRedirect(t *testing.T) {
	tests := []struct {
		name     string
		location string
		body     string
		want     bool
		wantTLS  bool
	}{
		{
			name:     "Gen2 unauthenticated",
			location: "https://127.0.0.1:{tlsPort}/",
			body:     `{"model":"SPSW-004PE16EU","gen":2,"auth_en":false}`,
			want:     true,
			wantTLS:  true,
		},
		{
			name:     "Gen2 authentication enabled",
			location: "https://127.0.0.1:{tlsPort}/",
			body:     `{"model":"SPSW-004PE16EU","gen":2,"auth_en":true}`,
			wantTLS:  true,
		},
		{
			name:     "different IP",
			location: "https://127.0.0.2:{tlsPort}/",
			body:     `{"model":"SPSW-004PE16EU","gen":2,"auth_en":false}`,
		},
		{
			name:     "HTTP URL",
			location: "http://127.0.0.1:{tlsPort}/",
			body:     `{"model":"SPSW-004PE16EU","gen":2,"auth_en":false}`,
		},
		{
			name:     "wrong port",
			location: "https://127.0.0.1:{wrongPort}/",
			body:     `{"model":"SPSW-004PE16EU","gen":2,"auth_en":false}`,
		},
		{
			name: "missing Location",
			body: `{"model":"SPSW-004PE16EU","gen":2,"auth_en":false}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tlsRequests atomic.Int32
			tlsServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				tlsRequests.Add(1)
				if request.Method != http.MethodGet || request.URL.Path != "/shelly" {
					t.Errorf("TLS request = %s %s, want GET /shelly", request.Method, request.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, test.body)
			}))
			defer tlsServer.Close()
			var wrongTLSRequests atomic.Int32
			wrongTLSServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				wrongTLSRequests.Add(1)
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{"model":"SPSW-004PE16EU","gen":2,"auth_en":false}`)
			}))
			defer wrongTLSServer.Close()
			tlsAddress, err := netip.ParseAddrPort(tlsServer.Listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			tlsPort := int(tlsAddress.Port())
			wrongTLSAddress, err := netip.ParseAddrPort(wrongTLSServer.Listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			location := strings.NewReplacer(
				"{tlsPort}", strconv.Itoa(tlsPort),
				"{wrongPort}", strconv.Itoa(int(wrongTLSAddress.Port())),
			).Replace(test.location)

			var httpRequests atomic.Int32
			httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				httpRequests.Add(1)
				if request.Method != http.MethodGet || request.URL.Path != "/shelly" {
					t.Errorf("HTTP request = %s %s, want GET /shelly", request.Method, request.URL.Path)
				}
				if location != "" {
					w.Header().Set("Location", location)
				}
				w.WriteHeader(http.StatusFound)
			}))
			defer httpServer.Close()
			httpAddress, err := netip.ParseAddrPort(httpServer.Listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}

			got := unauthenticated(
				context.Background(),
				httpAddress.Addr(),
				int(httpAddress.Port()),
				tlsPort,
				time.Second,
			)
			if got != test.want {
				t.Fatalf("unauthenticated() = %t, want %t", got, test.want)
			}
			if got := httpRequests.Load(); got != 1 {
				t.Errorf("HTTP request count = %d, want exactly one", got)
			}
			wantTLSRequests := int32(0)
			if test.wantTLS {
				wantTLSRequests = 1
			}
			if got := tlsRequests.Load(); got != wantTLSRequests {
				t.Errorf("TLS request count = %d, want %d", got, wantTLSRequests)
			}
			if got := wrongTLSRequests.Load(); got != 0 {
				t.Errorf("wrong-port TLS request count = %d, want zero", got)
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
		{name: "default HTTP port omitted", port: 80, want: "http://192.168.1.5/shelly"},
		{name: "non-default port included", port: 8080, want: "http://192.168.1.5:8080/shelly"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := requestURL(ip, test.port); got != test.want {
				t.Fatalf("requestURL() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestHTTPSRequestURL(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.5")
	tests := []struct {
		name string
		port int
		want string
	}{
		{name: "default HTTPS port omitted", port: 443, want: "https://192.168.1.5/shelly"},
		{name: "non-default port included", port: 8443, want: "https://192.168.1.5:8443/shelly"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := httpsRequestURL(ip, test.port); got != test.want {
				t.Fatalf("httpsRequestURL() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestRedirectPortMatches(t *testing.T) {
	tests := []struct {
		name         string
		locationPort string
		httpsPort    int
		want         bool
	}{
		{name: "443 omitted", httpsPort: 443, want: true},
		{name: "443 explicit", locationPort: "443", httpsPort: 443, want: true},
		{name: "443 mismatched", locationPort: "8443", httpsPort: 443},
		{name: "non-default explicit", locationPort: "8443", httpsPort: 8443, want: true},
		{name: "non-default omitted", httpsPort: 8443},
		{name: "non-default redirected to 443", locationPort: "443", httpsPort: 8443},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := redirectPortMatches(test.locationPort, test.httpsPort); got != test.want {
				t.Fatalf("redirectPortMatches(%q, %d) = %t, want %t", test.locationPort, test.httpsPort, got, test.want)
			}
		})
	}
}
