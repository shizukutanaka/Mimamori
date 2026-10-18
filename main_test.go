package main

import (
	"io"
	"net/netip"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/lan"
	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func TestMakeFindings(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	newIP := netip.MustParseAddr("192.168.1.11")
	for _, lang := range i18n.Supported() {
		t.Run(lang, func(t *testing.T) {
			findings := makeFindings(
				[]probe.Open{
					{IP: ip, Port: 23},
					{IP: ip, Port: 2323},
					{IP: ip, Port: 21},
					{IP: ip, Port: 5900},
				},
				[]ssdp.Response{
					{IP: ip, ST: "urn:schemas-upnp-org:device:InternetGatewayDevice:1"},
					{IP: ip, USN: "uuid:router::InternetGatewayDevice"},
				},
				map[netip.Addr]ssdp.Device{},
				map[netip.Addr]string{ip: "Router", newIP: "Camera"},
				map[netip.Addr]bool{newIP: true},
				lang,
			)
			if len(findings) != 5 {
				t.Fatalf("makeFindings() returned %d findings: %#v", len(findings), findings)
			}
			counts := make(map[string]int)
			for _, finding := range findings {
				counts[finding.ID]++
				if finding.Fix == finding.ID+"_fix" || finding.Message == finding.ID {
					t.Errorf("finding fell back to a catalog key: %#v", finding)
				}
				if finding.Fix == "" || finding.Message == "" {
					t.Errorf("finding is missing localized message or fix: %#v", finding)
				}
			}
			if counts["telnet_open"] != 1 || counts["ftp_open"] != 1 || counts["vnc_open"] != 1 || counts["router_upnp"] != 1 || counts["new_device"] != 1 {
				t.Fatalf("unexpected finding counts: %#v", counts)
			}
			if got := report.Score(findings); got != 51 {
				t.Fatalf("score = %d, want 51", got)
			}
		})
	}
}

func TestResolveLanguage(t *testing.T) {
	tests := []struct {
		value string
		env   map[string]string
		want  string
		err   bool
	}{
		{"auto", map[string]string{"LANG": "en_US.UTF-8"}, "en", false},
		{"pt-BR", nil, "pt", false},
		{"zh_TW", nil, "zh-Hant", false},
		{"klingon", nil, "", true},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			got, err := resolveLanguage(test.value, func(key string) string { return test.env[key] })
			if (err != nil) != test.err {
				t.Fatalf("resolveLanguage(%q) error = %v, want error %t", test.value, err, test.err)
			}
			if got != test.want {
				t.Fatalf("resolveLanguage(%q) = %q, want %q", test.value, got, test.want)
			}
		})
	}
}

func TestRunPauseMode(t *testing.T) {
	textResult := run([]string{"--version"}, io.Discard, io.Discard)
	if !textResult.pause {
		t.Fatal("text mode should allow a console pause")
	}
	jsonResult := run([]string{"--json", "--version"}, io.Discard, io.Discard)
	if jsonResult.pause {
		t.Fatal("JSON mode should not pause")
	}
}

func TestCollectHostsFiltersAndCaps(t *testing.T) {
	neighbors := make([]lan.Neighbor, 1025)
	for i := range neighbors {
		neighbors[i].IP = netip.AddrFrom4([4]byte{10, 0, byte(i >> 8), byte(i)})
	}
	neighbors = append(neighbors, lan.Neighbor{IP: netip.MustParseAddr("8.8.8.8")})
	got := collectHosts(neighbors, netip.MustParseAddr("192.168.1.1"), []ssdp.Response{{IP: netip.MustParseAddr("172.16.0.1")}})
	if len(got) != 1024 {
		t.Fatalf("collectHosts() returned %d hosts", len(got))
	}
	for i := 1; i < len(got); i++ {
		if !got[i-1].Less(got[i]) {
			t.Fatalf("hosts are not sorted: %v before %v", got[i-1], got[i])
		}
	}
}

func TestGatewayString(t *testing.T) {
	if got := gatewayString(netip.Addr{}); got != "" {
		t.Fatalf("gatewayString(invalid) = %q", got)
	}
}
