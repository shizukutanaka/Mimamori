package main

import (
	"context"
	"net"
	"net/http"
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
			select {
			case request := <-requests:
				if request != "GET /version" {
					t.Fatalf("request = %q, want GET /version", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Docker fake server received no request")
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
	requests := make(chan string, 1)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- request.Method + " " + request.URL.Path
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
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
