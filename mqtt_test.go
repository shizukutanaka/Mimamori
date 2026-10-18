package main

import (
	"net/netip"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func TestMQTTAnonymousFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	for _, lang := range []string{"en", "ja"} {
		t.Run(lang, func(t *testing.T) {
			findings := makeFindings(
				[]probe.Open{{IP: ip, Port: 1883}},
				nil,
				netip.Addr{},
				map[netip.Addr]ssdp.Description{},
				map[netip.Addr]string{ip: "Kitchen hub"},
				nil,
				nil,
				lang,
			)
			if len(findings) != 1 {
				t.Fatalf("makeFindings() returned %d findings, want 1: %#v", len(findings), findings)
			}
			finding := findings[0]
			if finding.ID != "mqtt_anonymous" || finding.Level != report.Caution || finding.Points != 8 {
				t.Fatalf("MQTT finding = %#v", finding)
			}
			if finding.Message != i18n.T(lang, "mqtt_anonymous", "Kitchen hub") {
				t.Fatalf("MQTT message = %q, want %q", finding.Message, i18n.T(lang, "mqtt_anonymous", "Kitchen hub"))
			}
		})
	}
}
