package main

import (
	"context"
	"net/netip"
	"strings"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func TestMarkUnauthenticatedVNCWithoutVNCEntries(t *testing.T) {
	opened := []probe.Open{
		{IP: netip.MustParseAddr("192.168.1.20"), Port: 21},
	}
	got := markUnauthenticatedVNC(context.Background(), opened)
	if len(got) != len(opened) || got[0] != opened[0] {
		t.Fatalf("markUnauthenticatedVNC() = %#v, want %#v", got, opened)
	}
	if &got[0] != &opened[0] {
		t.Fatal("markUnauthenticatedVNC() copied input without port 5900 entries")
	}
}

func TestMakeFindingsForUnauthenticatedVNC(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.20")
	device := "Living room computer"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: vncNoAuthPort}},
		nil,
		netip.Addr{},
		nil,
		map[netip.Addr]string{ip: device},
		nil,
		nil,
		"ja",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want 1: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "vnc_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("makeFindings() finding = %#v, want vnc_noauth Danger/25", finding)
	}
	if !strings.Contains(finding.Message, device) {
		t.Fatalf("finding message %q does not contain device name %q", finding.Message, device)
	}
	if finding.ID == "vnc_open" {
		t.Fatal("no-password VNC was reported as vnc_open")
	}
}
