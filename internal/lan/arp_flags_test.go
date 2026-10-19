package lan

import (
	"net/netip"
	"testing"
)

func TestParseARPProcFlags(t *testing.T) {
	text := `IP address       HW type     Flags       HW address            Mask     Device
192.168.1.1       0x1         0x2         a4:2b:b0:11:22:33     *        wlan0
192.168.1.20      0x1         0x6         a4:2b:b0:44:55:66     *        wlan0
192.168.1.30      0x1         0x0         3c:22:fb:77:88:99     *        wlan0
192.168.1.40      0x1         0xzz        3c:22:fb:aa:bb:cc     *        wlan0`
	got := ParseARP(text)
	want := []Neighbor{
		{IP: netip.MustParseAddr("192.168.1.1"), MAC: "a4:2b:b0:11:22:33"},
		{IP: netip.MustParseAddr("192.168.1.20"), MAC: "a4:2b:b0:44:55:66"},
	}
	if len(got) != len(want) {
		t.Fatalf("ParseARP() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ParseARP()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestParseARPMacOSAndWindowsLines(t *testing.T) {
	text := "? (192.168.1.5) at a4:2b:b0:0:0:5 on en0 ifscope [ethernet]\n" +
		"  192.168.1.6          a4-2b-b0-00-00-06     dynamic"
	got := ParseARP(text)
	want := []Neighbor{
		{IP: netip.MustParseAddr("192.168.1.5"), MAC: "a4:2b:b0:00:00:05"},
		{IP: netip.MustParseAddr("192.168.1.6"), MAC: "a4:2b:b0:00:00:06"},
	}
	if len(got) != len(want) {
		t.Fatalf("ParseARP() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ParseARP()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}
