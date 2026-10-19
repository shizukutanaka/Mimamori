package docker

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

const (
	modernVersion = `{"Platform":{"Name":"Docker Engine - Community"},"Components":[{"Name":"Engine","Version":"27.5.1","Details":{"ApiVersion":"1.47"}}],"Version":"27.5.1","ApiVersion":"1.47","Os":"linux"}`
	legacyVersion = `{"Version":"1.13.1","ApiVersion":"1.25","Os":"linux"}`
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "modern response", body: modernVersion, want: true},
		{name: "old 1.x response", body: legacyVersion, want: true},
		{name: "missing API version", body: `{"Version":"27.5.1","Os":"linux"}`},
		{name: "empty OS", body: `{"Version":"27.5.1","ApiVersion":"1.47","Os":""}`},
		{name: "version is not a string", body: `{"Version":27,"ApiVersion":"1.47","Os":"linux"}`},
		{name: "JSON array", body: `[{"Version":"27.5.1","ApiVersion":"1.47","Os":"linux"}]`},
		{name: "JSON null", body: `null`},
		{name: "non-JSON", body: `<html><body>Not Docker</body></html>`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseVersion([]byte(test.body)); got != test.want {
				t.Fatalf("ParseVersion() = %t, want %t", got, test.want)
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
		{name: "unauthenticated Docker version", status: http.StatusOK, body: modernVersion, want: true},
		{name: "authentication required", status: http.StatusUnauthorized, body: `{"message":"Unauthorized"}`},
		{name: "non-Docker JSON", status: http.StatusOK, body: `{"version":"27.5.1"}`},
		{name: "redirect", status: http.StatusFound, body: `<html>Redirecting</html>`, location: "/login"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.Method != http.MethodGet || request.URL.Path != "/version" {
					t.Errorf("request = %s %s, want GET /version", request.Method, request.URL.Path)
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
			if got := requests.Load(); got != 1 {
				t.Fatalf("request count = %d, want exactly one", got)
			}
		})
	}
}
