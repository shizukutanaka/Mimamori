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

const validScryptedLogin = `{"error":"Not logged in.","hasLogin":false,"externalAddresses":[],"addresses":["https://192.168.1.5:10443"],"hostname":"scrypted"}`

func TestMarkUnclaimedScrypted(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "unclaimed Scrypted server is kept", status: http.StatusOK, body: validScryptedLogin, want: true},
		{name: "claimed Scrypted server is dropped", status: http.StatusOK, body: `{"hasLogin":true,"addresses":[],"hostname":"scrypted"}`},
		{name: "unrelated service is dropped", status: http.StatusOK, body: `{"hasLogin":false,"addresses":[],"hostname":1}`},
		{name: "not found is dropped", status: http.StatusNotFound, body: `{"error":"Not Found"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scryptedOpen, requests := startScryptedTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: scryptedOpen.IP, Port: 5900}
			opened := []probe.Open{scryptedOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnclaimedScrypted(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnclaimedScrypted() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnclaimedScrypted() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /login" {
					t.Fatalf("request = %q, want GET /login", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Scrypted fake server received no request")
			}
		})
	}
}

func startScryptedTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "11080"))
	if err != nil {
		t.Skipf("cannot listen on %s:11080: %v", address, err)
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
			t.Error("Scrypted fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 11080}, requests
}
