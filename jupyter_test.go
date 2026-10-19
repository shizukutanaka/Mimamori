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

const validJupyterStatus = `{"connections": 0, "kernels": 1, "last_activity": "2026-10-18T10:00:00.000000Z", "started": "2026-10-18T09:00:00.000000Z"}`

func TestMarkUnauthenticatedJupyter(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "unauthenticated Jupyter is kept", status: http.StatusOK, body: validJupyterStatus, want: true},
		{name: "forbidden Jupyter is dropped", status: http.StatusForbidden, body: validJupyterStatus},
		{name: "public version endpoint is dropped", status: http.StatusOK, body: `{"version":"2.14.0"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jupyterOpen, requests := startJupyterTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: jupyterOpen.IP, Port: 5900}
			opened := []probe.Open{jupyterOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedJupyter(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedJupyter() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedJupyter() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /api/status" {
					t.Fatalf("request = %q, want GET /api/status", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Jupyter fake server received no request")
			}
		})
	}
}

func startJupyterTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "8888"))
	if err != nil {
		t.Skipf("cannot listen on %s:8888: %v", address, err)
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
			t.Error("Jupyter fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 8888}, requests
}
