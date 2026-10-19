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

const validMoonrakerInfo = `{"result":{"moonraker_version":"v0.9.3-12-gabcdef","api_version":[1,5,0],"api_version_string":"1.5.0"}}`

func TestMarkUnauthenticatedMoonraker(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "trusted Moonraker info is kept", status: http.StatusOK, body: validMoonrakerInfo, want: true},
		{name: "authentication required is dropped", status: http.StatusUnauthorized, body: `{"error":{"code":401,"message":"Unauthorized"}}`},
		{name: "generic JSON is dropped", status: http.StatusOK, body: `{"result":{"version":"1.0"}}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			moonrakerOpen, requests := startMoonrakerTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: moonrakerOpen.IP, Port: 5900}
			opened := []probe.Open{moonrakerOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedMoonraker(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedMoonraker() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedMoonraker() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /server/info" {
					t.Fatalf("request = %q, want GET /server/info", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Moonraker fake server received no request")
			}
		})
	}
}

func startMoonrakerTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "7125"))
	if err != nil {
		t.Skipf("cannot listen on %s:7125: %v", address, err)
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
			t.Error("Moonraker fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 7125}, requests
}
