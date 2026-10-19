package wifi

import (
	"encoding/binary"
	"testing"
	"unsafe"
)

type bssTestEntry struct {
	bssid    [6]byte
	elements []byte
	ieOffset uint32
	ieSize   uint32
}

func buildBSSList(t *testing.T, items uint32, entries []bssTestEntry) []byte {
	t.Helper()
	list := make([]byte, wlanBSSListEntriesOffset+len(entries)*wlanBSSEntrySize)
	binary.LittleEndian.PutUint32(list[4:], items)
	for i, entry := range entries {
		start := wlanBSSListEntriesOffset + i*wlanBSSEntrySize
		copy(list[start+wlanBSSEntryBSSIDOffset:], entry.bssid[:])
		ieOffset := entry.ieOffset
		if ieOffset == 0 {
			ieOffset = uint32(len(list) - start)
			list = append(list, entry.elements...)
		}
		ieSize := entry.ieSize
		if ieSize == 0 {
			ieSize = uint32(len(entry.elements))
		}
		binary.LittleEndian.PutUint32(list[start+wlanBSSEntryIEOffsetField:], ieOffset)
		binary.LittleEndian.PutUint32(list[start+wlanBSSEntryIESizeField:], ieSize)
	}
	binary.LittleEndian.PutUint32(list[0:], uint32(len(list)))
	return list
}

func TestBSSAdvertisesWPA1(t *testing.T) {
	bssid := [6]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	otherBSSID := [6]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	ssid := []byte("Family Wi-Fi")
	ssidElement := append([]byte{informationElementSSID, byte(len(ssid))}, ssid...)
	rsnElement := []byte{48, 2, 0x01, 0x00}
	wpa1Element := []byte{221, 6, 0x00, 0x50, 0xF2, 0x01, 0x01, 0x00}
	wmmElement := []byte{221, 6, 0x00, 0x50, 0xF2, 0x02, 0x01, 0x00}

	concat := func(parts ...[]byte) []byte {
		var all []byte
		for _, part := range parts {
			all = append(all, part...)
		}
		return all
	}

	tests := []struct {
		name    string
		items   uint32
		entries []bssTestEntry
		want    bool
	}{
		{
			name:  "RSN and WPA1 elements",
			items: 1,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat(ssidElement, rsnElement, wpa1Element)},
			},
			want: true,
		},
		{
			name:  "RSN only",
			items: 1,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat(ssidElement, rsnElement)},
			},
		},
		{
			name:  "WMM vendor element is not WPA1",
			items: 1,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat(ssidElement, rsnElement, wmmElement)},
			},
		},
		{
			name:  "different BSSID",
			items: 1,
			entries: []bssTestEntry{
				{bssid: otherBSSID, elements: concat(ssidElement, wpa1Element)},
			},
		},
		{
			name:  "different SSID element",
			items: 1,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat([]byte{informationElementSSID, 5, 'O', 't', 'h', 'e', 'r'}, wpa1Element)},
			},
		},
		{
			name:  "no SSID element",
			items: 1,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat(rsnElement, wpa1Element)},
			},
		},
		{
			name:  "truncated element",
			items: 1,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat(ssidElement, wpa1Element, []byte{48, 40, 0x01})},
			},
		},
		{
			name:  "information elements past the buffer",
			items: 1,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat(ssidElement, wpa1Element), ieOffset: wlanBSSEntrySize, ieSize: 4096},
			},
		},
		{
			name:  "second entry matches",
			items: 2,
			entries: []bssTestEntry{
				{bssid: otherBSSID, elements: concat(ssidElement, rsnElement)},
				{bssid: bssid, elements: concat(ssidElement, rsnElement, wpa1Element)},
			},
			want: true,
		},
		{
			name:  "item count larger than the buffer",
			items: 1000,
			entries: []bssTestEntry{
				{bssid: bssid, elements: concat(ssidElement, rsnElement)},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list := buildBSSList(t, test.items, test.entries)
			if got := bssAdvertisesWPA1(list, bssid, ssid); got != test.want {
				t.Fatalf("bssAdvertisesWPA1() = %t, want %t", got, test.want)
			}
		})
	}
	if bssAdvertisesWPA1(nil, bssid, ssid) {
		t.Fatal("bssAdvertisesWPA1(nil) = true, want false")
	}
}

func TestWLANSecurity(t *testing.T) {
	tests := []struct {
		name   string
		auth   uint32
		cipher uint32
		want   Security
	}{
		{name: "shared key is WEP", auth: dot11AuthAlgorithmSharedKey, cipher: dot11CipherAlgorithmNone, want: WEP},
		{name: "WEP40 cipher", auth: dot11AuthAlgorithmOpen, cipher: dot11CipherAlgorithmWEP40, want: WEP},
		{name: "open with WEP104 cipher", auth: dot11AuthAlgorithmOpen, cipher: dot11CipherAlgorithmWEP104, want: WEP},
		{name: "WEP cipher", auth: dot11AuthAlgorithmRSNAPSK, cipher: dot11CipherAlgorithmWEP, want: WEP},
		{name: "open network", auth: dot11AuthAlgorithmOpen, cipher: dot11CipherAlgorithmNone, want: Open},
		{name: "RSNA PSK with TKIP", auth: dot11AuthAlgorithmRSNAPSK, cipher: dot11CipherAlgorithmTKIP, want: WPA1orTKIP},
		{name: "WPA with CCMP", auth: dot11AuthAlgorithmWPA, cipher: dot11CipherAlgorithmCCMP, want: WPA1orTKIP},
		{name: "WPA PSK with CCMP", auth: dot11AuthAlgorithmWPAPSK, cipher: dot11CipherAlgorithmCCMP, want: WPA1orTKIP},
		{name: "WPA none with CCMP", auth: dot11AuthAlgorithmWPANone, cipher: dot11CipherAlgorithmCCMP, want: WPA1orTKIP},
		{name: "RSNA with CCMP", auth: dot11AuthAlgorithmRSNA, cipher: dot11CipherAlgorithmCCMP, want: Modern},
		{name: "RSNA PSK with GCMP", auth: dot11AuthAlgorithmRSNAPSK, cipher: dot11CipherAlgorithmGCMP, want: Modern},
		{name: "WPA3 enterprise 192 with GCMP", auth: dot11AuthAlgorithmWPA3Ent192, cipher: dot11CipherAlgorithmGCMP, want: Modern},
		{name: "WPA3 SAE with GCMP256", auth: dot11AuthAlgorithmWPA3SAE, cipher: dot11CipherAlgorithmGCMP256, want: Modern},
		{name: "OWE with CCMP", auth: dot11AuthAlgorithmOWE, cipher: dot11CipherAlgorithmCCMP, want: Modern},
		{name: "WPA3 enterprise with CCMP256", auth: dot11AuthAlgorithmWPA3Ent, cipher: dot11CipherAlgorithmCCMP256, want: Modern},
		{name: "RSNA PSK with no cipher", auth: dot11AuthAlgorithmRSNAPSK, cipher: dot11CipherAlgorithmNone, want: Unknown},
		{name: "IHV authentication", auth: 0x80000001, cipher: dot11CipherAlgorithmCCMP, want: Unknown},
		{name: "unknown cipher", auth: dot11AuthAlgorithmRSNA, cipher: 0x03, want: Unknown},
		{name: "open with CCMP", auth: dot11AuthAlgorithmOpen, cipher: dot11CipherAlgorithmCCMP, want: Unknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := wlanSecurity(test.auth, test.cipher); got != test.want {
				t.Fatalf("wlanSecurity(%#x, %#x) = %v, want %v", test.auth, test.cipher, got, test.want)
			}
		})
	}
}

func TestReduceWLAN(t *testing.T) {
	tests := []struct {
		name        string
		best        Security
		failed      bool
		want        Security
		wantHandled bool
	}{
		{name: "known result without failure", best: Modern, want: Modern, wantHandled: true},
		{name: "known result despite failure", best: Open, failed: true, want: Open, wantHandled: true},
		{name: "unknown without failure", want: Unknown, wantHandled: true},
		{name: "unknown with failure", failed: true, want: Unknown, wantHandled: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, handled := reduceWLAN(test.best, test.failed)
			if got != test.want || handled != test.wantHandled {
				t.Fatalf("reduceWLAN(%v, %t) = (%v, %t), want (%v, %t)",
					test.best, test.failed, got, handled, test.want, test.wantHandled)
			}
		})
	}
}

func TestWLANAPILayout(t *testing.T) {
	tests := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{name: "interface info size", got: unsafe.Sizeof(wlanInterfaceInfo{}), want: 532},
		{name: "interface info list items offset", got: unsafe.Offsetof(wlanInterfaceInfoList{}.InterfaceInfo), want: 8},
		{name: "connection attributes size", got: unsafe.Sizeof(wlanConnectionAttributes{}), want: 604},
		{name: "connection security offset", got: unsafe.Offsetof(wlanConnectionAttributes{}.Security), want: 588},
		{name: "association PHY type offset", got: unsafe.Offsetof(wlanAssociationAttributes{}.PhyType), want: 48},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("%s = %d, want %d", test.name, test.got, test.want)
			}
		})
	}
}
