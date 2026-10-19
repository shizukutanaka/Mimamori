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

const validFrigateStats = `{"detectors":{"cpu1":{"inference_speed":10.0,"detection_start":0.0,"pid":100}},"service":{"uptime":120,"version":"0.12.1-367d724","latest_version":"0.12.1","storage":{},"temperatures":{},"last_updated":1700000000},"front":{"camera_fps":5.0}}`

func TestMarkUnauthenticatedFrigate(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{
			name:   "unauthenticated Frigate API is kept",
			status: http.StatusOK,
			body:   validFrigateStats,
			want:   true,
		},
		{
			name:   "generic JSON on port 5000 is dropped",
			status: http.StatusOK,
			body:   `{"version":"2.0","service":{},"detectors":{}}`,
		},
		{
			name:   "authentication-required API is dropped",
			status: http.StatusUnauthorized,
			body:   `{"error":"Unauthorized"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			frigateOpen, requests := startFrigateTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: frigateOpen.IP, Port: 5900}
			opened := []probe.Open{frigateOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}

			got := markUnauthenticatedFrigate(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedFrigate() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedFrigate() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request.method != http.MethodGet || request.path != "/api/stats" {
					t.Fatalf("request = %s %s, want GET /api/stats", request.method, request.path)
				}
			case <-time.After(time.Second):
				t.Fatal("Frigate fake server received no API request")
			}
		})
	}
}

type frigateTestRequest struct {
	method string
	path   string
}

func startFrigateTestServer(t *testing.T, status int, body string) (probe.Open, <-chan frigateTestRequest) {
	t.Helper()
	address := "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "5000"))
	if err != nil {
		t.Skipf("cannot listen on %s:5000: %v", address, err)
	}
	requests := make(chan frigateTestRequest, 1)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- frigateTestRequest{method: request.Method, path: request.URL.Path}
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
			t.Error("Frigate fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 5000}, requests
}
