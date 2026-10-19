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

func TestMarkUnauthenticatedGo2RTC(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		auth   bool
		want   bool
	}{
		{
			name:   "unauthenticated go2rtc API is kept",
			status: http.StatusOK,
			body:   `{"host":"127.0.0.1:1984","revision":"a1b2c3d","version":"1.9.9"}`,
			want:   true,
		},
		{
			name:   "authentication-required go2rtc API is dropped",
			status: http.StatusUnauthorized,
			body:   `{"error":"Unauthorized"}`,
			auth:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			apiOpen, requests := startGo2RCTestServer(t, test.status, test.body, test.auth)
			otherPort := probe.Open{IP: apiOpen.IP, Port: 5900}
			opened := []probe.Open{apiOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}

			got := markUnauthenticatedGo2RTC(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedGo2RTC() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedGo2RTC() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request.method != http.MethodGet || request.path != "/api" {
					t.Fatalf("request = %s %s, want GET /api", request.method, request.path)
				}
			case <-time.After(time.Second):
				t.Fatal("go2rtc fake server received no API request")
			}
		})
	}
}

type go2RTCTestRequest struct {
	method string
	path   string
}

func startGo2RCTestServer(t *testing.T, status int, body string, auth bool) (probe.Open, <-chan go2RTCTestRequest) {
	t.Helper()
	address := "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "1984"))
	if err != nil {
		t.Skipf("cannot listen on %s:1984: %v", address, err)
	}
	requests := make(chan go2RTCTestRequest, 1)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- go2RTCTestRequest{method: request.Method, path: request.URL.Path}
		if auth {
			w.Header().Set("Www-Authenticate", `Basic realm="go2rtc"`)
		}
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
			t.Error("go2rtc fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 1984}, requests
}
