package homeassistant

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

const validOnboarding = `[{"step":"user","done":false},{"step":"core_config","done":false},{"step":"analytics","done":false},{"step":"integration","done":false}]`

func TestParseOnboarding(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "fresh onboarding", body: validOnboarding, want: true},
		{name: "only user step not done", body: `[{"step":"user","done":false},{"step":"core_config","done":true},{"step":"analytics","done":true},{"step":"integration","done":true}]`, want: true},
		{name: "user done but other steps are not", body: `[{"step":"user","done":true},{"step":"core_config","done":false},{"step":"analytics","done":false},{"step":"integration","done":false}]`},
		{name: "all steps done", body: `[{"step":"user","done":true},{"step":"core_config","done":true},{"step":"analytics","done":true},{"step":"integration","done":true}]`},
		{name: "missing user step", body: `[{"step":"core_config","done":false},{"step":"analytics","done":false},{"step":"integration","done":false}]`},
		{name: "done is a string", body: `[{"step":"user","done":"false"}]`},
		{name: "step is missing", body: `[{"done":false}]`},
		{name: "empty array", body: `[]`},
		{name: "JSON object", body: `{"step":"user","done":false}`},
		{name: "JSON null", body: `null`},
		{name: "non-JSON", body: `<html><body>Not Home Assistant</body></html>`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseOnboarding([]byte(test.body)); got != test.want {
				t.Fatalf("ParseOnboarding() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUserStepOpen(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		location string
		want     bool
	}{
		{name: "open user onboarding", status: http.StatusOK, body: validOnboarding, want: true},
		{name: "onboarding endpoint missing", status: http.StatusNotFound, body: `{"message":"Not Found"}`},
		{name: "redirect", status: http.StatusFound, body: `<html>Redirecting</html>`, location: "/login"},
		{name: "non-Home Assistant JSON", status: http.StatusOK, body: `{"status":"ok"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.Method != http.MethodGet || request.URL.Path != "/api/onboarding" {
					t.Errorf("request = %s %s, want GET /api/onboarding", request.Method, request.URL.Path)
				}
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()
			port := server.Listener.Addr().(*net.TCPAddr).Port
			if got := UserStepOpen(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second); got != test.want {
				t.Fatalf("UserStepOpen() = %t, want %t", got, test.want)
			}
			if got := requests.Load(); got != 1 {
				t.Fatalf("request count = %d, want exactly one", got)
			}
		})
	}
}
