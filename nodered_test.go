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

const validNodeREDSettings = `{"httpNodeRoot":"/","version":"4.0.9","context":{"default":"memory","stores":["memory"]},"libraries":[],"flowEncryptionType":"system"}`

func TestMarkUnauthenticatedNodeRED(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "unauthenticated Node-RED settings are kept", status: http.StatusOK, body: validNodeREDSettings, want: true},
		{name: "generic JSON on port 1880 is dropped", status: http.StatusOK, body: `{"version":"1"}`},
		{name: "authenticated editor is dropped", status: http.StatusUnauthorized, body: `{"error":"Unauthorized"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodeREDOpen, requests := startNodeREDTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: nodeREDOpen.IP, Port: 5900}
			opened := []probe.Open{nodeREDOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedNodeRED(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedNodeRED() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedNodeRED() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /settings" {
					t.Fatalf("request = %q, want GET /settings", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Node-RED fake server received no request")
			}
		})
	}
}

func startNodeREDTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "1880"))
	if err != nil {
		t.Skipf("cannot listen on %s:1880: %v", address, err)
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
			t.Error("Node-RED fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 1880}, requests
}
