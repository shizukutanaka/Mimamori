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

const validHomeAssistantOnboarding = `[{"step":"user","done":false},{"step":"core_config","done":false},{"step":"analytics","done":false},{"step":"integration","done":false}]`

func TestMarkOpenHomeAssistantOnboarding(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "open owner onboarding is kept", status: http.StatusOK, body: validHomeAssistantOnboarding, want: true},
		{name: "non-Home Assistant JSON is dropped", status: http.StatusOK, body: `{"status":"ok"}`},
		{name: "missing onboarding endpoint is dropped", status: http.StatusNotFound, body: `{"message":"Not Found"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			onboardingOpen, requests := startHomeAssistantTestServer(t, test.status, test.body)
			otherPort := probe.Open{IP: onboardingOpen.IP, Port: 5900}
			opened := []probe.Open{onboardingOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markOpenHomeAssistantOnboarding(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markOpenHomeAssistantOnboarding() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markOpenHomeAssistantOnboarding() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /api/onboarding" {
					t.Fatalf("request = %q, want GET /api/onboarding", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Home Assistant fake server received no request")
			}
		})
	}
}

func startHomeAssistantTestServer(t *testing.T, status int, body string) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "8123"))
	if err != nil {
		t.Skipf("cannot listen on %s:8123: %v", address, err)
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
			t.Error("Home Assistant fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 8123}, requests
}
