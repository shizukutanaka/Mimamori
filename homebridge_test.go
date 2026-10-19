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

const validHomebridgeSettings = `{"env":{"packageName":"homebridge-config-ui-x","version":"1.0.0"},"formAuth":false}`

func TestMarkUnauthenticatedHomebridge(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "no-auth Homebridge UI is kept", status: http.StatusOK, body: validHomebridgeSettings, want: true},
		{name: "form authentication enabled is dropped", status: http.StatusOK, body: `{"env":{"packageName":"homebridge-config-ui-x"},"formAuth":true}`},
		{name: "unrelated service is dropped", status: http.StatusOK, body: `{"env":{"packageName":"other"},"formAuth":false}`},
		{name: "authentication required is dropped", status: http.StatusUnauthorized, body: `{"error":"Unauthorized"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			homebridgeOpen, requests := startHomebridgeTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: homebridgeOpen.IP, Port: 5900}
			opened := []probe.Open{homebridgeOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedHomebridge(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedHomebridge() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedHomebridge() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /api/auth/settings" {
					t.Fatalf("request = %q, want GET /api/auth/settings", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Homebridge fake server received no request")
			}
		})
	}
}

func startHomebridgeTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "8581"))
	if err != nil {
		t.Skipf("cannot listen on %s:8581: %v", address, err)
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
			t.Error("Homebridge fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 8581}, requests
}
