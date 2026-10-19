package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/netip"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func TestMarkOpenOllama(t *testing.T) {
	listener, err := net.Listen("tcp4", net.JoinHostPort("127.0.0.1", "11434"))
	if err != nil {
		t.Skipf("cannot listen on 127.0.0.1:11434: %v", err)
	}
	requests := make(chan string, 2)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests <- request.Method + " " + request.URL.Path
		switch request.URL.Path {
		case "/":
			_, _ = w.Write([]byte("Ollama is running"))
		case "/api/version":
			_ = json.NewEncoder(w).Encode(map[string]string{"version": "0.12.3"})
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
			t.Error("Ollama fake server did not stop")
		}
	})

	ip := netip.MustParseAddr("127.0.0.1")
	ollamaOpen := probe.Open{IP: ip, Port: 11434}
	otherPort := probe.Open{IP: ip, Port: 5900}
	opened := []probe.Open{ollamaOpen, otherPort}
	got := markOpenOllama(context.Background(), opened)
	if len(got) != len(opened) || got[0] != ollamaOpen || got[1] != otherPort {
		t.Fatalf("markOpenOllama() = %v, want %v", got, opened)
	}
	for _, want := range []string{"GET /", "GET /api/version"} {
		select {
		case request := <-requests:
			if request != want {
				t.Fatalf("request = %q, want %q", request, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("fake Ollama did not receive %q", want)
		}
	}
}
