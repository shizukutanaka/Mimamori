package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func TestMarkUnauthenticatedWeb(t *testing.T) {
	tests := []struct {
		name        string
		rootBody    string
		shellyBody  string
		wledBody    string
		piholeBody  string
		wantService string
		wantPaths   []string
	}{
		{
			name:        "Tasmota",
			rootBody:    `<form action='up' method='get'><a>Tasmota 14.0 by Theo Arends</a>`,
			shellyBody:  `{"online":true}`,
			wantService: "tasmota",
			wantPaths:   []string{"/"},
		},
		{
			name:        "Shelly",
			rootBody:    `<html>Not Tasmota</html>`,
			shellyBody:  `{"type":"SHSW-21","mac":"5ECF7F1632E8","auth":false}`,
			wantService: "shelly",
			wantPaths:   []string{"/", "/shelly"},
		},
		{
			name:        "WLED",
			rootBody:    `<html>Unknown device</html>`,
			shellyBody:  `{"online":true}`,
			wledBody:    `{"ver":"0.15.0","vid":2406170,"leds":{"count":1}}`,
			wantService: "wled",
			wantPaths:   []string{"/", "/shelly", "/json/info", "/update"},
		},
		{
			name:        "Pi-hole",
			rootBody:    `<html>Unknown device</html>`,
			shellyBody:  `{"online":true}`,
			piholeBody:  `{"session":{"valid":true,"totp":false,"sid":null,"validity":-1,"message":"no password set"},"took":0.001}`,
			wantService: "pihole",
			wantPaths:   []string{"/", "/shelly", "/json/info", "/api/auth"},
		},
		{
			name:       "unmatched",
			rootBody:   `<html>Unknown device</html>`,
			shellyBody: `{"online":true}`,
			wledBody:   `{"online":true}`,
			wantPaths:  []string{"/", "/shelly", "/json/info", "/api/auth"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var pathsMu sync.Mutex
			var paths []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				pathsMu.Lock()
				paths = append(paths, request.URL.Path)
				pathsMu.Unlock()
				if request.Method != http.MethodGet {
					t.Errorf("request method = %s, want GET", request.Method)
				}
				body := test.rootBody
				switch request.URL.Path {
				case "/shelly":
					body = test.shellyBody
				case "/json/info":
					body = test.wledBody
				case "/update":
					body = ""
				case "/api/auth":
					body = test.piholeBody
				}
				_, _ = io.WriteString(w, body)
			}))
			defer server.Close()
			address, err := netip.ParseAddrPort(server.Listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			ip := address.Addr()
			candidate := probe.Open{IP: ip, Port: 80}
			otherPort := probe.Open{IP: ip, Port: 5900}
			got := markUnauthenticatedWebWithPort(
				context.Background(),
				[]probe.Open{candidate, otherPort},
				int(address.Port()),
			)
			pathsMu.Lock()
			gotPaths := append([]string(nil), paths...)
			pathsMu.Unlock()
			if len(gotPaths) != len(test.wantPaths) {
				t.Fatalf("request paths = %v, want %v", gotPaths, test.wantPaths)
			}
			for i, path := range test.wantPaths {
				if gotPaths[i] != path {
					t.Errorf("request path %d = %q, want %q", i, gotPaths[i], path)
				}
			}
			if test.wantService == "" {
				if len(got) != 1 || got[0] != otherPort {
					t.Fatalf("markUnauthenticatedWebWithPort() = %v, want only unrelated port %v", got, otherPort)
				}
				return
			}
			want := candidate
			want.Service = test.wantService
			if len(got) != 2 || got[0] != want || got[1] != otherPort {
				t.Fatalf("markUnauthenticatedWebWithPort() = %v, want %v and %v", got, want, otherPort)
			}
		})
	}
}

func TestMarkUnauthenticatedWebCanceledContextPreservesOtherPorts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ip := netip.MustParseAddr("192.0.2.10")
	webOpen := probe.Open{IP: ip, Port: 80}
	otherOpen := probe.Open{IP: ip, Port: 5900}
	got := markUnauthenticatedWeb(ctx, []probe.Open{webOpen, otherOpen})
	if len(got) != 1 || got[0] != otherOpen {
		t.Fatalf("markUnauthenticatedWeb() = %v, want only unrelated port %v", got, otherOpen)
	}
}

func TestMarkUnauthenticatedWebStopsAfterSilentHost(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	var connections atomic.Int32
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			connections.Add(1)
			go func() {
				time.Sleep(30 * time.Second)
				conn.Close()
			}()
		}
	}()
	port := listener.Addr().(*net.TCPAddr).Port
	ip := netip.MustParseAddr("127.0.0.1")
	start := time.Now()
	got := markUnauthenticatedWebWithPort(context.Background(), []probe.Open{{IP: ip, Port: 80}}, port)
	elapsed := time.Since(start)
	if len(got) != 0 {
		t.Fatalf("markUnauthenticatedWebWithPort() = %v, want no services for a silent host", got)
	}
	if connections.Load() != 1 {
		t.Fatalf("connections = %d, want 1 (later probes skipped once the host never answered)", connections.Load())
	}
	if elapsed >= 6*time.Second {
		t.Fatalf("elapsed = %s, want a single probe timeout (under 6s)", elapsed)
	}
}
