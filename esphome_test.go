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

const validESPHomeDevices = `{"configured":[],"importable":[]}`

func TestMarkUnauthenticatedESPHome(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "unauthenticated ESPHome devices are kept", status: http.StatusOK, body: validESPHomeDevices, want: true},
		{name: "generic JSON on port 6052 is dropped", status: http.StatusOK, body: `{"version":"2025.6.0"}`},
		{name: "authenticated dashboard is dropped", status: http.StatusUnauthorized, body: `{"error":"Unauthorized"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			esphomeOpen, requests := startESPHomeTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: esphomeOpen.IP, Port: 5900}
			opened := []probe.Open{esphomeOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedESPHome(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedESPHome() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedESPHome() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /devices" {
					t.Fatalf("request = %q, want GET /devices", request)
				}
			case <-time.After(time.Second):
				t.Fatal("ESPHome fake server received no request")
			}
		})
	}
}

func startESPHomeTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "6052"))
	if err != nil {
		t.Skipf("cannot listen on %s:6052: %v", address, err)
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
			t.Error("ESPHome fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 6052}, requests
}
