package main

import (
	"context"
	"net/netip"
	"strings"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func TestMarkUnauthenticatedADBWithoutADBEntries(t *testing.T) {
	opened := []probe.Open{
		{IP: netip.MustParseAddr("192.168.1.20"), Port: 21},
	}
	got := markUnauthenticatedADB(context.Background(), opened)
	if len(got) != len(opened) || got[0] != opened[0] {
		t.Fatalf("markUnauthenticatedADB() = %#v, want %#v", got, opened)
	}
	if &got[0] != &opened[0] {
		t.Fatal("markUnauthenticatedADB() copied input without port 5555 entries")
	}
}

func TestMakeFindingsForUnauthenticatedADB(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.20")
	device := "Living room TV"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: adbNoAuthPort}},
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
	if finding.ID != "adb_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("makeFindings() finding = %#v, want adb_noauth Danger/25", finding)
	}
	if !strings.Contains(finding.Message, device) {
		t.Fatalf("finding message %q does not contain device name %q", finding.Message, device)
	}
	if finding.ID == "adb_open" {
		t.Fatal("unauthenticated ADB was reported as adb_open")
	}
}
