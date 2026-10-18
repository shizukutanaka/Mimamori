package main

import (
	"net/netip"
	"reflect"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/lan"
)

func TestHomeNeighbors(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	tests := []struct {
		name      string
		neighbors []lan.Neighbor
		prefixes  []netip.Prefix
		gateway   netip.Addr
		want      []lan.Neighbor
	}{
		{
			name: "keeps only the gateway home prefix",
			neighbors: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.1"), MAC: "aa:bb:cc:00:00:01"},
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
				{IP: netip.MustParseAddr("172.20.50.3"), MAC: "aa:bb:cc:00:00:03"},
				{IP: netip.MustParseAddr("192.168.56.101"), MAC: "aa:bb:cc:00:00:04"},
			},
			prefixes: []netip.Prefix{
				netip.MustParsePrefix("192.168.1.0/24"),
				netip.MustParsePrefix("172.20.48.0/20"),
			},
			gateway: gateway,
			want: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.1"), MAC: "aa:bb:cc:00:00:01"},
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
			},
		},
		{
			name: "gateway in no prefix leaves neighbors unchanged",
			neighbors: []lan.Neighbor{
				{IP: netip.MustParseAddr("172.20.50.3"), MAC: "aa:bb:cc:00:00:03"},
				{IP: netip.MustParseAddr("192.168.56.101"), MAC: "aa:bb:cc:00:00:04"},
			},
			prefixes: []netip.Prefix{
				netip.MustParsePrefix("192.168.2.0/24"),
				netip.MustParsePrefix("172.20.48.0/20"),
			},
			gateway: gateway,
			want: []lan.Neighbor{
				{IP: netip.MustParseAddr("172.20.50.3"), MAC: "aa:bb:cc:00:00:03"},
				{IP: netip.MustParseAddr("192.168.56.101"), MAC: "aa:bb:cc:00:00:04"},
			},
		},
		{
			name: "invalid gateway leaves neighbors unchanged",
			neighbors: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
			},
			prefixes: []netip.Prefix{netip.MustParsePrefix("192.168.1.0/24")},
			want: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
			},
		},
		{
			name: "nil prefixes leave neighbors unchanged",
			neighbors: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
			},
			gateway: gateway,
			want: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
			},
		},
		{
			name: "union of matching prefixes",
			neighbors: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
				{IP: netip.MustParseAddr("192.168.1.200"), MAC: "aa:bb:cc:00:00:05"},
				{IP: netip.MustParseAddr("172.20.50.3"), MAC: "aa:bb:cc:00:00:03"},
			},
			prefixes: []netip.Prefix{
				netip.MustParsePrefix("192.168.1.0/25"),
				netip.MustParsePrefix("192.168.1.0/24"),
				netip.MustParsePrefix("172.20.48.0/20"),
			},
			gateway: gateway,
			want: []lan.Neighbor{
				{IP: netip.MustParseAddr("192.168.1.20"), MAC: "aa:bb:cc:00:00:02"},
				{IP: netip.MustParseAddr("192.168.1.200"), MAC: "aa:bb:cc:00:00:05"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := homeNeighbors(test.neighbors, test.prefixes, test.gateway)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("homeNeighbors() = %#v, want %#v", got, test.want)
			}
		})
	}
}
