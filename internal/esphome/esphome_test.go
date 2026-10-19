package esphome

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

const validLegacyDevices = `{"configured":[{"name":"kitchen","friendly_name":"Kitchen","configuration":"kitchen.yaml","loaded_integrations":["api","wifi"],"deployed_version":"2024.12.0","current_version":"2024.12.0","path":"/config/kitchen.yaml","comment":null,"address":"kitchen.local","web_port":null,"target_platform":"ESP32"}],"importable":[{"name":"m5stack-atom","friendly_name":"M5Stack Atom","package_import_url":"github://esphome/x.yaml","project_name":"esphome.atom","project_version":"1.0","network":"wifi"}]}`

func TestParseDevices(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "empty device lists", body: `{"configured":[],"importable":[]}`, want: true},
		{name: "legacy dashboard devices", body: validLegacyDevices, want: true},
		{name: "empty object", body: `{}`},
		{name: "missing importable", body: `{"configured":[]}`},
		{name: "configured is an object", body: `{"configured":{},"importable":[]}`},
		{name: "configured is null", body: `{"configured":null,"importable":[]}`},
		{name: "JSON array", body: `[{"configured":[],"importable":[]}]`},
		{name: "HTML", body: `<html><body>Not ESPHome</body></html>`},
		{name: "truncated JSON", body: `{"configured":[],"importable":`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseDevices([]byte(test.body)); got != test.want {
				t.Fatalf("ParseDevices() = %t, want %t", got, test.want)
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
		want     bool
	}{
		{name: "device list", status: http.StatusOK, body: `{"configured":[],"importable":[]}`, want: true},
		{name: "authentication required", status: http.StatusUnauthorized, body: `{"error":"Unauthorized"}`},
		{name: "legacy login redirect", status: http.StatusFound, body: `<html>Redirecting</html>`, location: "./login"},
		{name: "generic JSON", status: http.StatusOK, body: `{"version":"2025.6.0"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var loginRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/login" {
					loginRequests.Add(1)
				}
				if request.Method != http.MethodGet || request.URL.Path != "/devices" {
					t.Errorf("request = %s %s, want GET /devices", request.Method, request.URL.Path)
				}
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()
			port := server.Listener.Addr().(*net.TCPAddr).Port
			if got := Unauthenticated(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second); got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			if loginRequests.Load() != 0 {
				t.Fatalf("followed redirect %d times", loginRequests.Load())
			}
		})
	}
}
