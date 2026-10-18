package wifi

type wlanInterfaceInfo struct {
	InterfaceGUID [16]byte
	Description   [256]uint16
	State         uint32
}

type wlanInterfaceInfoList struct {
	NumberOfItems uint32
	Index         uint32
	InterfaceInfo [1]wlanInterfaceInfo
}

type dot11SSID struct {
	Length uint32
	SSID   [32]byte
}

type wlanAssociationAttributes struct {
	SSID          dot11SSID
	BSSType       uint32
	BSSID         [6]byte
	PhyType       uint32
	PhyIndex      uint32
	SignalQuality uint32
	RxRate        uint32
	TxRate        uint32
}

type wlanSecurityAttributes struct {
	SecurityEnabled int32
	OneXEnabled     int32
	AuthAlgorithm   uint32
	CipherAlgorithm uint32
}

type wlanConnectionAttributes struct {
	State          uint32
	ConnectionMode uint32
	ProfileName    [256]uint16
	Association    wlanAssociationAttributes
	Security       wlanSecurityAttributes
}

const (
	wlanInterfaceStateConnected     uint32 = 1
	wlanIntfOpcodeCurrentConnection uint32 = 7
	dot11AuthAlgorithmOpen          uint32 = 1
	dot11AuthAlgorithmSharedKey     uint32 = 2
	dot11AuthAlgorithmWPA           uint32 = 3
	dot11AuthAlgorithmWPAPSK        uint32 = 4
	dot11AuthAlgorithmWPANone       uint32 = 5
	dot11AuthAlgorithmRSNA          uint32 = 6
	dot11AuthAlgorithmRSNAPSK       uint32 = 7
	dot11AuthAlgorithmWPA3Ent192    uint32 = 8
	dot11AuthAlgorithmWPA3SAE       uint32 = 9
	dot11AuthAlgorithmOWE           uint32 = 10
	dot11AuthAlgorithmWPA3Ent       uint32 = 11
	dot11CipherAlgorithmNone        uint32 = 0x00
	dot11CipherAlgorithmWEP40       uint32 = 0x01
	dot11CipherAlgorithmTKIP        uint32 = 0x02
	dot11CipherAlgorithmCCMP        uint32 = 0x04
	dot11CipherAlgorithmWEP104      uint32 = 0x05
	dot11CipherAlgorithmGCMP        uint32 = 0x08
	dot11CipherAlgorithmGCMP256     uint32 = 0x09
	dot11CipherAlgorithmCCMP256     uint32 = 0x0a
	dot11CipherAlgorithmWEP         uint32 = 0x101
)

func reduceWLAN(best Security, failed bool) (Security, bool) {
	if best != Unknown {
		return best, true
	}
	if failed {
		return Unknown, false
	}
	return Unknown, true
}

func wlanSecurity(auth, cipher uint32) Security {
	switch {
	case auth == dot11AuthAlgorithmSharedKey ||
		cipher == dot11CipherAlgorithmWEP40 ||
		cipher == dot11CipherAlgorithmWEP104 ||
		cipher == dot11CipherAlgorithmWEP:
		return WEP
	case auth == dot11AuthAlgorithmOpen && cipher == dot11CipherAlgorithmNone:
		return Open
	case cipher == dot11CipherAlgorithmTKIP:
		return WPA1orTKIP
	case auth == dot11AuthAlgorithmWPA ||
		auth == dot11AuthAlgorithmWPAPSK ||
		auth == dot11AuthAlgorithmWPANone:
		return WPA1orTKIP
	case isModernWLANAuthentication(auth) && isModernWLANCipher(cipher):
		return Modern
	default:
		return Unknown
	}
}

func isModernWLANAuthentication(auth uint32) bool {
	switch auth {
	case dot11AuthAlgorithmRSNA, dot11AuthAlgorithmRSNAPSK,
		dot11AuthAlgorithmWPA3Ent192, dot11AuthAlgorithmWPA3SAE,
		dot11AuthAlgorithmOWE, dot11AuthAlgorithmWPA3Ent:
		return true
	default:
		return false
	}
}

func isModernWLANCipher(cipher uint32) bool {
	switch cipher {
	case dot11CipherAlgorithmCCMP, dot11CipherAlgorithmGCMP,
		dot11CipherAlgorithmGCMP256, dot11CipherAlgorithmCCMP256:
		return true
	default:
		return false
	}
}
