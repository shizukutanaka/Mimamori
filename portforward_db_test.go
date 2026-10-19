package main

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func TestDatabasePortForwardFindings(t *testing.T) {
	client := netip.MustParseAddr("192.168.1.30")
	device := "Home NAS"
	tests := []struct {
		name       string
		mappings   []ssdp.PortMapping
		wantID     string
		wantLevel  report.Level
		wantCount  int
		wantDevice bool
	}{
		{
			name: "unrestricted TCP database",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 27017, Protocol: "TCP", InternalPort: 27017,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "forward_database", wantLevel: report.Danger, wantCount: 1, wantDevice: true,
		},
		{
			name: "UDP database remains plain forward",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 6379, Protocol: "UDP", InternalPort: 6379,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1,
		},
		{
			name: "restricted database remote host remains plain forward",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 3306, Protocol: "TCP", RemoteHost: "203.0.113.5",
				InternalPort: 3306, InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward", wantLevel: report.Caution, wantCount: 1,
		},
		{
			name: "NetBIOS is risky",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 139, Protocol: "TCP", InternalPort: 139,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1,
		},
		{
			name: "Docker API is risky",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 2375, Protocol: "TCP", InternalPort: 2375,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantCount: 1,
		},
		{
			name: "database classification wins during deduplication",
			mappings: []ssdp.PortMapping{
				{
					ExternalPort: 4444, Protocol: "TCP", InternalPort: 8080,
					InternalClient: client.String(), Enabled: true,
				},
				{
					ExternalPort: 4444, Protocol: "TCP", InternalPort: 5432,
					InternalClient: client.String(), Enabled: true,
				},
			},
			wantID: "forward_database", wantLevel: report.Danger, wantCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lang := "en"
			if test.wantDevice {
				lang = "ja"
			}
			findings := makeFindings(
				nil, nil, netip.Addr{}, nil,
				map[netip.Addr]string{client: device}, nil, test.mappings, lang,
			)
			if len(findings) != test.wantCount {
				t.Fatalf("makeFindings() returned %d findings, want %d: %#v", len(findings), test.wantCount, findings)
			}
			finding := findings[0]
			if finding.ID != test.wantID || finding.Level != test.wantLevel {
				t.Fatalf("finding = %#v, want ID %q at level %q", finding, test.wantID, test.wantLevel)
			}
			if test.wantDevice && (!strings.Contains(finding.Message, device) || !strings.Contains(finding.Message, "27017")) {
				t.Fatalf("database message = %q, want device %q and port 27017", finding.Message, device)
			}
		})
	}
}
