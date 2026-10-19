package nodered

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

const validSettings = `{"httpNodeRoot":"/","version":"4.0.9","context":{"default":"memory","stores":["memory"]},"libraries":[],"flowEncryptionType":"system"}`

func TestParseSettings(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "no user", body: validSettings, want: true},
		{name: "Node-RED 1.3.7", body: `{"httpNodeRoot":"/","version":"1.3.7"}`, want: true},
		{name: "wildcard permission", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":"*"}}`, want: true},
		{name: "read permission only", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":"read"}}`},
		{name: "wildcard read permission only", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":"*.read"}}`},
		{name: "write permission", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":"write"}}`, want: true},
		{name: "wildcard write permission", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":"*.write"}}`, want: true},
		{name: "flows write permission", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":"flows.write"}}`, want: true},
		{name: "array without write permission", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":["flows.read"]}}`},
		{name: "array with write permission", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":["flows.read","flows.write"]}}`, want: true},
		{name: "empty permission array", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":[]}}`},
		{name: "user without permissions", body: `{"httpNodeRoot":"/","version":"1","user":{}}`},
		{name: "numeric permissions", body: `{"httpNodeRoot":"/","version":"1","user":{"permissions":1}}`},
		{name: "null user", body: `{"httpNodeRoot":"/","version":"1","user":null}`},
		{name: "empty object", body: `{}`},
		{name: "missing httpNodeRoot", body: `{"version":"1"}`},
		{name: "httpNodeRoot is not a string", body: `{"httpNodeRoot":1,"version":"1"}`},
		{name: "empty version", body: `{"httpNodeRoot":"/","version":""}`},
		{name: "version is not a string", body: `{"httpNodeRoot":"/","version":1}`},
		{name: "array", body: `[{"httpNodeRoot":"/","version":"1"}]`},
		{name: "HTML", body: `<html><body>Not Node-RED</body></html>`},
		{name: "truncated JSON", body: `{"httpNodeRoot":"/","version":`},
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
		name     string
		status   int
		body     string
		location string
		want     bool
	}{
		{name: "unauthenticated settings", status: http.StatusOK, body: validSettings, want: true},
		{name: "authentication required", status: http.StatusUnauthorized, body: `{"error":"Unauthorized"}`},
		{name: "redirect", status: http.StatusFound, body: `<html>Redirecting</html>`, location: "/login"},
		{name: "generic JSON", status: http.StatusOK, body: `{"version":"1","other":"service"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var loginRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/login" {
					loginRequests.Add(1)
				}
				if request.Method != http.MethodGet || request.URL.Path != "/settings" {
					t.Errorf("request = %s %s, want GET /settings", request.Method, request.URL.Path)
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
