package main

import (
	"context"
	"net/netip"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func TestMarkUnauthenticatedTasmotaPreservesOtherPorts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ip := netip.MustParseAddr("192.0.2.10")
	tasmotaOpen := probe.Open{IP: ip, Port: 80}
	otherOpen := probe.Open{IP: ip, Port: 5900}
	got := markUnauthenticatedTasmota(ctx, []probe.Open{tasmotaOpen, otherOpen})
	if len(got) != 1 || got[0] != otherOpen {
		t.Fatalf("markUnauthenticatedTasmota() = %v, want only unrelated port %v", got, otherOpen)
	}
}
