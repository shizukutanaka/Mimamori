package lan

import (
	"net/netip"
	"testing"
)

func TestSweepHosts(t *testing.T) {
	tests := []struct {
		name        string
		prefixes    []netip.Prefix
		gateway     netip.Addr
		wantCount   int
		wantInclude []netip.Addr
		wantExclude []netip.Addr
	}{
		{
			name:      "/24 excludes network, broadcast, and self",
			prefixes:  []netip.Prefix{netip.MustParsePrefix("192.168.1.23/24")},
			gateway:   netip.MustParseAddr("192.168.1.1"),
			wantCount: 253,
			wantInclude: []netip.Addr{
				netip.MustParseAddr("192.168.1.1"),
				netip.MustParseAddr("192.168.1.254"),
			},
			wantExclude: []netip.Addr{
				netip.MustParseAddr("192.168.1.0"),
				netip.MustParseAddr("192.168.1.23"),
				netip.MustParseAddr("192.168.1.255"),
			},
		},
		{
			name:      "/16 is capped to own /24",
			prefixes:  []netip.Prefix{netip.MustParsePrefix("10.0.5.7/16")},
			gateway:   netip.MustParseAddr("10.0.5.1"),
			wantCount: 253,
			wantInclude: []netip.Addr{
				netip.MustParseAddr("10.0.5.1"),
				netip.MustParseAddr("10.0.5.254"),
			},
			wantExclude: []netip.Addr{netip.MustParseAddr("10.0.5.7")},
		},
		{
			name:      "/16 with gateway outside own /24 is skipped",
			prefixes:  []netip.Prefix{netip.MustParsePrefix("10.0.5.7/16")},
			gateway:   netip.MustParseAddr("10.0.0.1"),
			wantCount: 0,
		},
		{
			name:        "/30 includes the other usable host",
			prefixes:    []netip.Prefix{netip.MustParsePrefix("192.168.1.1/30")},
			gateway:     netip.MustParseAddr("192.168.1.2"),
			wantCount:   1,
			wantInclude: []netip.Addr{netip.MustParseAddr("192.168.1.2")},
		},
		{
			name:      "/31 is skipped",
			prefixes:  []netip.Prefix{netip.MustParsePrefix("192.168.1.1/31")},
			gateway:   netip.MustParseAddr("192.168.1.0"),
			wantCount: 0,
		},
		{
			name:      "public prefix is skipped",
			prefixes:  []netip.Prefix{netip.MustParsePrefix("8.8.8.10/24")},
			gateway:   netip.MustParseAddr("8.8.8.1"),
			wantCount: 0,
		},
		{
			name:      "prefix not containing gateway is skipped",
			prefixes:  []netip.Prefix{netip.MustParsePrefix("172.17.0.1/16")},
			gateway:   netip.MustParseAddr("192.168.1.1"),
			wantCount: 0,
		},
		{
			name: "same subnet on two interfaces is deduplicated",
			prefixes: []netip.Prefix{
				netip.MustParsePrefix("192.168.1.23/24"),
				netip.MustParsePrefix("192.168.1.24/24"),
			},
			gateway:     netip.MustParseAddr("192.168.1.1"),
			wantCount:   252,
			wantInclude: []netip.Addr{netip.MustParseAddr("192.168.1.1")},
			wantExclude: []netip.Addr{
				netip.MustParseAddr("192.168.1.23"),
				netip.MustParseAddr("192.168.1.24"),
			},
		},
		{
			name:      "invalid gateway",
			prefixes:  []netip.Prefix{netip.MustParsePrefix("192.168.1.23/24")},
			wantCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := SweepHosts(test.prefixes, test.gateway)
			if len(got) != test.wantCount {
				t.Fatalf("SweepHosts() returned %d addresses, want %d", len(got), test.wantCount)
			}
			for i := 1; i < len(got); i++ {
				if !got[i-1].Less(got[i]) {
					t.Fatalf("SweepHosts() is not sorted: %v before %v", got[i-1], got[i])
				}
			}
			seen := make(map[netip.Addr]struct{}, len(got))
			for _, address := range got {
				seen[address] = struct{}{}
			}
			for _, address := range test.wantInclude {
				if _, ok := seen[address]; !ok {
					t.Errorf("SweepHosts() does not include %s", address)
				}
			}
			for _, address := range test.wantExclude {
				if _, ok := seen[address]; ok {
					t.Errorf("SweepHosts() unexpectedly includes %s", address)
				}
			}
		})
	}
}
