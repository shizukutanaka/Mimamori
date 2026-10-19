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

const validTransmissionConflict = `<html><body><h1>409: Conflict</h1><p>Your request has been rejected for reasons listed below:</p><ul><li><code>X-Transmission-Session-Id: session123</code></li></ul></body></html>`

func TestMarkUnauthenticatedTransmission(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		id     string
		want   bool
	}{
		{name: "unauthenticated RPC is kept", status: http.StatusConflict, body: validTransmissionConflict, id: "session123", want: true},
		{name: "mismatched body is dropped", status: http.StatusConflict, body: validTransmissionConflict, id: "another-session"},
		{name: "unauthorized RPC is dropped", status: http.StatusUnauthorized, body: validTransmissionConflict, id: "session123"},
		{name: "forbidden RPC is dropped", status: http.StatusForbidden, body: validTransmissionConflict, id: "session123"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transmissionOpen, requests := startTransmissionTestServer(t, test.status, test.body, test.id)
			otherPort := probe.Open{IP: transmissionOpen.IP, Port: 5900}
			opened := []probe.Open{transmissionOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedTransmission(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedTransmission() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedTransmission() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /transmission/rpc" {
					t.Fatalf("request = %q, want GET /transmission/rpc", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Transmission fake server received no request")
			}
		})
	}
}

func startTransmissionTestServer(t *testing.T, status int, body, id string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "9091"))
	if err != nil {
		t.Skipf("cannot listen on %s:9091: %v", address, err)
	}
	requests := make(chan string, 1)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- request.Method + " " + request.URL.Path
		if id != "" {
			w.Header().Set("X-Transmission-Session-Id", id)
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
			t.Error("Transmission fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 9091}, requests
}
