package syncthing

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

const (
	testDeviceID = "ABCDEFG-HIJKLMN-OPQRSTU-VWXYZ12-3456789-ABCDEFG-HIJKLMN-OPQRSTU"
	validMeta    = `var metadata = {"authenticated":true,"deviceID":"ABCDEFG-HIJKLMN-OPQRSTU-VWXYZ12-3456789-ABCDEFG-HIJKLMN-OPQRSTU","deviceIDShort":"ABCDEFG"};`
)

func TestParseMeta(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		headerID string
		want     bool
	}{
		{name: "real metadata response", body: validMeta, headerID: testDeviceID, want: true},
		{name: "trailing newline", body: validMeta + "\n", headerID: testDeviceID, want: true},
		{name: "authenticated field is not required", body: `var metadata = {"deviceID":"ABCDEFG-HIJKLMN-OPQRSTU-VWXYZ12-3456789-ABCDEFG-HIJKLMN-OPQRSTU"};`, headerID: testDeviceID, want: true},
		{name: "missing prefix", body: `metadata = {"deviceID":"` + testDeviceID + `"};`, headerID: testDeviceID},
		{name: "missing semicolon", body: `var metadata = {"deviceID":"` + testDeviceID + `"}`, headerID: testDeviceID},
		{name: "malformed JSON", body: `var metadata = {invalid};`, headerID: testDeviceID},
		{name: "device ID missing", body: `var metadata = {};`, headerID: testDeviceID},
		{name: "device ID empty", body: `var metadata = {"deviceID":""};`, headerID: testDeviceID},
		{name: "device ID is not a string", body: `var metadata = {"deviceID":123};`, headerID: testDeviceID},
		{name: "device ID differs from header", body: `var metadata = {"deviceID":"another-device"};`, headerID: testDeviceID},
		{name: "header ID is empty", body: `var metadata = {"deviceID":"` + testDeviceID + `"};`, headerID: ""},
		{name: "top-level array", body: `var metadata = [{"deviceID":"` + testDeviceID + `"}];`, headerID: testDeviceID},
		{name: "leading whitespace", body: " " + validMeta, headerID: testDeviceID},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseMeta([]byte(test.body), test.headerID); got != test.want {
				t.Fatalf("ParseMeta() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		headerID string
		location string
		oversize bool
		want     bool
	}{
		{name: "no password", status: http.StatusOK, body: validMeta, headerID: testDeviceID, want: true},
		{name: "device ID mismatch", status: http.StatusOK, body: `var metadata = {"deviceID":"other"};`, headerID: testDeviceID},
		{name: "forbidden", status: http.StatusForbidden, body: validMeta, headerID: testDeviceID},
		{name: "unauthorized", status: http.StatusUnauthorized, body: validMeta, headerID: testDeviceID},
		{name: "redirect not followed", status: http.StatusFound, body: validMeta, headerID: testDeviceID, location: "/redirect"},
		{name: "oversized body", status: http.StatusOK, body: validMeta, headerID: testDeviceID, oversize: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			var invalidRequest atomic.Bool
			var forbiddenHeader atomic.Bool
			var redirectRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				if request.URL.Path == "/redirect" {
					redirectRequests.Add(1)
					w.Header().Set("X-Syncthing-ID", test.headerID)
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, validMeta)
					return
				}
				if request.Method != http.MethodGet || request.URL.Path != "/meta.js" {
					invalidRequest.Store(true)
				}
				for _, name := range []string{"X-API-Key", "X-CSRF-Token", "Authorization", "Origin", "Referer"} {
					if request.Header.Get(name) != "" {
						forbiddenHeader.Store(true)
					}
				}
				w.Header().Set("X-Syncthing-ID", test.headerID)
				if test.location != "" {
					w.Header().Set("Location", test.location)
				}
				w.WriteHeader(test.status)
				body := test.body
				if test.oversize {
					body += strings.Repeat(" ", maxResponseBytes+1-len(body))
				}
				_, _ = io.WriteString(w, body)
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
			if count := requests.Load(); count != 1 {
				t.Errorf("received %d requests, want exactly one", count)
			}
			if invalidRequest.Load() {
				t.Error("request was not GET /meta.js")
			}
			if forbiddenHeader.Load() {
				t.Error("request included an authentication, CSRF, API-key, Origin, or Referer header")
			}
			if count := redirectRequests.Load(); count != 0 {
				t.Errorf("followed redirect %d times", count)
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
		{name: "nil context", ip: ipv4, port: 8384, timeout: time.Second},
		{name: "non-IPv4 address", ctx: context.Background(), ip: netip.MustParseAddr("::1"), port: 8384, timeout: time.Second},
		{name: "invalid address", ctx: context.Background(), port: 8384, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ipv4, timeout: time.Second},
		{name: "negative port", ctx: context.Background(), ip: ipv4, port: -1, timeout: time.Second},
		{name: "port above range", ctx: context.Background(), ip: ipv4, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ipv4, port: 8384},
		{name: "negative timeout", ctx: context.Background(), ip: ipv4, port: 8384, timeout: -time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Unauthenticated(test.ctx, test.ip, test.port, test.timeout); got {
				t.Fatal("Unauthenticated() = true, want false")
			}
		})
	}
}
