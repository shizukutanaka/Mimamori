package main

import (
	"net/netip"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func TestMQTTPortForwardFindings(t *testing.T) {
	client := netip.MustParseAddr("192.168.1.30")
	tests := []struct {
		name       string
		mappings   []ssdp.PortMapping
		wantID     string
		wantLevel  report.Level
		wantPoints int
	}{
		{
			name: "unrestricted TCP MQTT",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 1883, Protocol: "TCP", InternalPort: 1883,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantPoints: 20,
		},
		{
			name: "TCP MQTT on remapped external port",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 41883, Protocol: "TCP", InternalPort: 1883,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward_risky", wantLevel: report.Danger, wantPoints: 20,
		},
		{
			name: "UDP MQTT remains plain forward",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 1883, Protocol: "UDP", InternalPort: 1883,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward", wantLevel: report.Caution, wantPoints: 3,
		},
		{
			name: "restricted TCP MQTT remote host remains plain forward",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 1883, Protocol: "TCP", RemoteHost: "203.0.113.8",
				InternalPort: 1883, InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward", wantLevel: report.Caution, wantPoints: 3,
		},
		{
			name: "TCP MQTT over TLS remains plain forward",
			mappings: []ssdp.PortMapping{{
				ExternalPort: 8883, Protocol: "TCP", InternalPort: 8883,
				InternalClient: client.String(), Enabled: true,
			}},
			wantID: "port_forward", wantLevel: report.Caution, wantPoints: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := makeFindings(
				nil, nil, netip.Addr{}, nil,
				map[netip.Addr]string{client: "Home hub"}, nil, test.mappings, "en",
			)
			if len(findings) != 1 {
				t.Fatalf("makeFindings() returned %d findings, want 1: %#v", len(findings), findings)
			}
			finding := findings[0]
			if finding.ID != test.wantID || finding.Level != test.wantLevel || finding.Points != test.wantPoints {
				t.Fatalf(
					"finding = %#v, want ID %q, level %q, points %d",
					finding, test.wantID, test.wantLevel, test.wantPoints,
				)
			}
		})
	}
}
