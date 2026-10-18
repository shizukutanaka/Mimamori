package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
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
	"github.com/shizukutanaka/Mimamori/internal/wifi"
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
				map[netip.Addr]ssdp.Description{},
				map[netip.Addr]string{ip: "Router", newIP: "Camera"},
				map[netip.Addr]bool{newIP: true},
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
				test.names, map[netip.Addr]bool{}, test.mappings, "en",
			)
			portFindings := make([]report.Finding, 0)
			for _, finding := range findings {
				if finding.ID == "port_forward" || finding.ID == "port_forward_risky" {
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
		mappings, err := fetchPortMappings(context.Background(), nil, gateway, nil, nil)
		if err != nil || mappings != nil {
			t.Fatalf("fetchPortMappings() = %#v, %v; want nil, nil", mappings, err)
		}
	})

	t.Run("non-IGD description without IGD responder", func(t *testing.T) {
		description := routerDescription("")
		description.Device.DeviceType = "urn:schemas-upnp-org:device:Basic:1"
		mappings, err := fetchPortMappings(context.Background(), nil, gateway, nil, map[netip.Addr]ssdp.Description{
			gateway: description,
		})
		if err != nil || mappings != nil {
			t.Fatalf("fetchPortMappings() = %#v, %v; want nil, nil", mappings, err)
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
		_, err := fetchPortMappings(context.Background(), nil, gateway, nil, map[netip.Addr]ssdp.Description{
			gateway: description,
		})
		if err == nil || err.Error() != "router advertises UPnP but its port forwarding list cannot be queried" {
			t.Fatalf("fetchPortMappings() error = %v; want rejected control URL error", err)
		}
	})

	t.Run("IGD without WAN service", func(t *testing.T) {
		_, err := fetchPortMappings(context.Background(), nil, gateway, nil, map[netip.Addr]ssdp.Description{
			gateway: routerDescription(""),
		})
		if err == nil || err.Error() != "router advertises UPnP but its port forwarding list cannot be queried" {
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
			_, err := fetchPortMappings(context.Background(), client, gateway, nil, map[netip.Addr]ssdp.Description{
				gateway: routerDescription(controlURL),
			})
			if err == nil || err.Error() != "router advertises UPnP but its port forwarding list cannot be queried" {
				t.Fatalf("fetchPortMappings() error = %v", err)
			}
			if requests != 0 {
				t.Fatalf("HTTP requests = %d, want 0", requests)
			}
		})
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (roundTripper roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTripper(request)
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
			findings := makeFindings(nil, test.responses, test.gateway, test.descriptions, map[netip.Addr]string{}, map[netip.Addr]bool{}, nil, "en")
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
	name := deviceName("en", ip, lan.Neighbor{}, netip.Addr{}, map[netip.Addr]ssdp.Description{
		ip: {Device: ssdp.Device{Manufacturer: "Brand\x1b", ModelName: "Model\u200f"}},
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
		{IP: ip1, MAC: "aa:aa:aa:aa:aa:01"},
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
		collectHosts(neighbors, gateway, nil),
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

func TestRenderTextShowsLocalizedIncompleteNotice(t *testing.T) {
	scoreLine := i18n.T("ja", "score", i18n.T("ja", "safe"), 100)
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
			renderText(&buffer, output{Score: 100, Status: "safe", Incomplete: test.incomplete}, "ja")
			lines := strings.Split(strings.TrimSuffix(buffer.String(), "\n"), "\n")
			if len(lines) < 3 {
				t.Fatalf("renderText() output has too few lines: %q", buffer.String())
			}
			if lines[1] != scoreLine {
				t.Fatalf("score line = %q, want %q", lines[1], scoreLine)
			}
			if test.incomplete {
				if lines[2] != notice {
					t.Fatalf("line after score = %q, want notice %q", lines[2], notice)
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
			map[netip.Addr]bool{},
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
	incompleteNotice := i18n.T("ja", "incomplete_notice")
	locationNotice := i18n.T("ja", "wifi_location_notice")

	var textOutput bytes.Buffer
	renderText(&textOutput, output{
		Score: 100, Status: "safe", Incomplete: true, WiFiLocationRequired: true,
	}, "ja")
	lines := strings.Split(strings.TrimSuffix(textOutput.String(), "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("renderText() output has too few lines: %q", textOutput.String())
	}
	if lines[1] != scoreLine || lines[2] != incompleteNotice || lines[3] != locationNotice {
		t.Fatalf("notice order = %q, %q, %q; want score, incomplete notice, location notice", lines[1], lines[2], lines[3])
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
