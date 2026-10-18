package wifi

import (
	"testing"
	"unsafe"
)

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
