package frigate

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"sync/atomic"
	"testing"
	"time"
)

const validStatsV012 = `{"detectors":{"cpu1":{"inference_speed":10.0,"detection_start":0.0,"pid":100}},"service":{"uptime":120,"version":"0.12.1-367d724","latest_version":"0.12.1","storage":{},"temperatures":{},"last_updated":1700000000},"front":{"camera_fps":5.0}}`

const validStatsDev = `{"cameras":{},"detectors":{},"service":{"uptime":120,"version":"dev","latest_version":"0.16.0","storage":{},"last_updated":1700000000,"retention_unmet":false}}`

func TestParseStats(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "Frigate v0.12.1 stats", body: validStatsV012, want: true},
		{name: "Frigate dev stats without temperatures", body: validStatsDev, want: true},
		{name: "empty object", body: `{}`},
		{name: "version without service", body: `{"version":"1"}`},
		{name: "service without required stats", body: `{"service":{"version":"1"}}`},
		{
			name: "missing detectors",
			body: `{"service":{"uptime":120,"version":"1","latest_version":"1","storage":{}}}`,
		},
		{
			name: "missing latest version",
			body: `{"detectors":{},"service":{"uptime":120,"version":"1","storage":{}}}`,
		},
		{
			name: "empty version",
			body: `{"detectors":{},"service":{"uptime":120,"version":"","latest_version":"1","storage":{}}}`,
		},
		{
			name: "version is not a string",
			body: `{"detectors":{},"service":{"uptime":120,"version":1,"latest_version":"1","storage":{}}}`,
		},
		{
			name: "uptime is not a number",
			body: `{"detectors":{},"service":{"uptime":"120","version":"1","latest_version":"1","storage":{}}}`,
		},
		{
			name: "storage is null",
			body: `{"detectors":{},"service":{"uptime":120,"version":"1","latest_version":"1","storage":null}}`,
		},
		{name: "JSON array", body: `[{"detectors":{},"service":{"uptime":120,"version":"1","latest_version":"1","storage":{}}}]`},
		{name: "HTML", body: `<html><body>Not Frigate</body></html>`},
		{name: "truncated JSON", body: `{"detectors":{},"service":{"uptime":120,"version":"1"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseStats([]byte(test.body)); got != test.want {
				t.Fatalf("ParseStats() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		location   string
		want       bool
		wantMethod string
		wantPath   string
	}{
		{
			name:       "Frigate stats",
			status:     http.StatusOK,
			body:       validStatsV012,
			want:       true,
			wantMethod: http.MethodGet,
			wantPath:   "/api/stats",
		},
		{
			name:       "authentication required",
			status:     http.StatusUnauthorized,
			body:       `{"error":"Unauthorized"}`,
			wantMethod: http.MethodGet,
			wantPath:   "/api/stats",
		},
		{
			name:       "AirPlay forbidden",
			status:     http.StatusForbidden,
			body:       `{"error":"Forbidden"}`,
			wantMethod: http.MethodGet,
			wantPath:   "/api/stats",
		},
		{
			name:       "Synology DSM redirect",
			status:     http.StatusFound,
			body:       `<html>Redirecting</html>`,
			location:   "/login",
			wantMethod: http.MethodGet,
			wantPath:   "/api/stats",
		},
		{
			name:       "generic JSON",
			status:     http.StatusOK,
			body:       `{"version":"2.0","service":{},"detectors":{}}`,
			wantMethod: http.MethodGet,
			wantPath:   "/api/stats",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requests := make(chan struct {
				method string
				path   string
			}, 2)
			var loginRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/login" {
					loginRequests.Add(1)
				}
				requests <- struct {
					method string
					path   string
				}{method: request.Method, path: request.URL.Path}
				if test.location != "" {
					w.Header().Set("Location", test.location)
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
				if request.method != test.wantMethod || request.path != test.wantPath {
					t.Fatalf("request = %s %s, want %s %s", request.method, request.path, test.wantMethod, test.wantPath)
				}
			case <-time.After(time.Second):
				t.Fatal("server received no stats request")
			}
			if loginRequests.Load() != 0 {
				t.Fatalf("followed redirect %d times, want no redirect", loginRequests.Load())
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
		{name: "nil context", ip: ip, port: 5000, timeout: time.Second},
		{name: "IPv6 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 5000, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ip, timeout: time.Second},
		{name: "port out of range", ctx: context.Background(), ip: ip, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ip, port: 5000},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if Unauthenticated(test.ctx, test.ip, test.port, test.timeout) {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
