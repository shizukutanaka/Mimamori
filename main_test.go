package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/netip"
	"strings"
	"testing"
	"unicode"

	"github.com/shizukutanaka/Mimamori/internal/baseline"
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
				ip,
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

func TestRouterUPnPAttribution(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	secondary := netip.MustParseAddr("192.168.1.2")
	igd := ssdp.Response{IP: secondary, ST: "urn:schemas-upnp-org:device:InternetGatewayDevice:1"}
	tests := []struct {
		name         string
		gateway      netip.Addr
		responses    []ssdp.Response
		descriptions map[netip.Addr]ssdp.Device
		want         int
	}{
		{name: "secondary IGD ignored", gateway: gateway, responses: []ssdp.Response{igd}},
		{name: "known gateway response", gateway: gateway, responses: []ssdp.Response{{IP: gateway, ST: igd.ST}}, want: 1},
		{
			name:         "known gateway description",
			gateway:      gateway,
			descriptions: map[netip.Addr]ssdp.Device{gateway: {DeviceType: "urn:schemas-upnp-org:device:InternetGatewayDevice:1"}},
			want:         1,
		},
		{name: "unknown gateway preserves first IGD", responses: []ssdp.Response{igd}, want: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := makeFindings(nil, test.responses, test.gateway, test.descriptions, map[netip.Addr]string{}, map[netip.Addr]bool{}, "en")
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
		if unicode.IsControl(r) || r >= 0x200B && r <= 0x200F ||
			r >= 0x202A && r <= 0x202E || r >= 0x2066 && r <= 0x2069 || r == 0xFEFF {
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
	name := deviceName("en", ip, lan.Neighbor{}, netip.Addr{}, map[netip.Addr]ssdp.Device{
		ip: {Manufacturer: "Brand\x1b", ModelName: "Model\u200f"},
	})
	if name != "Brand Model" {
		t.Fatalf("deviceName() = %q, want sanitized manufacturer and model", name)
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
				if hosts := collectHosts(test.neighbors, test.gateway, test.responses); len(hosts) != 1 {
					t.Fatalf("gateway-only scan hosts = %v, want only the gateway", hosts)
				}
			}
			if got := scanObservedNetwork(test.neighbors, test.responses, test.opened); got != test.want {
				t.Fatalf("scanObservedNetwork() = %v, want %v", got, test.want)
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

func TestFindingsCatalogHitForAllSupportedLanguages(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	for _, lang := range i18n.Supported() {
		findings := makeFindings(
			[]probe.Open{{IP: ip, Port: 23}, {IP: ip, Port: 21}, {IP: ip, Port: 5900}},
			nil,
			netip.Addr{},
			map[netip.Addr]ssdp.Device{},
			map[netip.Addr]string{ip: "Router"},
			map[netip.Addr]bool{},
			lang,
		)
		for _, finding := range findings {
			if finding.Fix == finding.ID+"_fix" || finding.Message == finding.ID {
				t.Errorf("%s finding fell back to a catalog key: %#v", lang, finding)
			}
		}
	}
}
