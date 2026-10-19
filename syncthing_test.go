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

const validSyncthingMeta = `var metadata = {"authenticated":true,"deviceID":"ABCDEFG-HIJKLMN-OPQRSTU-VWXYZ12-3456789-ABCDEFG-HIJKLMN-OPQRSTU","deviceIDShort":"ABCDEFG"};`

func TestMarkUnauthenticatedSyncthing(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		id     string
		want   bool
	}{
		{name: "unauthenticated Syncthing GUI is kept", status: http.StatusOK, body: validSyncthingMeta, id: "ABCDEFG-HIJKLMN-OPQRSTU-VWXYZ12-3456789-ABCDEFG-HIJKLMN-OPQRSTU", want: true},
		{name: "unmatched metadata is dropped", status: http.StatusOK, body: `var metadata = {"deviceID":"other"};`, id: "expected-device"},
		{name: "authentication required is dropped", status: http.StatusForbidden, body: validSyncthingMeta, id: "ABCDEFG-HIJKLMN-OPQRSTU-VWXYZ12-3456789-ABCDEFG-HIJKLMN-OPQRSTU"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			syncthingOpen, requests := startSyncthingTestServer(t, test.status, test.body, test.id)
			otherPort := probe.Open{IP: syncthingOpen.IP, Port: 5900}
			opened := []probe.Open{syncthingOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedSyncthing(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedSyncthing() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedSyncthing() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /meta.js" {
					t.Fatalf("request = %q, want GET /meta.js", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Syncthing fake server received no request")
			}
		})
	}
}

func startSyncthingTestServer(t *testing.T, status int, body, id string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "8384"))
	if err != nil {
		t.Skipf("cannot listen on %s:8384: %v", address, err)
	}
	requests := make(chan string, 1)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- request.Method + " " + request.URL.Path
		w.Header().Set("X-Syncthing-ID", id)
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
			t.Error("Syncthing fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 8384}, requests
}
