package lan

import (
	"net/netip"
	"testing"
)

func TestParseARPFixtures(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []Neighbor
	}{
		{
			name: "Windows English",
			text: `Interface: 192.168.1.10 --- 0x5
  Internet Address      Physical Address      Type
  192.168.1.1           aa-bb-cc-dd-ee-ff     dynamic
  192.168.1.2           0-1a-2b-3-c-5         dynamic`,
			want: []Neighbor{
				{IP: netip.MustParseAddr("192.168.1.1"), MAC: "aa:bb:cc:dd:ee:ff"},
				{IP: netip.MustParseAddr("192.168.1.2"), MAC: "00:1a:2b:03:0c:05"},
			},
		},
		{
			name: "Windows Japanese",
			text: `インターフェイス: 192.168.1.10 --- 0x5
インターネット アドレス 物理アドレス 種類
192.168.1.1 aa-bb-cc-dd-ee-ff 動的
192.168.1.2 00-1a-2b-3c-4d-5e 静的
192.168.1.255 ff-ff-ff-ff-ff-ff 静的
224.0.0.22 01-00-5e-00-00-16 静的`,
			want: []Neighbor{
				{IP: netip.MustParseAddr("192.168.1.1"), MAC: "aa:bb:cc:dd:ee:ff"},
				{IP: netip.MustParseAddr("192.168.1.2"), MAC: "00:1a:2b:3c:4d:5e"},
			},
		},
		{
			name: "macOS",
			text: `? (192.168.1.1) at 0:1a:2b:3c:4d:5e on en0 ifscope [ethernet]
? (192.168.1.2) at (incomplete) on en0 ifscope [ethernet]`,
			want: []Neighbor{
				{IP: netip.MustParseAddr("192.168.1.1"), MAC: "00:1a:2b:3c:4d:5e"},
			},
		},
		{
			name: "Linux",
			text: `IP address       HW type     Flags       HW address            Mask     Device
192.168.1.1       0x1         0x2         aa:bb:cc:dd:ee:ff     *        eth0
192.168.1.2       0x1         0x0         00:00:00:00:00:00     *        eth0`,
			want: []Neighbor{
				{IP: netip.MustParseAddr("192.168.1.1"), MAC: "aa:bb:cc:dd:ee:ff"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ParseARP(test.text)
			if len(got) != len(test.want) {
				t.Fatalf("ParseARP() = %#v, want %#v", got, test.want)
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("ParseARP()[%d] = %#v, want %#v", i, got[i], test.want[i])
				}
			}
		})
	}
}

func TestParseARPSkipsPublicAndDeduplicates(t *testing.T) {
	text := `8.8.8.8 aa:bb:cc:dd:ee:ff
192.168.1.1 aa:bb:cc:dd:ee:ff
192.168.1.1 11:22:33:44:55:66`
	got := ParseARP(text)
	if len(got) != 1 || got[0].MAC != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("ParseARP() = %#v", got)
	}
}

func TestIsRandomizedMAC(t *testing.T) {
	if !IsRandomizedMAC("02:11:22:33:44:55") || IsRandomizedMAC("00:11:22:33:44:55") {
		t.Fatal("unexpected randomized MAC detection")
	}
}

func TestRouteParsers(t *testing.T) {
	t.Run("linux", func(t *testing.T) {
		text := "Iface Destination Gateway Flags RefCnt Use Metric Mask\n" +
			"eth0 00000000 0201A8C0 0003 0 0 600 00000000\n" +
			"eth0 00000000 0101A8C0 0003 0 0 100 00000000\n" +
			"eth1 00000000 0301A8C0 0001 0 0 1 00000000\n"
		got, ok := ParseLinuxRoute(text)
		if !ok || got.String() != "192.168.1.1" {
			t.Fatalf("ParseLinuxRoute() = %v, %v", got, ok)
		}
	})
	t.Run("BSD", func(t *testing.T) {
		got, ok := ParseBSDRoute("   route to: default\n destination: default\n    gateway: 192.168.1.1\n")
		if !ok || got.String() != "192.168.1.1" {
			t.Fatalf("ParseBSDRoute() = %v, %v", got, ok)
		}
		if _, ok := ParseBSDRoute("gateway: router.local\n"); ok {
			t.Fatal("hostname gateway accepted")
		}
	})
	t.Run("Windows", func(t *testing.T) {
		text := "Network Destination        Netmask          Gateway        Interface       Metric\n" +
			"0.0.0.0          0.0.0.0      On-link        192.168.1.10    1\n" +
			"0.0.0.0          0.0.0.0      192.168.1.2    192.168.1.10    600\n" +
			"0.0.0.0          0.0.0.0      192.168.1.1    192.168.1.10    100\n"
		got, ok := ParseWindowsRoute(text)
		if !ok || got.String() != "192.168.1.1" {
			t.Fatalf("ParseWindowsRoute() = %v, %v", got, ok)
		}
	})
}
