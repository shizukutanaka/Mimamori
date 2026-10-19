package docker

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"reflect"
	"sync"
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
		name          string
		versionStatus int
		versionBody   string
		optionsStatus int
		optionsBody   string
		location      string
		want          bool
		wantRequests  []string
	}{
		{
			name:          "write request forbidden",
			versionStatus: http.StatusOK,
			versionBody:   modernVersion,
			optionsStatus: http.StatusForbidden,
			optionsBody:   `{"message":"Forbidden"}`,
			wantRequests:  []string{"GET /version", "OPTIONS /containers/create"},
		},
		{
			name:          "write request requires authentication",
			versionStatus: http.StatusOK,
			versionBody:   modernVersion,
			optionsStatus: http.StatusUnauthorized,
			optionsBody:   `{"message":"Unauthorized"}`,
			wantRequests:  []string{"GET /version", "OPTIONS /containers/create"},
		},
		{
			name:          "write request returns Docker not found",
			versionStatus: http.StatusOK,
			versionBody:   modernVersion,
			optionsStatus: http.StatusNotFound,
			optionsBody:   `{"message":"page not found"}`,
			want:          true,
			wantRequests:  []string{"GET /version", "OPTIONS /containers/create"},
		},
		{
			name:          "write request method is not allowed",
			versionStatus: http.StatusOK,
			versionBody:   modernVersion,
			optionsStatus: http.StatusMethodNotAllowed,
			optionsBody:   `{"message":"method not allowed"}`,
			want:          true,
			wantRequests:  []string{"GET /version", "OPTIONS /containers/create"},
		},
		{
			name:          "write request redirect is not followed",
			versionStatus: http.StatusOK,
			versionBody:   modernVersion,
			optionsStatus: http.StatusFound,
			optionsBody:   `<html>Redirecting</html>`,
			location:      "/redirect-target",
			want:          true,
			wantRequests:  []string{"GET /version", "OPTIONS /containers/create"},
		},
		{
			name:          "non-Docker JSON does not send a second request",
			versionStatus: http.StatusOK,
			versionBody:   `{"version":"27.5.1"}`,
			wantRequests:  []string{"GET /version"},
		},
		{
			name:          "version not found does not send a second request",
			versionStatus: http.StatusNotFound,
			versionBody:   `{"message":"page not found"}`,
			wantRequests:  []string{"GET /version"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mu sync.Mutex
			var requests []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requestBody, err := io.ReadAll(request.Body)
				if err != nil {
					t.Errorf("read request body: %v", err)
				}
				if len(requestBody) != 0 {
					t.Errorf("request body = %q, want empty", requestBody)
				}
				mu.Lock()
				requests = append(requests, request.Method+" "+request.URL.Path)
				mu.Unlock()
				switch request.URL.Path {
				case "/version":
					w.WriteHeader(test.versionStatus)
					_, _ = w.Write([]byte(test.versionBody))
				case "/containers/create":
					if test.location != "" {
						w.Header().Set("Location", test.location)
					}
					w.WriteHeader(test.optionsStatus)
					_, _ = w.Write([]byte(test.optionsBody))
				default:
					w.WriteHeader(http.StatusOK)
				}
			}))
			defer server.Close()
			port := server.Listener.Addr().(*net.TCPAddr).Port
			if got := Unauthenticated(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second); got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			mu.Lock()
			gotRequests := append([]string(nil), requests...)
			mu.Unlock()
			if !reflect.DeepEqual(gotRequests, test.wantRequests) {
				t.Fatalf("requests = %v, want %v", gotRequests, test.wantRequests)
			}
		})
	}
}
