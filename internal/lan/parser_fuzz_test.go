package lan

import (
	"net"
	"strings"
	"testing"
)

func FuzzParseARP(f *testing.F) {
	seeds := []string{
		`Interface: 192.168.1.10 --- 0x5
  Internet Address      Physical Address      Type
  192.168.1.1           aa-bb-cc-dd-ee-ff     dynamic
  192.168.1.2           0-1a-2b-3-c-5         dynamic`,
		`インターフェイス: 192.168.1.10 --- 0x5
インターネット アドレス 物理アドレス 種類
192.168.1.1 aa-bb-cc-dd-ee-ff 動的
192.168.1.2 00-1a-2b-3c-4d-5e 静的
192.168.1.255 ff-ff-ff-ff-ff-ff 静的
224.0.0.22 01-00-5e-00-00-16 静的`,
		`? (192.168.1.1) at 0:1a:2b:3c:4d:5e on en0 ifscope [ethernet]
? (192.168.1.2) at (incomplete) on en0 ifscope [ethernet]`,
		`IP address       HW type     Flags       HW address            Mask     Device
192.168.1.1       0x1         0x2         aa:bb:cc:dd:ee:ff     *        eth0
192.168.1.2       0x1         0x0         00:00:00:00:00:00     *        eth0`,
		`8.8.8.8 aa:bb:cc:dd:ee:ff
192.168.1.1 aa:bb:cc:dd:ee:ff
192.168.1.1 11:22:33:44:55:66`,
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		for _, neighbor := range ParseARP(text) {
			if !neighbor.IP.Is4() || !neighbor.IP.IsPrivate() {
				t.Fatalf("ParseARP() returned non-private IPv4 address %v", neighbor.IP)
			}
			mac, err := net.ParseMAC(neighbor.MAC)
			if err != nil || len(mac) != 6 || mac.String() != neighbor.MAC {
				t.Fatalf("ParseARP() returned non-normalized MAC %q", neighbor.MAC)
			}
		}
	})
}

func FuzzParseLinuxRoute(f *testing.F) {
	seeds := []string{
		"Iface Destination Gateway Flags RefCnt Use Metric Mask\n" +
			"eth0 00000000 0201A8C0 0003 0 0 600 00000000\n" +
			"eth0 00000000 0101A8C0 0003 0 0 100 00000000\n" +
			"eth1 00000000 0301A8C0 0001 0 0 1 00000000\n",
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		gateway, ok := ParseLinuxRoute(text)
		if ok && !gateway.Is4() {
			t.Fatalf("ParseLinuxRoute() returned non-IPv4 address %v", gateway)
		}
	})
}

func FuzzParseBSDRoute(f *testing.F) {
	for _, seed := range []string{
		"   route to: default\n destination: default\n    gateway: 192.168.1.1\n",
		"gateway: router.local\n",
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		gateway, ok := ParseBSDRoute(text)
		if ok && !gateway.Is4() {
			t.Fatalf("ParseBSDRoute() returned non-IPv4 address %v", gateway)
		}
	})
}

func FuzzParseWindowsRoute(f *testing.F) {
	seeds := []string{
		"Network Destination        Netmask          Gateway        Interface       Metric\n" +
			"0.0.0.0          0.0.0.0      On-link        192.168.1.10    1\n" +
			"0.0.0.0          0.0.0.0      192.168.1.2    192.168.1.10    600\n" +
			"0.0.0.0          0.0.0.0      192.168.1.1    192.168.1.10    100\n",
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		gateway, ok := ParseWindowsRoute(text)
		if ok && !gateway.Is4() {
			t.Fatalf("ParseWindowsRoute() returned non-IPv4 address %v", gateway)
		}
	})
}
