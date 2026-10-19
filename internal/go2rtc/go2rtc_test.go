package go2rtc

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"
)

func TestParseAPIInfo(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{
			name:   "go2rtc API info",
			status: http.StatusOK,
			body:   `{"host":"192.168.1.5:1984","revision":"a1b2c3d","version":"1.9.9"}`,
			want:   true,
		},
		{
			name:   "unauthorized",
			status: http.StatusUnauthorized,
			body:   `{"host":"192.168.1.5:1984","revision":"a1b2c3d","version":"1.9.9"}`,
		},
		{
			name:   "missing revision",
			status: http.StatusOK,
			body:   `{"version":"1.9.9"}`,
		},
		{
			name:   "empty version",
			status: http.StatusOK,
			body:   `{"revision":"a1b2c3d","version":""}`,
		},
		{
			name:   "non-JSON",
			status: http.StatusOK,
			body:   `not JSON`,
		},
		{
			name:   "JSON array",
			status: http.StatusOK,
			body:   `[{"revision":"a1b2c3d","version":"1.9.9"}]`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseAPIInfo(test.status, []byte(test.body)); got != test.want {
				t.Fatalf("ParseAPIInfo() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		authHeader string
		want       bool
	}{
		{
			name:   "go2rtc API info",
			status: http.StatusOK,
			body:   `{"host":"127.0.0.1:1984","revision":"a1b2c3d","version":"1.9.9"}`,
			want:   true,
		},
		{
			name:       "authentication required",
			status:     http.StatusUnauthorized,
			body:       `{"error":"Unauthorized"}`,
			authHeader: `Basic realm="go2rtc"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requests := make(chan struct {
				method string
				path   string
			}, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests <- struct {
					method string
					path   string
				}{method: request.Method, path: request.URL.Path}
				if test.authHeader != "" {
					w.Header().Set("Www-Authenticate", test.authHeader)
				}
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()

			port := server.Listener.Addr().(*net.TCPAddr).Port
			ip := netip.MustParseAddr("127.0.0.1")
			if got := Unauthenticated(context.Background(), ip, port, time.Second); got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			select {
			case request := <-requests:
				if request.method != http.MethodGet || request.path != "/api" {
					t.Fatalf("request = %s %s, want GET /api", request.method, request.path)
				}
			case <-time.After(time.Second):
				t.Fatal("server received no API request")
			}
		})
	}
}

func TestUnauthenticatedRejectsInvalidArguments(t *testing.T) {
	ip := netip.MustParseAddr("127.0.0.1")
	tests := []struct {
		name    string
		ctx     context.Context
		ip      netip.Addr
		port    int
		timeout time.Duration
	}{
		{name: "nil context", ip: ip, port: 1984, timeout: time.Second},
		{name: "IPv6 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 1984, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ip, port: 0, timeout: time.Second},
		{name: "port out of range", ctx: context.Background(), ip: ip, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ip, port: 1984},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if Unauthenticated(test.ctx, test.ip, test.port, test.timeout) {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
