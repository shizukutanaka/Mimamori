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

func TestMarkUnauthenticatedCodeServer(t *testing.T) {
	listener, err := net.Listen("tcp4", net.JoinHostPort("127.0.0.1", "8443"))
	if err != nil {
		t.Skipf("cannot listen on 127.0.0.1:8443: %v", err)
	}
	requests := make(chan string, 2)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- request.Method + " " + request.URL.Path
		switch request.URL.Path {
		case "/healthz":
			_, _ = w.Write([]byte(`{"status":"alive","lastHeartbeat":123}`))
		case "/login":
			w.Header().Set("Location", "./")
			w.WriteHeader(http.StatusFound)
		default:
			http.NotFound(w, request)
		}
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
			t.Error("code-server fake server did not stop")
		}
	})

	ip := netip.MustParseAddr("127.0.0.1")
	codeServerOpen := probe.Open{IP: ip, Port: 8443}
	otherPort := probe.Open{IP: ip, Port: 5900}
	opened := []probe.Open{codeServerOpen, otherPort}
	got := markUnauthenticatedCodeServer(context.Background(), opened)
	if len(got) != len(opened) || got[0] != codeServerOpen || got[1] != otherPort {
		t.Fatalf("markUnauthenticatedCodeServer() = %v, want %v", got, opened)
	}
	for _, want := range []string{"GET /healthz", "GET /login"} {
		select {
		case request := <-requests:
			if request != want {
				t.Fatalf("request = %q, want %q", request, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("fake code-server did not receive %q", want)
		}
	}
}
