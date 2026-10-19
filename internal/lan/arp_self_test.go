package lan

import (
	"net/netip"
	"testing"
)

func TestWithoutAddresses(t *testing.T) {
	first := Neighbor{IP: netip.MustParseAddr("172.16.4.1"), MAC: "02:44:45:56:00:00"}
	self := Neighbor{IP: netip.MustParseAddr("172.16.4.2"), MAC: "06:00:ac:10:04:02"}
	last := Neighbor{IP: netip.MustParseAddr("172.16.4.3"), MAC: "02:44:45:56:00:03"}
	neighbors := []Neighbor{first, self, last}

	tests := []struct {
		name     string
		excluded map[netip.Addr]struct{}
		want     []Neighbor
	}{
		{name: "nil exclusions", want: neighbors},
		{
			name: "remove excluded address and preserve order",
			excluded: map[netip.Addr]struct{}{
				self.IP: {},
			},
			want: []Neighbor{first, last},
		},
		{
			name: "remove all addresses",
			excluded: map[netip.Addr]struct{}{
				first.IP: {},
				self.IP:  {},
				last.IP:  {},
			},
			want: []Neighbor{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := withoutAddresses(neighbors, test.excluded)
			if len(got) != len(test.want) {
				t.Fatalf("withoutAddresses() returned %d neighbors, want %d", len(got), len(test.want))
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("withoutAddresses()[%d] = %#v, want %#v", i, got[i], test.want[i])
				}
			}
			if test.excluded == nil && &got[0] != &neighbors[0] {
				t.Fatal("withoutAddresses() copied the input with no exclusions")
			}
		})
	}
}

func TestWithoutAddressesRemovesMacOSSelfEntry(t *testing.T) {
	arpOutput := `? (172.16.4.2) at 6:0:ac:10:4:2 on en0 ifscope permanent [ethernet]
? (172.16.4.1) at 02:44:45:56:00:00 on en0 ifscope [ethernet]`
	excluded := map[netip.Addr]struct{}{
		netip.MustParseAddr("172.16.4.2"): {},
	}

	got := withoutAddresses(ParseARP(arpOutput), excluded)
	want := Neighbor{IP: netip.MustParseAddr("172.16.4.1"), MAC: "02:44:45:56:00:00"}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("filtered ARP neighbors = %#v, want only %#v", got, want)
	}
}

func TestLocalAddressesIncludesLoopback(t *testing.T) {
	if _, ok := LocalAddresses()[netip.MustParseAddr("127.0.0.1")]; !ok {
		t.Fatal("LocalAddresses() does not contain 127.0.0.1")
	}
}
