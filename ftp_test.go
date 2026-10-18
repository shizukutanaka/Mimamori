package main

import (
	"net/netip"
	"reflect"
	"strings"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func TestMarkAnonymousFTP(t *testing.T) {
	ipA := netip.MustParseAddr("192.168.1.10")
	ipB := netip.MustParseAddr("192.168.1.11")
	ftpA := report.Finding{
		ID: "ftp_open", Level: report.Caution, Points: 8, IP: ipA.String(), Device: "ルーター",
		Message: "original message", Fix: "original fix",
	}
	ftpB := report.Finding{
		ID: "ftp_open", Level: report.Caution, Points: 8, IP: ipB.String(), Device: "NAS",
		Message: "original message", Fix: "original fix",
	}
	telnetA := report.Finding{
		ID: "telnet_open", Level: report.Danger, Points: 25, IP: ipA.String(), Device: "ルーター",
		Message: "telnet message", Fix: "telnet fix",
	}
	findings := []report.Finding{ftpA, ftpB, telnetA}

	got := markAnonymousFTP(findings, map[netip.Addr]bool{ipA: true}, "ja")
	want := []report.Finding{
		{
			ID: "ftp_anonymous", Level: report.Danger, Points: 25, IP: ipA.String(), Device: ftpA.Device,
			Message: i18n.T("ja", "ftp_anonymous", ftpA.Device),
			Fix:     i18n.T("ja", "ftp_anonymous_fix"),
		},
		ftpB,
		telnetA,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("markAnonymousFTP() = %#v, want %#v", got, want)
	}
	if !strings.Contains(got[0].Message, ftpA.Device) {
		t.Fatalf("Japanese message %q does not contain device %q", got[0].Message, ftpA.Device)
	}
}
