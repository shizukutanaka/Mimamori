package wled

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"
)

const validInfoResponse = `{"ver":"0.15.0","vid":2406170,"leds":{"count":1,"rgbw":false,"wv":false,"cct":false}}`

func TestParseInfo(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "v0.15-like info", body: validInfoResponse, want: true},
		{name: "missing leds", body: `{"ver":"0.15.0","vid":2406170}`},
		{name: "leds is not an object", body: `{"ver":"0.15.0","vid":2406170,"leds":[]}`},
		{name: "leds is null", body: `{"ver":"0.15.0","vid":2406170,"leds":null}`},
		{name: "vid is a string", body: `{"ver":"0.15.0","vid":"2406170","leds":{}}`},
		{name: "vid is null", body: `{"ver":"0.15.0","vid":null,"leds":{}}`},
		{name: "empty version", body: `{"ver":"","vid":2406170,"leds":{}}`},
		{name: "missing version", body: `{"vid":2406170,"leds":{}}`},
		{name: "unrelated JSON", body: `{"brand":"WLED"}`},
		{name: "malformed JSON", body: `{"ver":"0.15.0","vid":`},
		{name: "JSON array", body: `[]`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseInfo([]byte(test.body)); got != test.want {
				t.Fatalf("ParseInfo() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name           string
		infoStatus     int
		infoBody       string
		infoLocation   string
		updateStatus   int
		updateBody     string
		oversizeInfo   bool
		oversizeUpdate bool
		want           bool
		wantPaths      []string
	}{
		{
			name:         "info and update available",
			infoStatus:   http.StatusOK,
			infoBody:     validInfoResponse,
			updateStatus: http.StatusOK,
			want:         true,
			wantPaths:    []string{"GET /json/info", "GET /update"},
		},
		{
			name:         "update requires authentication",
			infoStatus:   http.StatusOK,
			infoBody:     validInfoResponse,
			updateStatus: http.StatusUnauthorized,
			wantPaths:    []string{"GET /json/info", "GET /update"},
		},
		{
			name:         "update not supported",
			infoStatus:   http.StatusOK,
			infoBody:     validInfoResponse,
			updateStatus: http.StatusNotImplemented,
			wantPaths:    []string{"GET /json/info", "GET /update"},
		},
		{
			name:         "update missing",
			infoStatus:   http.StatusOK,
			infoBody:     validInfoResponse,
			updateStatus: http.StatusNotFound,
			wantPaths:    []string{"GET /json/info", "GET /update"},
		},
		{
			name:         "non-WLED info does not request update",
			infoStatus:   http.StatusOK,
			infoBody:     `{"brand":"WLED"}`,
			updateStatus: http.StatusOK,
			wantPaths:    []string{"GET /json/info"},
		},
		{
			name:         "info unauthorized",
			infoStatus:   http.StatusUnauthorized,
			infoBody:     `{"error":"unauthorized"}`,
			updateStatus: http.StatusOK,
			wantPaths:    []string{"GET /json/info"},
		},
		{
			name:         "info redirect is not followed",
			infoStatus:   http.StatusFound,
			infoBody:     validInfoResponse,
			infoLocation: "/redirect",
			updateStatus: http.StatusOK,
			wantPaths:    []string{"GET /json/info"},
		},
		{
			name:         "oversized info",
			infoStatus:   http.StatusOK,
			infoBody:     validInfoResponse,
			updateStatus: http.StatusOK,
			oversizeInfo: true,
			wantPaths:    []string{"GET /json/info"},
		},
		{
			name:           "oversized update",
			infoStatus:     http.StatusOK,
			infoBody:       validInfoResponse,
			updateStatus:   http.StatusOK,
			oversizeUpdate: true,
			wantPaths:      []string{"GET /json/info", "GET /update"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requestsMu sync.Mutex
			var requests []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requestsMu.Lock()
				requests = append(requests, request.Method+" "+request.URL.Path)
				requestsMu.Unlock()

				switch request.URL.Path {
				case "/json/info":
					if test.infoLocation != "" {
						w.Header().Set("Location", test.infoLocation)
					}
					w.WriteHeader(test.infoStatus)
					body := test.infoBody
					if test.oversizeInfo {
						body += strings.Repeat(" ", maxResponseBytes+1-len(body))
					}
					_, _ = io.WriteString(w, body)
				case "/update":
					w.WriteHeader(test.updateStatus)
					body := test.updateBody
					if test.oversizeUpdate {
						body += strings.Repeat(" ", maxResponseBytes+1-len(body))
					}
					_, _ = io.WriteString(w, body)
				case "/redirect":
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, validInfoResponse)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()
			address, err := netip.ParseAddrPort(server.Listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			got := Unauthenticated(context.Background(), address.Addr(), int(address.Port()), time.Second)
			if got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			requestsMu.Lock()
			gotRequests := append([]string(nil), requests...)
			requestsMu.Unlock()
			if len(gotRequests) != len(test.wantPaths) {
				t.Fatalf("requests = %v, want %v", gotRequests, test.wantPaths)
			}
			for index, want := range test.wantPaths {
				if gotRequests[index] != want {
					t.Errorf("request %d = %q, want %q", index, gotRequests[index], want)
				}
			}
		})
	}
}

func TestUnauthenticatedGuards(t *testing.T) {
	ipv4 := netip.MustParseAddr("127.0.0.1")
	tests := []struct {
		name    string
		ctx     context.Context
		ip      netip.Addr
		port    int
		timeout time.Duration
	}{
		{name: "nil context", ctx: nil, ip: ipv4, port: 80, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 80, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), ip: netip.Addr{}, port: 80, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, port: 0, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 80},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 80, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}

func TestRequestURL(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.5")
	tests := []struct {
		name string
		port int
		path string
		want string
	}{
		{name: "port 80 info", port: 80, path: "/json/info", want: "http://192.168.1.5/json/info"},
		{name: "port 80 update", port: 80, path: "/update", want: "http://192.168.1.5/update"},
		{name: "non-default port", port: 8080, path: "/json/info", want: "http://192.168.1.5:8080/json/info"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := requestURL(ip, test.port, test.path); got != test.want {
				t.Fatalf("requestURL() = %q, want %q", got, test.want)
			}
		})
	}
}
