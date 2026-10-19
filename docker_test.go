package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
)

const validDockerVersion = `{"Platform":{"Name":"Docker Engine - Community"},"Components":[{"Name":"Engine","Version":"27.5.1","Details":{"ApiVersion":"1.47"}}],"Version":"27.5.1","ApiVersion":"1.47","Os":"linux"}`

func TestMarkUnauthenticatedDocker(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "unauthenticated Docker API is kept", status: http.StatusOK, body: validDockerVersion, want: true},
		{name: "non-Docker JSON on port 2375 is dropped", status: http.StatusOK, body: `{"version":"27.5.1"}`},
		{name: "authenticated API is dropped", status: http.StatusUnauthorized, body: `{"message":"Unauthorized"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dockerOpen, requests := startDockerTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: dockerOpen.IP, Port: 5900}
			opened := []probe.Open{dockerOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedDocker(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedDocker() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedDocker() returned %v, want %v", got, want)
				}
			}
			wantRequests := []string{"GET /version"}
			if test.want {
				wantRequests = append(wantRequests, "OPTIONS /containers/create")
			}
			for _, wantRequest := range wantRequests {
				select {
				case request := <-requests:
					if request != wantRequest {
						t.Fatalf("request = %q, want %q", request, wantRequest)
					}
				case <-time.After(time.Second):
					t.Fatalf("Docker fake server received no request; want %q", wantRequest)
				}
			}
			select {
			case request := <-requests:
				t.Fatalf("Docker fake server received unexpected request %q", request)
			default:
			}
		})
	}
}

func TestMarkUnauthenticatedDockerTLS(t *testing.T) {
	tests := []struct {
		name                     string
		requireClientCertificate bool
		want                     bool
	}{
		{name: "unauthenticated TLS Docker API is kept", want: true},
		{name: "client certificate is required and the API is dropped", requireClientCertificate: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requests := startDockerTLSTestServer(t, test.requireClientCertificate, validDockerVersion)
			dockerOpen := probe.Open{IP: netip.MustParseAddr("127.0.0.1"), Port: 2376}
			otherPort := probe.Open{IP: dockerOpen.IP, Port: 5900}
			opened := []probe.Open{dockerOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedDocker(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedDocker() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedDocker() returned %v, want %v", got, want)
				}
			}
			if test.want {
				for _, wantRequest := range []string{"GET /version", "OPTIONS /containers/create"} {
					select {
					case request := <-requests:
						if request != wantRequest {
							t.Fatalf("request = %q, want %q", request, wantRequest)
						}
					case <-time.After(time.Second):
						t.Fatalf("TLS Docker fake server received no request; want %q", wantRequest)
					}
				}
				select {
				case request := <-requests:
					t.Fatalf("TLS Docker fake server received unexpected request %q", request)
				default:
				}
			} else {
				select {
				case request := <-requests:
					t.Fatalf("client-certificate-required server received unexpected request %q", request)
				default:
				}
			}
		})
	}
}

func startDockerTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "2375"))
	if err != nil {
		t.Skipf("cannot listen on %s:2375: %v", address, err)
	}
	requests := make(chan string, 2)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- request.Method + " " + request.URL.Path
		switch request.URL.Path {
		case "/version":
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		case "/containers/create":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"page not found"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})}
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		_ = server.Close()
		select {
		case <-done:
		case <-time.After(4 * time.Second):
			t.Error("Docker fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 2375}, requests
}

func startDockerTLSTestServer(t *testing.T, requireClientCertificate bool, body string) <-chan string {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "2376"))
	if err != nil {
		t.Skipf("cannot listen on %s:2376: %v", address, err)
	}
	requests := make(chan string, 2)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- request.Method + " " + request.URL.Path
		switch request.URL.Path {
		case "/version":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
		case "/containers/create":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"page not found"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	if err := server.Listener.Close(); err != nil {
		_ = listener.Close()
		t.Fatalf("close default TLS test listener: %v", err)
	}
	server.Listener = listener
	if requireClientCertificate {
		server.TLS = &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			MinVersion: tls.VersionTLS13,
		}
	}
	server.StartTLS()
	t.Cleanup(server.Close)
	return requests
}
