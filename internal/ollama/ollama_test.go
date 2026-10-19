package ollama

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestParseRoot(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "exact response", status: http.StatusOK, body: "Ollama is running", want: true},
		{name: "trailing newline", status: http.StatusOK, body: "Ollama is running\n", want: true},
		{name: "different response", status: http.StatusOK, body: "Ollama is running!"},
		{name: "not found", status: http.StatusNotFound, body: "Ollama is running"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseRoot(test.status, []byte(test.body)); got != test.want {
				t.Fatalf("ParseRoot() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{name: "version response", status: http.StatusOK, body: `{"version":"0.12.3"}`, want: true},
		{name: "empty version", status: http.StatusOK, body: `{"version":""}`},
		{name: "missing version", status: http.StatusOK, body: `{}`},
		{name: "version is a number", status: http.StatusOK, body: `{"version":12}`},
		{name: "not an object", status: http.StatusOK, body: `["0.12.3"]`},
		{name: "not found", status: http.StatusNotFound, body: `{"version":"0.12.3"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseVersion(test.status, []byte(test.body)); got != test.want {
				t.Fatalf("ParseVersion() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	tests := []struct {
		name          string
		rootStatus    int
		rootBody      string
		versionStatus int
		versionBody   string
		redirect      bool
		oversize      bool
		want          bool
		wantRequests  int32
	}{
		{name: "open Ollama", rootStatus: http.StatusOK, rootBody: "Ollama is running", versionStatus: http.StatusOK, versionBody: `{"version":"0.12.3"}`, want: true, wantRequests: 2},
		{name: "root does not match", rootStatus: http.StatusOK, rootBody: "not Ollama", wantRequests: 1},
		{name: "root status rejected", rootStatus: http.StatusNotFound, rootBody: "Ollama is running", wantRequests: 1},
		{name: "redirect not followed", rootStatus: http.StatusFound, rootBody: "Ollama is running", redirect: true, wantRequests: 1},
		{name: "oversized root body", rootStatus: http.StatusOK, rootBody: "Ollama is running", oversize: true, wantRequests: 1},
		{name: "version rejected", rootStatus: http.StatusOK, rootBody: "Ollama is running", versionStatus: http.StatusUnauthorized, versionBody: `{"version":"0.12.3"}`, wantRequests: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidMethod atomic.Bool
			var redirectRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.Method != http.MethodGet {
					invalidMethod.Store(true)
				}
				switch request.URL.Path {
				case "/":
					if test.redirect {
						w.Header().Set("Location", "/api/version")
						http.Redirect(w, request, "/api/version", http.StatusFound)
						return
					}
					w.WriteHeader(test.rootStatus)
					body := test.rootBody
					if test.oversize {
						body += strings.Repeat(" ", maxResponseBytes+1-len(body))
					}
					_, _ = io.WriteString(w, body)
				case "/api/version":
					w.WriteHeader(test.versionStatus)
					_, _ = io.WriteString(w, test.versionBody)
				default:
					redirectRequests.Add(1)
				}
			}))
			defer server.Close()
			address, err := netip.ParseAddrPort(server.Listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			got := Open(context.Background(), address.Addr(), int(address.Port()), time.Second)
			if got != test.want {
				t.Fatalf("Open() = %t, want %t", got, test.want)
			}
			if count := requests.Load(); count != test.wantRequests {
				t.Errorf("received %d requests, want %d", count, test.wantRequests)
			}
			if invalidMethod.Load() {
				t.Error("probe sent a non-GET request")
			}
			if redirectRequests.Load() != 0 {
				t.Error("probe followed a redirect")
			}
		})
	}
}

func TestOpenOversizedVersion(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests.Add(1)
		if request.URL.Path == "/" {
			_, _ = io.WriteString(w, "Ollama is running")
			return
		}
		_, _ = io.WriteString(w, `{"version":"0.12.3"}`+strings.Repeat(" ", maxResponseBytes))
	}))
	defer server.Close()
	address, err := netip.ParseAddrPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	if Open(context.Background(), address.Addr(), int(address.Port()), time.Second) {
		t.Fatal("Open() = true for an oversized version response")
	}
	if requests.Load() != 2 {
		t.Fatalf("received %d requests, want 2", requests.Load())
	}
}

func TestOpenGuards(t *testing.T) {
	ipv4 := netip.MustParseAddr("127.0.0.1")
	tests := []struct {
		name    string
		ctx     context.Context
		ip      netip.Addr
		port    int
		timeout time.Duration
	}{
		{name: "nil context", ip: ipv4, port: 11434, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 11434, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 11434, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 11434},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 11434, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Open(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Open() = true, want false")
			}
		})
	}
}
