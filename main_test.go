package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/shizukutanaka/Mimamori/internal/baseline"
	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/lan"
	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
	"github.com/shizukutanaka/Mimamori/internal/wifi"
)

func TestExitCode(t *testing.T) {
	danger := report.Finding{Level: report.Danger}
	caution := report.Finding{Level: report.Caution}
	tests := []struct {
		name              string
		findings          []report.Finding
		loadState         historyLoadState
		historySaveFailed bool
		want              int
	}{
		{name: "danger and history save failed", findings: []report.Finding{danger}, loadState: historyLoadReady, historySaveFailed: true, want: 2},
		{name: "caution and history save failed", findings: []report.Finding{caution}, loadState: historyLoadReady, historySaveFailed: true, want: 1},
		{name: "danger and history saved", findings: []report.Finding{danger}, loadState: historyLoadReady, want: 2},
		{name: "caution and history saved", findings: []report.Finding{caution}, loadState: historyLoadReady, want: 0},
		{name: "no findings and history saved", loadState: historyLoadReady, want: 0},
		{name: "caution and history unreadable", findings: []report.Finding{caution}, loadState: historyLoadUnavailable, want: 1},
		{name: "no findings and history unreadable", loadState: historyLoadUnavailable, want: 1},
		{name: "danger and history unreadable", findings: []report.Finding{danger}, loadState: historyLoadUnavailable, want: 2},
		{name: "quarantined history", loadState: historyLoadQuarantined, want: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := exitCode(test.findings, test.loadState, test.historySaveFailed); got != test.want {
				t.Fatalf("exitCode() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestLoadHistoryPathError(t *testing.T) {
	t.Run("returns path error without creating a file", func(t *testing.T) {
		tempDir := t.TempDir()
		workingDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.Chdir(workingDir); err != nil {
				t.Errorf("restore working directory: %v", err)
			}
		})

		wantErr := errors.New("no home")
		_, existed, err := loadHistory("", wantErr)
		if err != wantErr {
			t.Fatalf("loadHistory() error = %v, want %v", err, wantErr)
		}
		if existed {
			t.Fatal("loadHistory() reported history exists, want false")
		}
		entries, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			t.Fatalf("working directory contains files: %v", entries)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing.json")
		_, existed, err := loadHistory(path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if existed {
			t.Fatal("loadHistory() reported history exists, want false")
		}
	})
}

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
				ip,
				map[netip.Addr]ssdp.Description{},
				map[netip.Addr]string{ip: "Router", newIP: "Camera"},
				map[netip.Addr]string{newIP: "new_device"},
				nil,
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

func TestVulnerableUPnPComponent(t *testing.T) {
	tests := []struct {
		name       string
		server     string
		want       string
		vulnerable bool
	}{
		{name: "libupnp 1.6.6", server: "Portable SDK for UPnP devices/1.6.6", want: "libupnp 1.6.6", vulnerable: true},
		{name: "libupnp 1.6.17", server: "Portable SDK for UPnP devices/1.6.17", want: "libupnp 1.6.17", vulnerable: true},
		{name: "libupnp 1.6.18", server: "Portable SDK for UPnP devices/1.6.18", want: "libupnp 1.6.18", vulnerable: true},
		{name: "libupnp 1.6.20", server: "Portable SDK for UPnP devices/1.6.20", want: "libupnp 1.6.20", vulnerable: true},
		{name: "libupnp threshold", server: "Portable SDK for UPnP devices/1.6.21"},
		{name: "libupnp newer minor", server: "Portable SDK for UPnP devices/1.8.0"},
		{name: "libupnp numeric minor", server: "Portable SDK for UPnP devices/1.14.0"},
		{name: "Intel SDK case insensitive", server: "intel sdk for upnp devices/", want: "Intel SDK for UPnP devices (all versions)", vulnerable: true},
		{name: "embedded Intel SDK", server: "FooIntel SDK for UPnP devices/"},
		{name: "MiniUPnPd 1.0", server: "miniupnpd/1.0", want: "MiniUPnPd 1.0", vulnerable: true},
		{name: "MiniUPnPd 1.3", server: "MINIUPNPD/1.3", want: "MiniUPnPd 1.3", vulnerable: true},
		{name: "MiniUPnPd after server tokens", server: "Linux/2.6 UPnP/1.0 miniupnpd/1.3", want: "MiniUPnPd 1.3", vulnerable: true},
		{name: "MiniUPnPd after labelled token", server: "Linux/2.6 UPnP/1.0 (daemon:miniupnpd/1.3)", want: "MiniUPnPd 1.3", vulnerable: true},
		{name: "MiniUPnPd in parentheses", server: "(miniupnpd/1.3)", want: "MiniUPnPd 1.3", vulnerable: true},
		{name: "embedded MiniUPnPd", server: "notminiupnpd/1.3"},
		{name: "MiniUPnPd after underscore", server: "x_miniupnpd/1.3"},
		{name: "MiniUPnPd after digit", server: "2miniupnpd/1.3"},
		{name: "Portable SDK after comma", server: "Linux, UPnP/1.0, Portable SDK for UPnP devices/1.6.6", want: "libupnp 1.6.6", vulnerable: true},
		{name: "Portable SDK after labelled token", server: "Linux UPnP/1.0 sdk:Portable SDK for UPnP devices/1.6.6", want: "libupnp 1.6.6", vulnerable: true},
		{name: "embedded Portable SDK", server: "xPortable SDK for UPnP devices/1.6.6"},
		{name: "MiniUPnPd threshold", server: "miniupnpd/1.4"},
		{name: "MiniUPnPd newer major", server: "miniupnpd/2.1"},
		{name: "unparseable portable version", server: "Portable SDK for UPnP devices/1.6.x"},
		{name: "unparseable miniupnpd version", server: "miniupnpd/1.x"},
		{name: "garbage", server: "unrelated UPnP/1.0"},
		{name: "empty", server: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, vulnerable := vulnerableUPnPComponent(test.server)
			if got != test.want || vulnerable != test.vulnerable {
				t.Fatalf("vulnerableUPnPComponent(%q) = (%q, %t), want (%q, %t)",
					test.server, got, vulnerable, test.want, test.vulnerable)
			}
		})
	}
}

func TestVulnerableUPnPStackFindingsDeduplicateAndScoreByGateway(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	client := netip.MustParseAddr("192.168.1.20")
	responses := []ssdp.Response{
		{IP: gateway, Server: "miniupnpd/1.0"},
		{IP: gateway, Server: "Portable SDK for UPnP devices/1.6.6"},
		{IP: client, Server: "Portable SDK for UPnP devices/1.6.17"},
	}
	findings := vulnerableUPnPStackFindings(
		responses, gateway, map[netip.Addr]string{gateway: "Home Router"}, "en",
	)
	if len(findings) != 2 {
		t.Fatalf("vulnerableUPnPStackFindings() returned %d findings, want 2: %#v", len(findings), findings)
	}
	if finding := findings[0]; finding.ID != "upnp_stack_vulnerable" ||
		finding.Level != report.Danger || finding.Points != 15 ||
		finding.IP != gateway.String() || finding.Device != "Home Router" {
		t.Errorf("gateway finding = %#v, want Danger/15 for Home Router", finding)
	}
	if finding := findings[1]; finding.ID != "upnp_stack_vulnerable" ||
		finding.Level != report.Caution || finding.Points != 8 ||
		finding.IP != client.String() || finding.Device != client.String() {
		t.Errorf("client finding = %#v, want Caution/8 with IP fallback", finding)
	}
	if got := report.Score(findings); got != 77 {
		t.Errorf("report.Score() = %d, want 77", got)
	}
}

func TestVulnerableUPnPStackMessageFormatting(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.1")
	tests := []struct {
		lang        string
		device      string
		wantMessage string
		wantFix     string
	}{
		{
			lang:        "en",
			device:      "Router",
			wantMessage: "Router uses an old UPnP component (libupnp 1.6.6) with known remotely exploitable vulnerabilities.",
			wantFix:     "Update the device's firmware. If no update exists, turn off UPnP on it (and on the router), and consider replacing it.",
		},
		{
			lang:        "ja",
			device:      "ルーター",
			wantMessage: "「ルーター」は、遠隔から悪用できる既知の脆弱性がある古い UPnP 部品（libupnp 1.6.6）を使っています。",
			wantFix:     "ファームウェアを更新してください。更新がない場合は、その機器とルーターの UPnP を無効にし、買い替えも検討してください。",
		},
	}
	for _, test := range tests {
		t.Run(test.lang, func(t *testing.T) {
			findings := vulnerableUPnPStackFindings(
				[]ssdp.Response{{IP: ip, Server: "Portable SDK for UPnP devices/1.6.6"}},
				ip, map[netip.Addr]string{ip: test.device}, test.lang,
			)
			if len(findings) != 1 {
				t.Fatalf("vulnerableUPnPStackFindings() returned %d findings, want 1", len(findings))
			}
			if findings[0].Message != test.wantMessage || findings[0].Fix != test.wantFix {
				t.Fatalf("finding text = (%q, %q), want (%q, %q)",
					findings[0].Message, findings[0].Fix, test.wantMessage, test.wantFix)
			}
		})
	}
}

func TestADBOpenFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Android TV"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 5555}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		map[netip.Addr]string{},
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want exactly one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "adb_open" || finding.Level != report.Caution || finding.Points != 8 {
		t.Fatalf("ADB finding = %#v, want adb_open/Caution/8", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("ADB finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "adb_open", deviceName) || finding.Fix != i18n.T("en", "adb_open_fix") {
		t.Fatalf("ADB finding text = %#v, want localized message and fix", finding)
	}
}

func TestGo2RTCNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Home Assistant"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 1984}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "go2rtc_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("go2rtc finding = %#v, want go2rtc_noauth/Danger/25", finding)
	}
	if !strings.Contains(finding.Message, deviceName) {
		t.Fatalf("go2rtc message = %q, want it to contain %q", finding.Message, deviceName)
	}
}

func TestFrigateNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Frigate NVR"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 5000}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "frigate_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Frigate finding = %#v, want frigate_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Frigate finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "frigate_noauth", deviceName) || finding.Fix != i18n.T("en", "frigate_noauth_fix") {
		t.Fatalf("Frigate finding text = %#v, want localized message and fix", finding)
	}
}

func TestNodeREDNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Node-RED"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 1880}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "nodered_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Node-RED finding = %#v, want nodered_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Node-RED finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "nodered_noauth", deviceName) || finding.Fix != i18n.T("en", "nodered_noauth_fix") {
		t.Fatalf("Node-RED finding text = %#v, want localized message and fix", finding)
	}
}

func TestHomebridgeNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Homebridge"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 8581}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "homebridge_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Homebridge finding = %#v, want homebridge_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Homebridge finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "homebridge_noauth", deviceName) || finding.Fix != i18n.T("en", "homebridge_noauth_fix") {
		t.Fatalf("Homebridge finding text = %#v, want localized message and fix", finding)
	}
}

func TestESPHomeNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "ESPHome"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 6052}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "esphome_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("ESPHome finding = %#v, want esphome_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("ESPHome finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "esphome_noauth", deviceName) || finding.Fix != i18n.T("en", "esphome_noauth_fix") {
		t.Fatalf("ESPHome finding text = %#v, want localized message and fix", finding)
	}
}

func TestZigbee2MQTTNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Zigbee2MQTT"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 8080}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "zigbee2mqtt_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Zigbee2MQTT finding = %#v, want zigbee2mqtt_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Zigbee2MQTT finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "zigbee2mqtt_noauth", deviceName) || finding.Fix != i18n.T("en", "zigbee2mqtt_noauth_fix") {
		t.Fatalf("Zigbee2MQTT finding text = %#v, want localized message and fix", finding)
	}
}

func TestMoonrakerNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Moonraker"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 7125}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "moonraker_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Moonraker finding = %#v, want moonraker_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Moonraker finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "moonraker_noauth", deviceName) || finding.Fix != i18n.T("en", "moonraker_noauth_fix") {
		t.Fatalf("Moonraker finding text = %#v, want localized message and fix", finding)
	}
}

func TestTasmotaNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Smart plug"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 80}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "tasmota_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Tasmota finding = %#v, want tasmota_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Tasmota finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "tasmota_noauth", deviceName) || finding.Fix != i18n.T("en", "tasmota_noauth_fix") {
		t.Fatalf("Tasmota finding text = %#v, want localized message and fix", finding)
	}
}

func TestShellyNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Smart plug"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 80, Service: "shelly"}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "shelly_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Shelly finding = %#v, want shelly_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Shelly finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "shelly_noauth", deviceName) || finding.Fix != i18n.T("en", "shelly_noauth_fix") {
		t.Fatalf("Shelly finding text = %#v, want localized message and fix", finding)
	}
}

func TestWLEDNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Living room lights"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 80, Service: "wled"}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "wled_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("WLED finding = %#v, want wled_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("WLED finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "wled_noauth", deviceName) || finding.Fix != i18n.T("en", "wled_noauth_fix") {
		t.Fatalf("WLED finding text = %#v, want localized message and fix", finding)
	}
}

func TestDockerNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Home Server"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 2375}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "docker_noauth" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Docker finding = %#v, want docker_noauth/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Docker finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "docker_noauth", deviceName) || finding.Fix != i18n.T("en", "docker_noauth_fix") {
		t.Fatalf("Docker finding text = %#v, want localized message and fix", finding)
	}
}

func TestDockerTLSNoVerifyFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Home Server"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 2376}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "docker_tls_noverify" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Docker TLS finding = %#v, want docker_tls_noverify/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Docker TLS finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "docker_tls_noverify", deviceName) || finding.Fix != i18n.T("en", "docker_tls_noverify_fix") {
		t.Fatalf("Docker TLS finding text = %#v, want localized message and fix", finding)
	}
}

func TestHomeAssistantOnboardingFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	deviceName := "Home Assistant"
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 8123}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: deviceName},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "homeassistant_onboarding" || finding.Level != report.Danger || finding.Points != 25 {
		t.Fatalf("Home Assistant finding = %#v, want homeassistant_onboarding/Danger/25", finding)
	}
	if finding.IP != ip.String() || finding.Device != deviceName {
		t.Fatalf("Home Assistant finding target = %#v, want IP %s and device %q", finding, ip, deviceName)
	}
	if finding.Message != i18n.T("en", "homeassistant_onboarding", deviceName) || finding.Fix != i18n.T("en", "homeassistant_onboarding_fix") {
		t.Fatalf("Home Assistant finding text = %#v, want localized message and fix", finding)
	}
}

func TestRTSP8554Finding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.23")
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 8554}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: "Home Assistant"},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	finding := findings[0]
	if finding.ID != "rtsp_noauth" || finding.Level != report.Caution || finding.Points != 8 {
		t.Fatalf("RTSP finding = %#v, want rtsp_noauth/Caution/8", finding)
	}
}

func TestRTSPFindingDeduplicatedAcrossPorts(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.23")
	findings := makeFindings(
		[]probe.Open{{IP: ip, Port: 554}, {IP: ip, Port: 8554}},
		nil,
		netip.Addr{},
		map[netip.Addr]ssdp.Description{},
		map[netip.Addr]string{ip: "Home Assistant"},
		nil,
		nil,
		"en",
	)
	if len(findings) != 1 {
		t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
	}
	if findings[0].ID != "rtsp_noauth" || findings[0].IP != ip.String() {
		t.Fatalf("RTSP finding = %#v, want one rtsp_noauth finding for %s", findings[0], ip)
	}
}

func TestNewDeviceFindingClassification(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.11")
	deviceName := "My phone"
	tests := []struct {
		name       string
		firstRun   bool
		isNew      bool
		mac        string
		wantID     string
		wantLevel  report.Level
		wantPoints int
	}{
		{
			name: "randomized new device", isNew: true, mac: "02:00:00:00:00:01",
			wantID: "new_device_private", wantLevel: report.Caution, wantPoints: 3,
		},
		{
			name: "universal new device", isNew: true, mac: "00:11:22:33:44:55",
			wantID: "new_device", wantLevel: report.Caution, wantPoints: 3,
		},
		{
			name: "first run", firstRun: true, isNew: true, mac: "02:00:00:00:00:01",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findingID := newDeviceFindingID(test.firstRun, test.isNew, test.mac)
			if findingID != test.wantID {
				t.Fatalf("newDeviceFindingID() = %q, want %q", findingID, test.wantID)
			}
			newDeviceFindings := map[netip.Addr]string{}
			if findingID != "" {
				newDeviceFindings[ip] = findingID
			}
			for _, lang := range []string{"en", "ja"} {
				findings := makeFindings(
					nil, nil, netip.Addr{}, map[netip.Addr]ssdp.Description{},
					map[netip.Addr]string{ip: deviceName}, newDeviceFindings, nil, lang,
				)
				if test.wantID == "" {
					if len(findings) != 0 {
						t.Fatalf("first-run findings = %#v, want none", findings)
					}
					continue
				}
				if len(findings) != 1 {
					t.Fatalf("findings = %#v, want one", findings)
				}
				finding := findings[0]
				if finding.ID != test.wantID || finding.Level != test.wantLevel || finding.Points != test.wantPoints {
					t.Fatalf("finding = %#v, want %s/%s/%d", finding, test.wantID, test.wantLevel, test.wantPoints)
				}
				if got := report.Score(findings); got != 97 {
					t.Errorf("score = %d, want 97", got)
				}
				if got := report.Status(report.Score(findings), findings); got != "caution" {
					t.Errorf("status = %q, want caution", got)
				}
				if test.wantID == "new_device_private" {
					wantMessages := map[string]string{
						"en": "A device using a private (changing) address connected: My phone. Phones and PCs change this address over time, so it may be one of your own devices.",
						"ja": "プライベートアドレス（変化するアドレス）を使う機器が接続しました: My phone。スマホや PC はこのアドレスを定期的に変えるため、ご家族の機器の可能性があります。",
					}
					wantFixes := map[string]string{
						"en": "Check whether everyone's phones/PCs are accounted for. If you see devices you don't recognize repeatedly, change the Wi-Fi password and reconnect only your own devices.",
						"ja": "家族のスマホや PC か確認してください。心当たりのない機器が繰り返し出る場合は、Wi-Fi のパスワードを変更し、自分の機器だけ再接続してください。",
					}
					if finding.Message != wantMessages[lang] || finding.Fix != wantFixes[lang] {
						t.Errorf("%s private finding text = (%q, %q), want (%q, %q)", lang, finding.Message, finding.Fix, wantMessages[lang], wantFixes[lang])
					}
				}
			}
		})
	}
}

func TestPortForwardFindings(t *testing.T) {
	internalIP := netip.MustParseAddr("192.168.1.30")
	publicIP := netip.MustParseAddr("8.8.8.8")
	tests := []struct {
		name      string
		mappings  []ssdp.PortMapping
		names     map[netip.Addr]string
		wantID    string
		wantLevel report.Level
		wantCount int
		wantPort  int
		wantName  string
	}{
		{
			name: "risky remote access",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 44338, Protocol: "TCP", InternalPort: 3389, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Office PC"},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 44338, wantName: "Office PC",
		},
		{
			name: "restricted TCP remote host is caution",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 44338, Protocol: "TCP", RemoteHost: "203.0.113.8", InternalPort: 3389, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Office PC"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 44338, wantName: "Office PC",
		},
		{
			name: "UDP risky port is caution",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 22, Protocol: "UDP", InternalPort: 22, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Office PC"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 22, wantName: "Office PC",
		},
		{
			name: "unrestricted TCP risky port is danger",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 3389, Protocol: "tCp", RemoteHost: "  ", InternalPort: 3389, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Office PC"},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 3389, wantName: "Office PC",
		},
		{
			name: "unrestricted TCP ADB port is danger",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 5555, Protocol: "TCP", InternalPort: 5555, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Android TV"},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 5555, wantName: "Android TV",
		},
		{
			name: "unrestricted TCP RTSP port is camera danger",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 554, Protocol: "TCP", InternalPort: 554, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Living room camera"},
			wantID: "forward_camera", wantLevel: report.Danger, wantCount: 1, wantPort: 554, wantName: "Living room camera",
		},
		{
			name: "unrestricted TCP Dahua DVR port is camera danger",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 37777, Protocol: "TCP", InternalPort: 37777, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "DVR"},
			wantID: "forward_camera", wantLevel: report.Danger, wantCount: 1, wantPort: 37777, wantName: "DVR",
		},
		{
			name: "unrestricted TCP Xiongmai DVR port is camera danger",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 34567, Protocol: "TCP", InternalPort: 34567, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "DVR"},
			wantID: "forward_camera", wantLevel: report.Danger, wantCount: 1, wantPort: 34567, wantName: "DVR",
		},
		{
			name: "unrestricted TCP go2rtc RTSP port is camera danger",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 8554, Protocol: "TCP", InternalPort: 8554, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Home Assistant"},
			wantID: "forward_camera", wantLevel: report.Danger, wantCount: 1, wantPort: 8554, wantName: "Home Assistant",
		},
		{
			name: "unrestricted TCP go2rtc API port is risky",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 1984, Protocol: "TCP", InternalPort: 1984, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Home Assistant"},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 1984, wantName: "Home Assistant",
		},
		{
			name: "unrestricted TCP Node-RED editor port is risky",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 1880, Protocol: "TCP", InternalPort: 1880, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Node-RED"},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 1880, wantName: "Node-RED",
		},
		{
			name: "UDP Node-RED editor port is caution",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 1880, Protocol: "UDP", InternalPort: 1880, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Node-RED"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 1880, wantName: "Node-RED",
		},
		{
			name: "UDP RTSP port is not risky",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 554, Protocol: "UDP", InternalPort: 554, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Camera"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 554, wantName: "Camera",
		},
		{
			name: "restricted Dahua port is not risky",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 37777, Protocol: "TCP", RemoteHost: "203.0.113.8", InternalPort: 37777, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "DVR"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 37777, wantName: "DVR",
		},
		{
			name: "TCP 8000 is not risky",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 8000, Protocol: "TCP", InternalPort: 8000, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Camera"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 8000, wantName: "Camera",
		},
		{
			name: "ordinary forward",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 3074, InternalPort: 3074, InternalClient: internalIP.String(), Enabled: true,
			}},
			names:  map[netip.Addr]string{internalIP: "Game console"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 3074, wantName: "Game console",
		},
		{
			name: "TCP and UDP deduplicated",
			mappings: []ssdp.PortMapping{
				{ExternalPort: 3074, Protocol: "TCP", InternalPort: 3074, InternalClient: internalIP.String(), Enabled: true},
				{ExternalPort: 3074, Protocol: "UDP", InternalPort: 3074, InternalClient: internalIP.String(), Enabled: true},
			},
			names:  map[netip.Addr]string{internalIP: "Game console"},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 3074, wantName: "Game console",
		},
		{
			name: "UDP then TCP keeps higher classification",
			mappings: []ssdp.PortMapping{
				{ExternalPort: 4000, Protocol: "UDP", InternalPort: 4000, InternalClient: internalIP.String(), Enabled: true},
				{ExternalPort: 4000, Protocol: "TCP", InternalPort: 3389, InternalClient: internalIP.String(), Enabled: true},
			},
			names:  map[netip.Addr]string{internalIP: "Game console"},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 4000, wantName: "Game console",
		},
		{
			name: "upgraded group uses the risky mapping description",
			mappings: []ssdp.PortMapping{
				{ExternalPort: 4000, Protocol: "UDP", InternalPort: 4000, InternalClient: internalIP.String(), Enabled: true, Description: "Game"},
				{ExternalPort: 4000, Protocol: "TCP", InternalPort: 3389, InternalClient: internalIP.String(), Enabled: true, Description: "RDP"},
			},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 4000, wantName: "RDP",
		},
		{
			name: "TCP then UDP keeps higher classification",
			mappings: []ssdp.PortMapping{
				{ExternalPort: 4000, Protocol: "TCP", InternalPort: 3389, InternalClient: internalIP.String(), Enabled: true},
				{ExternalPort: 4000, Protocol: "UDP", InternalPort: 4000, InternalClient: internalIP.String(), Enabled: true},
			},
			names:  map[netip.Addr]string{internalIP: "Game console"},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1, wantPort: 4000, wantName: "Game console",
		},
		{
			name: "disabled mapping ignored",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 3389, InternalPort: 3389, InternalClient: internalIP.String(),
			}},
		},
		{
			name: "public internal client ignored",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 3389, InternalPort: 3389, InternalClient: publicIP.String(), Enabled: true,
			}},
		},
		{
			name: "description used when device name is unknown",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 3074, InternalPort: 3074, InternalClient: internalIP.String(), Enabled: true, Description: "  Game\x1b Console\u200f ",
			}},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1, wantPort: 3074, wantName: "Game Console",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := makeFindings(
				nil, nil, netip.Addr{}, map[netip.Addr]ssdp.Description{},
				test.names, map[netip.Addr]string{}, test.mappings, "en",
			)
			portFindings := make([]report.Finding, 0)
			for _, finding := range findings {
				if finding.ID == "port_forward" || finding.ID == "port_forward_risky" || finding.ID == "forward_camera" {
					portFindings = append(portFindings, finding)
				}
			}
			if len(portFindings) != test.wantCount {
				t.Fatalf("port forwarding findings = %#v, want %d", portFindings, test.wantCount)
			}
			if test.wantCount == 1 {
				finding := portFindings[0]
				if finding.ID != test.wantID || finding.Level != test.wantLevel || finding.Port != test.wantPort || finding.Device != test.wantName {
					t.Fatalf("port forwarding finding = %#v", finding)
				}
				wantPoints := 3
				if test.wantLevel == report.Danger {
					wantPoints = 20
				}
				if finding.Points != wantPoints || finding.IP != internalIP.String() {
					t.Fatalf("port forwarding finding = %#v", finding)
				}
				if finding.ID == "forward_camera" {
					wantMessage := i18n.T("en", "forward_camera", test.wantName, test.wantPort)
					wantFix := i18n.T("en", "forward_camera_fix")
					if finding.Message != wantMessage || finding.Fix != wantFix {
						t.Fatalf("camera forward text = %#v, want message %q and fix %q", finding, wantMessage, wantFix)
					}
				}
			}
		})
	}
}

func TestFetchPortMappingsQueryability(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	routerDescription := func(controlURL string) ssdp.Description {
		description := ssdp.Description{
			Device: ssdp.Device{
				DeviceType: "urn:schemas-upnp-org:device:InternetGatewayDevice:1",
			},
			Location: "http://192.168.1.1:1900/description.xml",
		}
		if controlURL != "" {
			description.Device.Services = []ssdp.Service{{
				ServiceType: "urn:schemas-upnp-org:service:WANIPConnection:1",
				ControlURL:  controlURL,
			}}
		}
		return description
	}

	t.Run("no description", func(t *testing.T) {
		mappings, queried, err := fetchPortMappings(context.Background(), nil, gateway, nil, nil)
		if err != nil || mappings != nil || queried {
			t.Fatalf("fetchPortMappings() = %#v, %t, %v; want nil, false, nil", mappings, queried, err)
		}
	})

	t.Run("non-IGD description without IGD responder", func(t *testing.T) {
		description := routerDescription("")
		description.Device.DeviceType = "urn:schemas-upnp-org:device:Basic:1"
		mappings, queried, err := fetchPortMappings(context.Background(), nil, gateway, nil, map[netip.Addr]ssdp.Description{
			gateway: description,
		})
		if err != nil || mappings != nil || queried {
			t.Fatalf("fetchPortMappings() = %#v, %t, %v; want nil, false, nil", mappings, queried, err)
		}
	})

	t.Run("nested IGD in basic gateway description", func(t *testing.T) {
		description := ssdp.Description{
			Device: ssdp.Device{
				DeviceType: "urn:schemas-upnp-org:device:Basic:1",
				Devices: []ssdp.Device{{
					DeviceType: "urn:schemas-upnp-org:device:InternetGatewayDevice:1",
					Services: []ssdp.Service{{
						ServiceType: "urn:schemas-upnp-org:service:WANIPConnection:1",
						ControlURL:  "http://8.8.8.8/control",
					}},
				}},
			},
			Location: "http://192.168.1.1:1900/description.xml",
		}
		_, queried, err := fetchPortMappings(context.Background(), nil, gateway, nil, map[netip.Addr]ssdp.Description{
			gateway: description,
		})
		if queried || err == nil || err.Error() != "router advertises UPnP but its port forwarding list cannot be queried" {
			t.Fatalf("fetchPortMappings() error = %v; want rejected control URL error", err)
		}
	})

	t.Run("IGD without WAN service", func(t *testing.T) {
		_, queried, err := fetchPortMappings(context.Background(), nil, gateway, nil, map[netip.Addr]ssdp.Description{
			gateway: routerDescription(""),
		})
		if queried || err == nil || err.Error() != "router advertises UPnP but its port forwarding list cannot be queried" {
			t.Fatalf("fetchPortMappings() error = %v", err)
		}
	})

	for _, controlURL := range []string{
		"http://8.8.8.8/control",
		"http://192.168.1.2/control",
	} {
		t.Run(controlURL, func(t *testing.T) {
			requests := 0
			client := &http.Client{Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
				requests++
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader("")),
					Request:    request,
				}, nil
			})}
			_, queried, err := fetchPortMappings(context.Background(), client, gateway, nil, map[netip.Addr]ssdp.Description{
				gateway: routerDescription(controlURL),
			})
			if queried || err == nil || err.Error() != "router advertises UPnP but its port forwarding list cannot be queried" {
				t.Fatalf("fetchPortMappings() error = %v", err)
			}
			if requests != 0 {
				t.Fatalf("HTTP requests = %d, want 0", requests)
			}
		})
	}

	t.Run("successful default gateway query", func(t *testing.T) {
		requests := 0
		client := &http.Client{Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
			requests++
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`<UPnPError><errorCode>713</errorCode></UPnPError>`)),
				Request:    request,
			}, nil
		})}
		mappings, queried, err := fetchPortMappings(context.Background(), client, gateway, nil, map[netip.Addr]ssdp.Description{
			gateway: routerDescription("/control"),
		})
		if err != nil || !queried || len(mappings) != 0 {
			t.Fatalf("fetchPortMappings() = %#v, %t, %v; want empty mappings, true, nil", mappings, queried, err)
		}
		if requests != 1 {
			t.Fatalf("HTTP requests = %d, want 1", requests)
		}
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (roundTripper roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTripper(request)
}

func TestFetchDescriptionsBoundsSpoofedFlood(t *testing.T) {
	previousBudget := descriptionFetchBudget
	descriptionFetchBudget = 200 * time.Millisecond
	t.Cleanup(func() {
		descriptionFetchBudget = previousBudget
	})

	gateway := netip.MustParseAddr("192.168.1.1")
	responses := make([]ssdp.Response, 0, 65)
	for host := 100; host <= 163; host++ {
		ip := netip.MustParseAddr(fmt.Sprintf("192.168.1.%d", host))
		responses = append(responses, ssdp.Response{
			IP:       ip,
			ST:       "upnp:rootdevice",
			Location: fmt.Sprintf("http://%s:49152/d.xml", ip),
		})
	}
	responses = append(responses, ssdp.Response{
		IP:       gateway,
		ST:       "upnp:rootdevice",
		Location: fmt.Sprintf("http://%s:49152/d.xml", gateway),
	})

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
			if request.URL.Hostname() == gateway.String() {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body: io.NopCloser(strings.NewReader(`<root><device>
<deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
<friendlyName>Home Router</friendlyName></device></root>`)),
				}, nil
			}
			<-request.Context().Done()
			return nil, request.Context().Err()
		}),
	}

	start := time.Now()
	devices, errs := fetchDescriptions(context.Background(), client, responses, gateway)
	elapsed := time.Since(start)
	if elapsed >= 2*time.Second {
		t.Fatalf("fetchDescriptions() took %s, want under 2s", elapsed)
	}
	if _, ok := devices[gateway]; !ok {
		t.Fatal("gateway description is missing")
	}
	if len(errs) > 9 {
		t.Fatalf("fetchDescriptions() returned %d errors, want at most 9", len(errs))
	}
	for _, err := range errs {
		if strings.Contains(err.Error(), "time limit reached") {
			return
		}
	}
	t.Fatal("fetchDescriptions() errors do not contain the time limit summary")
}

func TestRouterUPnPAttribution(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	secondary := netip.MustParseAddr("192.168.1.2")
	igd := ssdp.Response{IP: secondary, ST: "urn:schemas-upnp-org:device:InternetGatewayDevice:1"}
	tests := []struct {
		name         string
		gateway      netip.Addr
		responses    []ssdp.Response
		descriptions map[netip.Addr]ssdp.Description
		want         int
	}{
		{name: "secondary IGD ignored", gateway: gateway, responses: []ssdp.Response{igd}},
		{name: "known gateway response", gateway: gateway, responses: []ssdp.Response{{IP: gateway, ST: igd.ST}}, want: 1},
		{
			name:         "known gateway description",
			gateway:      gateway,
			descriptions: map[netip.Addr]ssdp.Description{gateway: {Device: ssdp.Device{DeviceType: "urn:schemas-upnp-org:device:InternetGatewayDevice:1"}}},
			want:         1,
		},
		{name: "unknown gateway preserves first IGD", responses: []ssdp.Response{igd}, want: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := makeFindings(nil, test.responses, test.gateway, test.descriptions, map[netip.Addr]string{}, map[netip.Addr]string{}, nil, "en")
			count := 0
			for _, finding := range findings {
				if finding.ID == "router_upnp" {
					count++
				}
			}
			if count != test.want {
				t.Fatalf("router_upnp findings = %d, want %d: %#v", count, test.want, findings)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	got := sanitizeName("\x1b[2J\x1b[31mCam")
	if got != "[2J[31mCam" {
		t.Fatalf("sanitizeName(escape sequences) = %q", got)
	}
	for _, r := range got {
		if unicode.IsControl(r) || unicode.Is(unicode.Bidi_Control, r) ||
			r >= 0x200B && r <= 0x200D || r == 0xFEFF {
			t.Errorf("sanitizeName left excluded rune U+%04X in %q", r, got)
		}
	}
	if got := sanitizeName(" \t Camera   One \u200b "); got != "Camera One" {
		t.Fatalf("sanitizeName whitespace = %q", got)
	}
	if got := sanitizeName(strings.Repeat("界", 70)); len([]rune(got)) != 64 {
		t.Fatalf("sanitizeName length = %d runes, want 64", len([]rune(got)))
	}

	ip := netip.MustParseAddr("192.168.1.20")
	name := deviceName("en", ip, lan.Neighbor{}, netip.Addr{}, map[netip.Addr]ssdp.Description{
		ip: {Device: ssdp.Device{Manufacturer: "Brand\x1b", ModelName: "Model\u200f"}},
	}, nil)
	if name != "Brand Model" {
		t.Fatalf("deviceName() = %q, want sanitized manufacturer and model", name)
	}
}

func TestStripUnsafeRunes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes controls and bidi overrides",
			input: "XML syntax error on line 1: invalid XML name: a\u009b8m\u202eevil\x1b",
			want:  "XML syntax error on line 1: invalid XML name: a8mevil",
		},
		{
			name:  "preserves spaces and Japanese",
			input: "家の device name と 通常のスペース",
			want:  "家の device name と 通常のスペース",
		},
		{
			name:  "removes Arabic letter mark",
			input: "Cam\u061Cera",
			want:  "Camera",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := stripUnsafeRunes(test.input); got != test.want {
				t.Fatalf("stripUnsafeRunes(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestPortMappingWarningSeparatesJoinedErrors(t *testing.T) {
	err := errors.Join(
		errors.New("WANIPConnection: HTTP 500"),
		errors.New("WANPPPConnection: HTTP 500"),
	)
	if got := portMappingWarning(err); !strings.Contains(got, "HTTP 500; WANPPPConnection:") {
		t.Fatalf("portMappingWarning() = %q, want semicolon between joined errors", got)
	}
}

func TestDeviceNamePrecedence(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.20")
	router := netip.MustParseAddr("192.168.1.1")
	randomized := lan.Neighbor{MAC: "02:11:22:33:44:55"}
	tests := []struct {
		name          string
		ip            netip.Addr
		neighbor      lan.Neighbor
		gateway       netip.Addr
		descriptions  map[netip.Addr]ssdp.Description
		resolvedNames map[netip.Addr]string
		want          string
	}{
		{
			name:     "SSDP description wins",
			ip:       router,
			neighbor: randomized,
			gateway:  router,
			descriptions: map[netip.Addr]ssdp.Description{
				router: {Device: ssdp.Device{FriendlyName: "Living Room"}},
			},
			resolvedNames: map[netip.Addr]string{router: "router.local"},
			want:          "Living Room",
		},
		{
			name:          "router label wins over resolved name",
			ip:            router,
			neighbor:      randomized,
			gateway:       router,
			resolvedNames: map[netip.Addr]string{router: "home-gateway"},
			want:          i18n.T("en", "router"),
		},
		{
			name:          "sanitized resolved name wins over randomized label",
			ip:            ip,
			neighbor:      randomized,
			resolvedNames: map[netip.Addr]string{ip: " Phone\u0000.local "},
			want:          "Phone.local",
		},
		{
			name:     "randomized label wins over IP",
			ip:       ip,
			neighbor: randomized,
			want:     i18n.T("en", "randomized_device"),
		},
		{
			name:          "uppercase NetBIOS name is preserved",
			ip:            ip,
			neighbor:      randomized,
			resolvedNames: map[netip.Addr]string{ip: "DESKTOP-AB12CD"},
			want:          "DESKTOP-AB12CD",
		},
		{name: "IP fallback", ip: ip, want: ip.String()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := deviceName(
				"en",
				test.ip,
				test.neighbor,
				test.gateway,
				test.descriptions,
				test.resolvedNames,
			)
			if got != test.want {
				t.Fatalf("deviceName() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestMergeDeviceNamesPrefersMDNSAndFallsBackToNetBIOS(t *testing.T) {
	mdnsIP := netip.MustParseAddr("192.168.1.20")
	netbiosIP := netip.MustParseAddr("192.168.1.21")
	emptyMDNSIP := netip.MustParseAddr("192.168.1.22")
	mdnsNames := map[netip.Addr]string{
		mdnsIP:      "iphone.local",
		emptyMDNSIP: "\x00",
	}
	netbiosNames := map[netip.Addr]string{
		mdnsIP:      "IPHONE",
		netbiosIP:   "DESKTOP-AB12CD",
		emptyMDNSIP: "PRINTER01",
	}

	got := mergeDeviceNames(mdnsNames, netbiosNames)
	want := map[netip.Addr]string{
		mdnsIP:      "iphone.local",
		netbiosIP:   "DESKTOP-AB12CD",
		emptyMDNSIP: "PRINTER01",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mergeDeviceNames() = %v, want %v", got, want)
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

func TestPreLanguage(t *testing.T) {
	getenv := func(key string) string {
		if key == "LANG" {
			return "fr_FR.UTF-8"
		}
		return ""
	}
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"short separate", []string{"-lang", "ja"}, "ja"},
		{"long separate", []string{"--lang", "ja"}, "ja"},
		{"short equals", []string{"-lang=ja"}, "ja"},
		{"long equals", []string{"--lang=ja"}, "ja"},
		{"language after boolean flag", []string{"--json", "--lang=ja"}, "ja"},
		{"language value can look like a flag", []string{"--state", "--", "--lang=ja"}, "ja"},
		{"state consumes language-looking value", []string{"--state", "--lang=ja"}, "fr"},
		{"help stops before language flag", []string{"--help", "--lang=ja"}, "fr"},
		{"positional argument stops before language flag", []string{"scan-pos", "--lang=ja"}, "fr"},
		{"lone dash stops before language flag", []string{"-", "--lang=ja"}, "fr"},
		{"unknown flag stops before language flag", []string{"--unknown", "--lang=ja"}, "fr"},
		{"unknown flag preserves earlier language", []string{"--lang=ja", "--unknown"}, "ja"},
		{"auto uses detected language", []string{"--lang=auto"}, "fr"},
		{"invalid uses detected language", []string{"--lang", "xx"}, "fr"},
		{"flag after separator is ignored", []string{"--", "--lang=ja"}, "fr"},
		{"missing value uses detected language", []string{"--lang"}, "fr"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := flag.NewFlagSet("mimamori", flag.ContinueOnError)
			flags.SetOutput(io.Discard)
			registeredFlags := registerRunFlags(flags)
			got := preLanguage(test.args, getenv, flags)

			_ = flags.Parse(test.args)
			parsed := i18n.Detect(getenv)
			if resolved, err := resolveLanguage(*registeredFlags.language, getenv); err == nil {
				parsed = resolved
			}
			if parsed != test.want {
				t.Fatalf("FlagSet.Parse(%v) selected %q, want %q", test.args, parsed, test.want)
			}
			if got != parsed {
				t.Fatalf("preLanguage(%v) = %q, want FlagSet language %q", test.args, got, parsed)
			}
		})
	}
}

func TestRunLocalizedHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
		env  string
		lang string
	}{
		{"detected Japanese", []string{"--help"}, "ja_JP.UTF-8", "ja"},
		{"explicit Korean", []string{"--lang=ko", "--help"}, "en_US.UTF-8", "ko"},
		{"help before language uses detected language", []string{"--help", "--lang=ja"}, "fr_FR.UTF-8", "fr"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			result := runWithEnv(test.args, &stdout, &stderr, func(key string) string {
				if key == "LANG" {
					return test.env
				}
				return ""
			})
			if result.code != 0 {
				t.Fatalf("runWithEnv() code = %d, want 0; stderr = %q", result.code, stderr.String())
			}
			if result.lang != test.lang {
				t.Fatalf("runWithEnv() language = %q, want %q", result.lang, test.lang)
			}
			for _, key := range []string{"usage_header", "flag_json", "flag_lang"} {
				var text string
				if key == "flag_lang" {
					text = i18n.T(test.lang, key, strings.Join(i18n.Supported(), ", "))
				} else {
					text = i18n.T(test.lang, key)
				}
				if !strings.Contains(stderr.String(), text) {
					t.Errorf("help output %q does not contain %q", stderr.String(), text)
				}
			}
		})
	}
}

func TestRunLocalizedStartupErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		env     string
		lang    string
		wantKey string
		wantArg string
	}{
		{"invalid language", []string{"--lang", "xx"}, "ja_JP.UTF-8", "ja", "invalid_lang", "xx"},
		{"unexpected argument", []string{"extra"}, "fr_FR.UTF-8", "fr", "unexpected_argument", "extra"},
		{"negative wait", []string{"--wait=-1s"}, "de_DE.UTF-8", "de", "wait_negative", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			result := runWithEnv(test.args, &stdout, &stderr, func(key string) string {
				if key == "LANG" {
					return test.env
				}
				return ""
			})
			if result.code != 1 {
				t.Fatalf("runWithEnv() code = %d, want 1; stderr = %q", result.code, stderr.String())
			}
			if result.lang != test.lang {
				t.Fatalf("runWithEnv() language = %q, want %q", result.lang, test.lang)
			}
			var want string
			if test.wantKey == "invalid_lang" {
				want = i18n.T(test.lang, test.wantKey, test.wantArg, strings.Join(i18n.Supported(), ", "))
			} else if test.wantArg != "" {
				want = i18n.T(test.lang, test.wantKey, test.wantArg)
			} else {
				want = i18n.T(test.lang, test.wantKey)
			}
			if got := stderr.String(); got != want+"\n" {
				t.Fatalf("runWithEnv() stderr = %q, want %q", got, want+"\n")
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

func TestPrintScanningNotice(t *testing.T) {
	const notice = "家のネットワークを確認しています。少し時間がかかります…"
	var textOutput bytes.Buffer
	printScanningNotice(&textOutput, "ja", false)
	if got, want := textOutput.String(), notice+"\n"; got != want {
		t.Fatalf("text scanning notice = %q, want %q", got, want)
	}

	var jsonOutput bytes.Buffer
	printScanningNotice(&jsonOutput, "ja", true)
	if jsonOutput.Len() != 0 {
		t.Fatalf("JSON scanning notice = %q, want empty", jsonOutput.String())
	}
}

func TestRunVersionOutputUnchanged(t *testing.T) {
	var stdout, stderr bytes.Buffer
	result := run([]string{"--version"}, &stdout, &stderr)
	if result.code != 0 {
		t.Fatalf("run(--version) code = %d, want 0", result.code)
	}
	if got, want := stdout.String(), "Mimamori "+version+"\n"; got != want {
		t.Fatalf("run(--version) stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("run(--version) stderr = %q, want empty", stderr.String())
	}
}

func TestCollectHostsFiltersAndCaps(t *testing.T) {
	neighbors := make([]lan.Neighbor, 1025)
	for i := range neighbors {
		neighbors[i].IP = netip.AddrFrom4([4]byte{10, 0, byte(i >> 8), byte(i)})
	}
	neighbors = append(neighbors, lan.Neighbor{IP: netip.MustParseAddr("8.8.8.8")})
	got := collectHosts(neighbors, netip.MustParseAddr("192.168.1.1"), []ssdp.Response{{IP: netip.MustParseAddr("172.16.0.1")}}, nil)
	if len(got) != 1024 {
		t.Fatalf("collectHosts() returned %d hosts", len(got))
	}
	for i := 1; i < len(got); i++ {
		if !got[i-1].Less(got[i]) {
			t.Fatalf("hosts are not sorted: %v before %v", got[i-1], got[i])
		}
	}
}

func TestCollectHostsSkipsOtherInterfaceSubnets(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	prefixes := []netip.Prefix{
		netip.MustParsePrefix("192.168.1.10/24"),
		netip.MustParsePrefix("10.0.0.5/24"),
	}
	scope := homePrefixes(prefixes, gateway)
	neighbors := []lan.Neighbor{
		{IP: netip.MustParseAddr("192.168.1.20")},
		{IP: netip.MustParseAddr("10.0.0.7")},
		{IP: netip.MustParseAddr("172.20.10.2")},
	}
	responses := []ssdp.Response{
		{IP: netip.MustParseAddr("10.0.0.9")},
		{IP: netip.MustParseAddr("192.168.1.30")},
	}
	got := collectHosts(neighbors, gateway, responses, scope)
	want := []netip.Addr{
		gateway,
		netip.MustParseAddr("192.168.1.20"),
		netip.MustParseAddr("192.168.1.30"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectHosts() = %v, want %v", got, want)
	}
}

func TestMergeHostsDeduplicatesSortsAndCaps(t *testing.T) {
	existing := []netip.Addr{
		netip.MustParseAddr("192.168.1.10"),
		netip.MustParseAddr("192.168.1.1"),
	}
	swept := []netip.Addr{
		netip.MustParseAddr("192.168.1.10"),
		netip.MustParseAddr("192.168.1.2"),
	}
	got := mergeHosts(existing, swept)
	want := []netip.Addr{
		netip.MustParseAddr("192.168.1.1"),
		netip.MustParseAddr("192.168.1.2"),
		netip.MustParseAddr("192.168.1.10"),
	}
	if len(got) != len(want) {
		t.Fatalf("mergeHosts() returned %d hosts, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("mergeHosts()[%d] = %s, want %s", i, got[i], want[i])
		}
	}

	full := make([]netip.Addr, maxScanHosts)
	for i := range full {
		full[i] = netip.AddrFrom4([4]byte{10, 0, byte(i >> 8), byte(i)})
	}
	extra := netip.MustParseAddr("10.0.4.1")
	capped := mergeHosts(full, []netip.Addr{extra})
	if len(capped) != maxScanHosts {
		t.Fatalf("mergeHosts() returned %d hosts, want cap of %d", len(capped), maxScanHosts)
	}
	if capped[len(capped)-1] == extra {
		t.Fatal("mergeHosts() exceeded the cap with a swept host")
	}
}

func TestMergeNeighborsOriginalWinsAndRefreshAdds(t *testing.T) {
	ip1 := netip.MustParseAddr("192.168.1.1")
	ip2 := netip.MustParseAddr("192.168.1.2")
	ip3 := netip.MustParseAddr("192.168.1.3")
	got := mergeNeighbors(
		[]lan.Neighbor{{IP: ip1, MAC: "aa:aa:aa:aa:aa:01"}, {IP: ip2, MAC: "aa:aa:aa:aa:aa:02"}},
		[]lan.Neighbor{{IP: ip1, MAC: "bb:bb:bb:bb:bb:01"}, {IP: ip3, MAC: "bb:bb:bb:bb:bb:03"}},
	)
	want := []lan.Neighbor{
		{IP: ip1, MAC: "bb:bb:bb:bb:bb:01"},
		{IP: ip2, MAC: "aa:aa:aa:aa:aa:02"},
		{IP: ip3, MAC: "bb:bb:bb:bb:bb:03"},
	}
	if len(got) != len(want) {
		t.Fatalf("mergeNeighbors() returned %d entries, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("mergeNeighbors()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestDeviceHostsIncludesOnlyObservedSweptHosts(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	neighbors := []lan.Neighbor{{IP: gateway, MAC: "aa:bb:cc:dd:ee:ff"}}
	scanTargets := mergeHosts(
		collectHosts(neighbors, gateway, nil, nil),
		lan.SweepHosts([]netip.Prefix{netip.MustParsePrefix("192.168.1.23/24")}, gateway),
	)
	if len(scanTargets) != 253 {
		t.Fatalf("scan target count = %d, want 253", len(scanTargets))
	}

	got := deviceHosts(neighbors, gateway, nil, []probe.Open{{
		IP: netip.MustParseAddr("192.168.1.30"), Port: 23,
	}})
	want := []netip.Addr{gateway, netip.MustParseAddr("192.168.1.30")}
	if len(got) != len(want) {
		t.Fatalf("deviceHosts() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("deviceHosts()[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestDeviceHostsKeepsOpenedHostsBeyondScanCap(t *testing.T) {
	neighbors := make([]lan.Neighbor, 0, 1000)
	for i := 0; i < 1000; i++ {
		ip := netip.AddrFrom4([4]byte{10, 0, byte(i / 250), byte(i%250 + 1)})
		neighbors = append(neighbors, lan.Neighbor{IP: ip, MAC: "aa:bb:cc:dd:ee:ff"})
	}
	opened := make([]probe.Open, 0, 100)
	for i := 0; i < 100; i++ {
		opened = append(opened, probe.Open{IP: netip.AddrFrom4([4]byte{10, 1, 0, byte(i + 1)}), Port: 21})
	}
	got := deviceHosts(neighbors, netip.Addr{}, nil, opened)
	if len(got) != len(neighbors)+len(opened) {
		t.Fatalf("deviceHosts() returned %d hosts, want %d", len(got), len(neighbors)+len(opened))
	}
	present := make(map[netip.Addr]bool, len(got))
	for i, host := range got {
		if i > 0 && !got[i-1].Less(host) {
			t.Fatalf("deviceHosts() is not sorted and unique at %s, %s", got[i-1], host)
		}
		present[host] = true
	}
	for _, neighbor := range neighbors {
		if !present[neighbor.IP] {
			t.Fatalf("observed neighbor %s was dropped", neighbor.IP)
		}
	}
	for _, item := range opened {
		if !present[item.IP] {
			t.Fatalf("opened host %s was dropped", item.IP)
		}
	}
}

func TestDeviceHostsKeepsOpenedHostsBeyondARPCap(t *testing.T) {
	neighbors := make([]lan.Neighbor, 0, maxScanHosts)
	for i := 0; i < maxScanHosts; i++ {
		ip := netip.AddrFrom4([4]byte{10, 0, byte(i / 254), byte(i%254 + 1)})
		neighbors = append(neighbors, lan.Neighbor{IP: ip, MAC: "aa:bb:cc:dd:ee:ff"})
	}
	openedOnly := netip.MustParseAddr("192.168.1.30")
	opened := []probe.Open{
		{IP: neighbors[0].IP, Port: 23},
		{IP: openedOnly, Port: 23},
	}

	got := deviceHosts(neighbors, netip.Addr{}, nil, opened)
	if len(got) != maxScanHosts+1 {
		t.Fatalf("deviceHosts() returned %d hosts, want %d", len(got), maxScanHosts+1)
	}
	present := make(map[netip.Addr]bool, len(got))
	for i, host := range got {
		if i > 0 && !got[i-1].Less(host) {
			t.Fatalf("deviceHosts() is not sorted and unique at %s, %s", got[i-1], host)
		}
		present[host] = true
	}
	if !present[openedOnly] {
		t.Fatalf("opened host %s is missing", openedOnly)
	}
}

func TestGatewayString(t *testing.T) {
	if got := gatewayString(netip.Addr{}); got != "" {
		t.Fatalf("gatewayString(invalid) = %q", got)
	}
}

func TestScanObservedNetwork(t *testing.T) {
	tests := []struct {
		name      string
		gateway   netip.Addr
		neighbors []lan.Neighbor
		responses []ssdp.Response
		opened    []probe.Open
		want      bool
	}{
		{name: "gateway only, nothing open", gateway: netip.MustParseAddr("192.168.1.1")},
		{name: "neighbor observed", neighbors: []lan.Neighbor{{IP: netip.MustParseAddr("192.168.1.2")}}, want: true},
		{name: "SSDP response observed", responses: []ssdp.Response{{IP: netip.MustParseAddr("192.168.1.2")}}, want: true},
		{name: "open service observed", opened: []probe.Open{{IP: netip.MustParseAddr("192.168.1.2"), Port: 23}}, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.gateway.IsValid() {
				if hosts := collectHosts(test.neighbors, test.gateway, test.responses, nil); len(hosts) != 1 {
					t.Fatalf("gateway-only scan hosts = %v, want only the gateway", hosts)
				}
			}
			if got := scanObservedNetwork(test.neighbors, test.responses, test.opened); got != test.want {
				t.Fatalf("scanObservedNetwork() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestReportable(t *testing.T) {
	tests := []struct {
		name                string
		observed            bool
		wifiFindings        []report.Finding
		wantReport          bool
		wantForceIncomplete bool
	}{
		{name: "observed network", observed: true, wantReport: true},
		{name: "unobserved with no Wi-Fi finding", wifiFindings: wifiFindings(wifi.Modern, "en")},
		{
			name:                "Wi-Fi danger allows unobserved report",
			wifiFindings:        wifiFindings(wifi.Open, "en"),
			wantReport:          true,
			wantForceIncomplete: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotReport, gotForceIncomplete := reportable(test.observed, test.wifiFindings)
			if gotReport != test.wantReport || gotForceIncomplete != test.wantForceIncomplete {
				t.Fatalf("reportable() = (%v, %v), want (%v, %v)",
					gotReport, gotForceIncomplete, test.wantReport, test.wantForceIncomplete)
			}
		})
	}
}

func TestShouldSaveHistory(t *testing.T) {
	populated := baseline.Store{Devices: map[string]baseline.Seen{"aa:bb:cc:dd:ee:ff": {}}}
	tests := []struct {
		name            string
		outputSucceeded bool
		firstRun        bool
		store           baseline.Store
		loadState       historyLoadState
		want            bool
	}{
		{name: "successful output and existing history", outputSucceeded: true, store: baseline.Store{}, loadState: historyLoadReady, want: true},
		{name: "output failure", outputSucceeded: false, store: populated, loadState: historyLoadReady},
		{name: "empty first run", outputSucceeded: true, firstRun: true, store: baseline.Store{Devices: map[string]baseline.Seen{}}, loadState: historyLoadReady},
		{name: "populated first run", outputSucceeded: true, firstRun: true, store: populated, loadState: historyLoadReady, want: true},
		{name: "quarantined history", outputSucceeded: true, firstRun: true, store: populated, loadState: historyLoadQuarantined, want: true},
		{name: "unavailable history", outputSucceeded: true, firstRun: true, store: populated, loadState: historyLoadUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := shouldSaveHistory(test.outputSucceeded, test.firstRun, test.store, test.loadState)
			if got != test.want {
				t.Fatalf("shouldSaveHistory() = %v, want %v", got, test.want)
			}
		})
	}
}

type failingOutputWriter struct{}

func (failingOutputWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

type shortOutputWriter struct{}

func (shortOutputWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) - 1, nil
}

func TestWriteOutputFailurePreventsHistorySave(t *testing.T) {
	result := output{Score: 100, Status: "safe"}
	store := baseline.Store{Devices: map[string]baseline.Seen{"aa:bb:cc:dd:ee:ff": {}}}
	if err := writeOutput(failingOutputWriter{}, result, false, "en"); err == nil {
		t.Fatal("writeOutput() succeeded with a failing writer")
	}
	if shouldSaveHistory(false, false, store, historyLoadReady) {
		t.Fatal("history should not be saved after a failed output write")
	}
	if err := writeOutput(shortOutputWriter{}, result, true, "en"); err != io.ErrShortWrite {
		t.Fatalf("short write error = %v, want %v", err, io.ErrShortWrite)
	}
}

func TestRunIncompleteWarningRouting(t *testing.T) {
	t.Run("text output uses localized notice only", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if err := writeOutput(&stdout, output{Score: 100, Status: "safe", Incomplete: true}, false, "ja"); err != nil {
			t.Fatal(err)
		}
		finishRun(&stderr, 0, "ja", false, true, true)
		if !strings.Contains(stdout.String(), i18n.T("ja", "incomplete_notice")) {
			t.Fatalf("text output = %q, missing localized incomplete notice", stdout.String())
		}
		if strings.Contains(stderr.String(), "warning: some checks could not run; this result may be incomplete") {
			t.Fatalf("text stderr = %q, unexpectedly includes generic incomplete warning", stderr.String())
		}
	})

	t.Run("JSON output retains generic warning", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if err := writeOutput(&stdout, output{Incomplete: true}, true, "en"); err != nil {
			t.Fatal(err)
		}
		finishRun(&stderr, 0, "en", true, true, false)
		if !strings.Contains(stderr.String(), "warning: some checks could not run; this result may be incomplete") {
			t.Fatalf("JSON stderr = %q, missing generic incomplete warning", stderr.String())
		}
	})

	t.Run("early exit retains generic warning", func(t *testing.T) {
		var stderr bytes.Buffer
		finishRun(&stderr, 1, "ja", false, true, false)
		if !strings.Contains(stderr.String(), "warning: some checks could not run; this result may be incomplete") {
			t.Fatalf("early-exit stderr = %q, missing generic incomplete warning", stderr.String())
		}
	})
}

func TestRunOutputFailureRetainsIncompleteWarning(t *testing.T) {
	var stderr bytes.Buffer
	err := writeOutput(failingOutputWriter{}, output{Incomplete: true}, false, "ja")
	if err == nil {
		t.Fatal("writeOutput() succeeded with a failing writer")
	}
	finishRun(&stderr, 1, "ja", false, true, false)
	if !strings.Contains(stderr.String(), "warning: some checks could not run; this result may be incomplete") {
		t.Fatalf("failed-output stderr = %q, missing generic incomplete warning", stderr.String())
	}
}

func TestWriteOutputIncludesIncompleteJSON(t *testing.T) {
	var buffer bytes.Buffer
	if err := writeOutput(&buffer, output{Incomplete: true}, true, "en"); err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["incomplete"] != true {
		t.Fatalf("JSON incomplete = %#v, want true", result["incomplete"])
	}
}

func TestWriteOutputHistoryUnavailableText(t *testing.T) {
	var buffer bytes.Buffer
	result := output{FirstRun: true, HistoryUnavailable: true}
	if err := writeOutput(&buffer, result, false, "ja"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buffer.String(), i18n.T("ja", "history_unavailable")) {
		t.Fatalf("text output = %q, missing history-unavailable message", buffer.String())
	}
	if strings.Contains(buffer.String(), i18n.T("ja", "first_run")) {
		t.Fatalf("text output = %q, includes first-run message", buffer.String())
	}
}

func TestWriteOutputFirstRunText(t *testing.T) {
	var buffer bytes.Buffer
	if err := writeOutput(&buffer, output{FirstRun: true}, false, "ja"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buffer.String(), i18n.T("ja", "first_run")) {
		t.Fatalf("text output = %q, missing first-run message", buffer.String())
	}
}

func TestWriteOutputHistoryUnavailableJSON(t *testing.T) {
	for _, test := range []struct {
		name        string
		unavailable bool
		wantPresent bool
	}{
		{name: "true is included", unavailable: true, wantPresent: true},
		{name: "false is omitted", unavailable: false, wantPresent: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			result := output{HistoryUnavailable: test.unavailable}
			if err := writeOutput(&buffer, result, true, "ja"); err != nil {
				t.Fatal(err)
			}
			var encoded map[string]any
			if err := json.Unmarshal(buffer.Bytes(), &encoded); err != nil {
				t.Fatal(err)
			}
			value, present := encoded["history_unavailable"]
			if present != test.wantPresent {
				t.Fatalf("JSON key presence = %t, want %t", present, test.wantPresent)
			}
			if test.wantPresent && value != true {
				t.Fatalf("JSON history_unavailable = %#v, want true", value)
			}
		})
	}
}

func TestRenderTextCheckedLine(t *testing.T) {
	tests := []struct {
		name   string
		lang   string
		result output
		want   string
	}{
		{
			name: "all checks",
			lang: "en",
			result: output{
				Score: 100, Status: "safe", Devices: make([]deviceOutput, 12),
				Checks: checkedSummary{Ports: true, PortForwards: true, WiFi: true},
			},
			want: "Checked: 12 devices · risky ports · router UPnP port forwards · Wi-Fi encryption",
		},
		{
			name: "devices and risky ports only",
			lang: "en",
			result: output{
				Score: 100, Status: "safe", Devices: []deviceOutput{{}},
				Checks: checkedSummary{Ports: true},
			},
			want: "Checked: 1 device · risky ports",
		},
		{
			name: "Japanese",
			lang: "ja",
			result: output{
				Score: 100, Status: "safe", Devices: make([]deviceOutput, 12),
				Checks: checkedSummary{Ports: true, PortForwards: true, WiFi: true},
			},
			want: "確認したこと: 機器 12台・危険なポート・ルーターの UPnP ポート転送・Wi-Fi の暗号化",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			renderText(&buffer, test.result, test.lang)
			lines := strings.Split(strings.TrimSuffix(buffer.String(), "\n"), "\n")
			if len(lines) < 3 {
				t.Fatalf("renderText() output has too few lines: %q", buffer.String())
			}
			if lines[2] != test.want {
				t.Fatalf("checked line = %q, want %q", lines[2], test.want)
			}
		})
	}
}

func TestRenderTextSkipsPortsWhenNoTargets(t *testing.T) {
	device := deviceOutput{IP: "192.168.1.20"}
	result := output{
		Score: 100, Status: "safe", Devices: []deviceOutput{device},
		Checks: checkedSummary{WiFi: true},
	}
	var buffer bytes.Buffer
	renderText(&buffer, result, "en")
	text := buffer.String()

	if strings.Contains(text, i18n.T("en", "checked_ports")) {
		t.Fatalf("renderText() claims a port check with no scan targets: %q", text)
	}
	if !strings.Contains(text, i18n.T("en", "checked_devices", i18n.CheckedDeviceCount("en", 1))) ||
		!strings.Contains(text, i18n.T("en", "checked_wifi")) ||
		!strings.Contains(text, device.IP) {
		t.Fatalf("renderText() omitted devices or Wi-Fi check: %q", text)
	}
}

func TestWriteOutputOmitsCheckedLineFromJSON(t *testing.T) {
	var buffer bytes.Buffer
	result := output{
		Checks:  checkedSummary{Ports: true},
		Devices: []deviceOutput{{}},
	}
	if err := writeOutput(&buffer, result, true, "en"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buffer.String(), "Checked:") {
		t.Fatalf("JSON output contains the human checked line: %s", buffer.String())
	}
	var encoded map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &encoded); err != nil {
		t.Fatal(err)
	}
	checked, ok := encoded["checked"].(map[string]any)
	if !ok {
		t.Fatalf("JSON checked = %#v, want object", encoded["checked"])
	}
	if len(checked) != 3 {
		t.Fatalf("JSON checked has %d keys, want exactly 3: %#v", len(checked), checked)
	}
	for _, field := range []struct {
		name string
		want bool
	}{
		{name: "ports", want: true},
		{name: "port_forwards", want: false},
		{name: "wifi", want: false},
	} {
		value, present := checked[field.name]
		if !present {
			t.Errorf("JSON checked is missing %q", field.name)
			continue
		}
		if value != field.want {
			t.Errorf("JSON checked.%s = %#v, want %t", field.name, value, field.want)
		}
	}
}

func TestRenderTextShowsLocalizedIncompleteNotice(t *testing.T) {
	scoreLine := i18n.T("ja", "score", i18n.T("ja", "safe"), 100)
	checkedLine := i18n.T("ja", "checked_devices", i18n.CheckedDeviceCount("ja", 0)) +
		i18n.T("ja", "checked_sep") + i18n.T("ja", "checked_ports")
	notice := i18n.T("ja", "incomplete_notice")
	for _, test := range []struct {
		name       string
		incomplete bool
	}{
		{name: "incomplete", incomplete: true},
		{name: "complete", incomplete: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			renderText(&buffer, output{
				Score: 100, Status: "safe", Incomplete: test.incomplete,
				Checks: checkedSummary{Ports: true},
			}, "ja")
			lines := strings.Split(strings.TrimSuffix(buffer.String(), "\n"), "\n")
			if len(lines) < 3 {
				t.Fatalf("renderText() output has too few lines: %q", buffer.String())
			}
			if lines[1] != scoreLine {
				t.Fatalf("score line = %q, want %q", lines[1], scoreLine)
			}
			if lines[2] != checkedLine {
				t.Fatalf("line after score = %q, want checked summary %q", lines[2], checkedLine)
			}
			if test.incomplete {
				if len(lines) < 4 {
					t.Fatalf("renderText() output has too few lines: %q", buffer.String())
				}
				if lines[3] != notice {
					t.Fatalf("line after checked summary = %q, want notice %q", lines[3], notice)
				}
				return
			}
			for _, line := range lines {
				if line == notice {
					t.Fatal("renderText() included incomplete notice for a complete result")
				}
			}
		})
	}
}

func TestFindingsCatalogHitForAllSupportedLanguages(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	for _, lang := range i18n.Supported() {
		findings := makeFindings(
			[]probe.Open{{IP: ip, Port: 23}, {IP: ip, Port: 21}, {IP: ip, Port: 5900}},
			nil,
			netip.Addr{},
			map[netip.Addr]ssdp.Description{},
			map[netip.Addr]string{ip: "Router"},
			map[netip.Addr]string{},
			[]ssdp.PortMapping{{
				ExternalPort: 44338, InternalPort: 3389, InternalClient: ip.String(), Enabled: true,
			}},
			lang,
		)
		for _, finding := range findings {
			if finding.Fix == finding.ID+"_fix" || finding.Message == finding.ID {
				t.Errorf("%s finding fell back to a catalog key: %#v", lang, finding)
			}
		}
	}
}

func TestWiFiFindings(t *testing.T) {
	tests := []struct {
		name     string
		security wifi.Security
		id       string
		level    report.Level
		points   int
	}{
		{name: "open", security: wifi.Open, id: "wifi_open", level: report.Danger, points: 25},
		{name: "WEP", security: wifi.WEP, id: "wifi_wep", level: report.Danger, points: 25},
		{name: "WPA1 or TKIP", security: wifi.WPA1orTKIP, id: "wifi_weak", level: report.Caution, points: 8},
		{name: "unknown", security: wifi.Unknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := wifiFindings(test.security, "en")
			if test.id == "" {
				if len(findings) != 0 {
					t.Fatalf("wifiFindings() = %#v, want no finding", findings)
				}
				return
			}
			if len(findings) != 1 {
				t.Fatalf("wifiFindings() returned %d findings, want 1", len(findings))
			}
			finding := findings[0]
			if finding.ID != test.id || finding.Level != test.level || finding.Points != test.points {
				t.Fatalf("wifi finding = %#v, want %s/%s/%d", finding, test.id, test.level, test.points)
			}
			if finding.IP != "" || finding.Device != "" {
				t.Errorf("Wi-Fi finding includes a device address/name: %#v", finding)
			}
			if finding.Message == "" || finding.Fix == "" || finding.Message == finding.ID || finding.Fix == finding.ID+"_fix" {
				t.Errorf("Wi-Fi finding is missing localized text: %#v", finding)
			}
		})
	}
}

func TestWiFiLocationRequiredTextAndJSON(t *testing.T) {
	scoreLine := i18n.T("ja", "score", i18n.T("ja", "safe"), 100)
	checkedLine := i18n.T("ja", "checked_devices", i18n.CheckedDeviceCount("ja", 0)) +
		i18n.T("ja", "checked_sep") + i18n.T("ja", "checked_ports")
	incompleteNotice := i18n.T("ja", "incomplete_notice")
	locationNotice := i18n.T("ja", "wifi_location_notice")

	var textOutput bytes.Buffer
	renderText(&textOutput, output{
		Score: 100, Status: "safe", Incomplete: true, WiFiLocationRequired: true,
		Checks: checkedSummary{Ports: true},
	}, "ja")
	lines := strings.Split(strings.TrimSuffix(textOutput.String(), "\n"), "\n")
	if len(lines) < 5 {
		t.Fatalf("renderText() output has too few lines: %q", textOutput.String())
	}
	if lines[1] != scoreLine || lines[2] != checkedLine || lines[3] != incompleteNotice || lines[4] != locationNotice {
		t.Fatalf("notice order = %q, %q, %q, %q; want score, checked summary, incomplete notice, location notice", lines[1], lines[2], lines[3], lines[4])
	}

	var noLocationNotice bytes.Buffer
	renderText(&noLocationNotice, output{
		Score: 100, Status: "safe", Incomplete: true, WiFiLocationRequired: false,
	}, "ja")
	if strings.Contains(noLocationNotice.String(), locationNotice) {
		t.Fatal("renderText() included the location notice when WiFiLocationRequired is false")
	}

	for _, test := range []struct {
		name     string
		required bool
		present  bool
	}{
		{name: "false is omitted", required: false, present: false},
		{name: "true is included", required: true, present: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var jsonOutput bytes.Buffer
			if err := writeOutput(&jsonOutput, output{WiFiLocationRequired: test.required}, true, "ja"); err != nil {
				t.Fatal(err)
			}
			var result map[string]any
			if err := json.Unmarshal(jsonOutput.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			value, exists := result["wifi_location_required"]
			if exists != test.present {
				t.Fatalf("JSON key presence = %t, want %t", exists, test.present)
			}
			if test.present && value != true {
				t.Fatalf("JSON wifi_location_required = %#v, want true", value)
			}
		})
	}
}

func TestPrintNoDeviceError(t *testing.T) {
	const errorLine = "error: the scan could not see any device on your network; the result would not be reliable"
	locationNotice := i18n.T("ja", "wifi_location_notice")
	tests := []struct {
		name       string
		wifiErr    error
		wantNotice bool
	}{
		{name: "location permission", wifiErr: wifi.ErrLocationPermission, wantNotice: true},
		{name: "wrapped location permission", wifiErr: fmt.Errorf("x: %w", wifi.ErrLocationPermission), wantNotice: true},
		{name: "nil error"},
		{name: "other error", wifiErr: context.Canceled},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			printNoDeviceError(&output, "ja", test.wifiErr)
			want := errorLine + "\n"
			if test.wantNotice {
				want += locationNotice + "\n"
			}
			if got := output.String(); got != want {
				t.Fatalf("printNoDeviceError() = %q, want %q", got, want)
			}
		})
	}
}

func TestFetchDescriptionsPrefersGatewayDescription(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	secondary := netip.MustParseAddr("192.168.1.50")
	response := func(ip netip.Addr, st string, port string) ssdp.Response {
		return ssdp.Response{
			IP:       ip,
			ST:       st,
			Location: fmt.Sprintf("http://%s:%s/rootDesc.xml", ip, port),
		}
	}
	igdST := "urn:schemas-upnp-org:device:InternetGatewayDevice:1"
	tests := []struct {
		name           string
		ip             netip.Addr
		responses      []ssdp.Response
		wantLocation   string
		wantGatewayIGD bool
	}{
		{
			name: "gateway IGD beats earlier media root device",
			ip:   gateway,
			responses: []ssdp.Response{
				response(gateway, "upnp:rootdevice", "8200"),
				response(gateway, "upnp:rootdevice", "5000"),
				response(gateway, igdST, "5000"),
			},
			wantLocation:   "http://192.168.1.1:5000/rootDesc.xml",
			wantGatewayIGD: true,
		},
		{
			name: "UPnP IGD beats DSL Forum IGD",
			ip:   gateway,
			responses: []ssdp.Response{
				response(gateway, "urn:dslforum-org:device:InternetGatewayDevice:1", "49443"),
				response(gateway, igdST, "5000"),
			},
			wantLocation:   "http://192.168.1.1:5000/rootDesc.xml",
			wantGatewayIGD: true,
		},
		{
			name: "UPnP WAN service reply beats DSL Forum IGD without IGD device reply",
			ip:   gateway,
			responses: []ssdp.Response{
				response(gateway, "urn:dslforum-org:device:InternetGatewayDevice:1", "49443"),
				response(gateway, "upnp:rootdevice", "5000"),
				response(gateway, "urn:schemas-upnp-org:service:WANIPConnection:1", "5000"),
			},
			wantLocation:   "http://192.168.1.1:5000/rootDesc.xml",
			wantGatewayIGD: true,
		},
		{
			name: "UPnP MediaServer does not beat root device",
			ip:   gateway,
			responses: []ssdp.Response{
				response(gateway, "urn:schemas-upnp-org:device:MediaServer:1", "8200"),
				response(gateway, "upnp:rootdevice", "5000"),
			},
			wantLocation:   "http://192.168.1.1:5000/rootDesc.xml",
			wantGatewayIGD: true,
		},
		{
			name: "root device beats UUID response for non-gateway",
			ip:   secondary,
			responses: []ssdp.Response{
				response(secondary, "uuid:device-1", "8200"),
				response(secondary, "upnp:rootdevice", "5000"),
			},
			wantLocation: "http://192.168.1.50:5000/rootDesc.xml",
		},
	}

	mediaDescription := `<root><device>
<deviceType>urn:schemas-upnp-org:device:MediaServer:1</deviceType>
<friendlyName>Media Server</friendlyName>
</device></root>`
	igdDescription := `<root><device>
<deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
<deviceList><device><deviceType>urn:schemas-upnp-org:device:WANDevice:1</deviceType>
<deviceList><device><deviceType>urn:schemas-upnp-org:device:WANConnectionDevice:1</deviceType>
<serviceList><service>
<serviceType>urn:schemas-upnp-org:service:WANIPConnection:1</serviceType>
<controlURL>/ctl/IPConn</controlURL>
</service></serviceList>
</device></deviceList></device></deviceList>
</device></root>`
	dslForumDescription := `<root><device>
<deviceType>urn:dslforum-org:device:InternetGatewayDevice:1</deviceType>
<friendlyName>DSL Forum Device</friendlyName>
</device></root>`

	client := &http.Client{Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		var body string
		switch request.URL.Port() {
		case "8200":
			body = mediaDescription
		case "5000":
			body = igdDescription
		case "49443":
			body = dslForumDescription
		default:
			t.Errorf("unexpected description URL %q", request.URL)
			body = mediaDescription
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    request,
		}, nil
	})}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			devices, errs := fetchDescriptions(context.Background(), client, test.responses, gateway)
			description, ok := devices[test.ip]
			if !ok {
				t.Fatalf("fetchDescriptions() has no device for %s; errors: %v", test.ip, errs)
			}
			if description.Location != test.wantLocation {
				t.Fatalf("description location = %q, want %q", description.Location, test.wantLocation)
			}
			if test.wantGatewayIGD && (!description.Device.IsIGD() || len(description.Device.WANConnections()) != 1) {
				t.Fatalf("gateway description = %#v, want IGD with one WAN connection", description.Device)
			}
		})
	}
}
