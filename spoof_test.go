package main

import (
	"context"
	"errors"
	"net/netip"
	"reflect"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/lan"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func TestGatewaySpoofFindings(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	sharedMAC := "aa:bb:cc:00:11:22"
	tests := []struct {
		name      string
		gateway   netip.Addr
		neighbors []lan.Neighbor
		wantIP    string
	}{
		{
			name:    "spoofed MAC normalizes case and dashes",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: sharedMAC},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: "AA-BB-CC-00-11-22"},
			},
			wantIP: "192.168.1.50",
		},
		{
			name:    "multiple duplicates select the lowest IP",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: sharedMAC},
				{IP: netip.MustParseAddr("192.168.1.60"), MAC: sharedMAC},
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: sharedMAC},
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "AA-BB-CC-00-11-22"},
			},
			wantIP: "192.168.1.20",
		},
		{
			name:    "unique MACs",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: sharedMAC},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: "aa:bb:cc:00:11:23"},
			},
		},
		{
			name:      "gateway missing from neighbors",
			gateway:   gateway,
			neighbors: []lan.Neighbor{{IP: netip.MustParseAddr("192.168.1.50"), MAC: sharedMAC}},
		},
		{
			name: "invalid gateway",
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: sharedMAC},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: sharedMAC},
			},
		},
		{
			name:    "broadcast MAC",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: "ff:ff:ff:ff:ff:ff"},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: "ff:ff:ff:ff:ff:ff"},
			},
		},
		{
			name:    "multicast MAC",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: "01:00:5e:00:00:fb"},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: "01:00:5e:00:00:fb"},
			},
		},
		{
			name:    "all-zero MAC",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: "00:00:00:00:00:00"},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: "00:00:00:00:00:00"},
			},
		},
		{
			name:    "malformed gateway MAC",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: "aa:bb:cc"},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: "aa:bb:cc"},
			},
		},
		{
			name:    "public address is ignored",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: sharedMAC},
				{IP: netip.MustParseAddr("8.8.8.8"), MAC: sharedMAC},
			},
		},
		{
			name:    "unpadded macOS octets match",
			gateway: gateway,
			neighbors: []lan.Neighbor{
				{IP: gateway, MAC: "aa:b:c:0:11:22"},
				{IP: netip.MustParseAddr("192.168.1.50"), MAC: "AA:0b:0c:00:11:22"},
			},
			wantIP: "192.168.1.50",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := gatewaySpoofFindings(test.neighbors, test.gateway, "en")
			if test.wantIP == "" {
				if findings != nil {
					t.Fatalf("gatewaySpoofFindings() = %#v, want nil", findings)
				}
				return
			}
			if len(findings) != 1 {
				t.Fatalf("gatewaySpoofFindings() returned %d findings, want one: %#v", len(findings), findings)
			}
			finding := findings[0]
			if finding.ID != "gateway_spoof" || finding.Level != report.Caution || finding.Points != 8 {
				t.Fatalf("finding = %#v, want gateway_spoof/Caution/8", finding)
			}
			if finding.IP != test.wantIP {
				t.Fatalf("finding IP = %q, want %q", finding.IP, test.wantIP)
			}
			if finding.Message != i18n.T("en", "gateway_spoof", test.wantIP) {
				t.Fatalf("finding message = %q, want localized message", finding.Message)
			}
			if finding.Fix != i18n.T("en", "gateway_spoof_fix") {
				t.Fatalf("finding fix = %q, want localized fix", finding.Fix)
			}
		})
	}
}

func TestGatewaySpoofCatalogCoverage(t *testing.T) {
	englishMessage := i18n.T("en", "gateway_spoof", "192.168.1.50")
	englishFix := i18n.T("en", "gateway_spoof_fix")
	for _, lang := range i18n.Supported() {
		message := i18n.T(lang, "gateway_spoof", "192.168.1.50")
		fix := i18n.T(lang, "gateway_spoof_fix")
		if message == "" || message == "gateway_spoof" {
			t.Errorf("%s catalog is missing gateway_spoof", lang)
		}
		if fix == "" || fix == "gateway_spoof_fix" {
			t.Errorf("%s catalog is missing gateway_spoof_fix", lang)
		}
		if lang != "en" && (message == englishMessage || fix == englishFix) {
			t.Errorf("%s catalog fell back to English for a gateway spoof message", lang)
		}
	}
}

func TestGatewayNeighborsRefreshFindsSpoof(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	initial := []lan.Neighbor{{IP: netip.MustParseAddr("192.168.1.10"), MAC: "aa:bb:cc:00:11:33"}}
	refreshed := []lan.Neighbor{
		{IP: gateway, MAC: "aa:bb:cc:00:11:22"},
		{IP: netip.MustParseAddr("192.168.1.50"), MAC: "aa:bb:cc:00:11:22"},
	}
	calls := 0

	neighbors := gatewayNeighbors(context.Background(), initial, gateway, func(context.Context) ([]lan.Neighbor, error) {
		calls++
		return refreshed, nil
	})
	findings := gatewaySpoofFindings(neighbors, gateway, "en")

	if calls != 1 {
		t.Fatalf("readARP calls = %d, want 1", calls)
	}
	if len(findings) != 1 || findings[0].ID != "gateway_spoof" {
		t.Fatalf("gatewaySpoofFindings() = %#v, want one gateway_spoof finding", findings)
	}
}

func TestGatewayNeighborsSkipsRefreshWhenGatewayPresent(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	neighbors := []lan.Neighbor{{IP: gateway, MAC: "aa:bb:cc:00:11:22"}}
	calls := 0

	got := gatewayNeighbors(context.Background(), neighbors, gateway, func(context.Context) ([]lan.Neighbor, error) {
		calls++
		return nil, nil
	})

	if calls != 0 {
		t.Fatalf("readARP calls = %d, want 0", calls)
	}
	if !reflect.DeepEqual(got, neighbors) {
		t.Fatalf("gatewayNeighbors() = %#v, want unchanged neighbors %#v", got, neighbors)
	}
}

func TestGatewayNeighborsKeepsNeighborsOnRefreshError(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	neighbors := []lan.Neighbor{{IP: netip.MustParseAddr("192.168.1.10"), MAC: "aa:bb:cc:00:11:33"}}
	calls := 0

	got := gatewayNeighbors(context.Background(), neighbors, gateway, func(context.Context) ([]lan.Neighbor, error) {
		calls++
		return nil, errors.New("read failed")
	})

	if calls != 1 {
		t.Fatalf("readARP calls = %d, want 1", calls)
	}
	if !reflect.DeepEqual(got, neighbors) {
		t.Fatalf("gatewayNeighbors() = %#v, want unchanged neighbors %#v", got, neighbors)
	}
}

func TestGatewayNeighborsSkipsRefreshForInvalidGateway(t *testing.T) {
	neighbors := []lan.Neighbor{{IP: netip.MustParseAddr("192.168.1.10"), MAC: "aa:bb:cc:00:11:33"}}
	calls := 0

	got := gatewayNeighbors(context.Background(), neighbors, netip.Addr{}, func(context.Context) ([]lan.Neighbor, error) {
		calls++
		return nil, nil
	})

	if calls != 0 {
		t.Fatalf("readARP calls = %d, want 0", calls)
	}
	if !reflect.DeepEqual(got, neighbors) {
		t.Fatalf("gatewayNeighbors() = %#v, want unchanged neighbors %#v", got, neighbors)
	}
}
