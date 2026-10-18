package wifi

import (
	"strings"
	"testing"
)

func FuzzParseNetsh(f *testing.F) {
	seeds := []string{
		`Name                   : Wi-Fi
State                  : connected
Authentication         : Open
Cipher                 : None`,
		`Name                   : Wi-Fi
State                  : connected
Authentication         : Open
Cipher                 : WEP`,
		`Name                   : Wi-Fi
State                  : connected
SSID                   : WEP
Authentication         : WPA2-Personal
Cipher                 : CCMP`,
		`Name                   : Wi-Fi
State                  : connected
SSID                   : TKIP
Authentication         : WPA2-Personal
Cipher                 : CCMP`,
		`Name                   : Wi-Fi
State                  : connected
Authentication         : WPA-Personal
Cipher                 : TKIP`,
		`Name                   : Wi-Fi
State                  : connected
Authentication         : WPA-Personal
Cipher                 : CCMP`,
		`Name                   : Wi-Fi
State                  : connected
Authentication         : WPA-Enterprise
Cipher                 : CCMP`,
		`名前                   : Wi-Fi
状態                   : 接続済み
認証                   : WPA2-Personal
暗号                   : CCMP`,
		`Name                   : Wi-Fi
State                  : connected
Authentication         : WPA2-Personal
Cipher                 : TKIP`,
		`Name                   : Wi-Fi
State                  : disconnected
Authentication         : Open
Cipher                 : None`,
		`Name                   : Wi-Fi 1
State                  : disconnected
Authentication         : Open
Cipher                 : WEP

Name                   : Wi-Fi 2
State                  : connected
Authentication         : WPA2-Personal
Cipher                 : CCMP`,
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		_ = ParseNetsh(text)
	})
}

func FuzzParseNmcli(f *testing.F) {
	for _, seed := range []string{
		"no:WPA2\nyes:\n",
		"yes:WEP\n",
		"yes:WPA1\n",
		"no:WEP\nyes:WPA2\n",
		"yes:WPA3\n",
		"yes:WPA1 WPA2\n",
		"yes:WPA2 WPA3\n",
		"はい:WPA2\n",
		"no:WPA2\n",
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		_ = ParseNmcli(text)
	})
}

func FuzzParseMacOS(f *testing.F) {
	for _, mode := range []string{
		"spairport_security_mode_none",
		"spairport_security_mode_wep",
		"spairport_security_mode_wpa_personal",
		"spairport_security_mode_wpa_enterprise",
		"spairport_security_mode_wpa2_personal",
		"spairport_security_mode_wpa3_personal",
		"spairport_security_mode_unknown",
	} {
		f.Add([]byte(`{
			"SPAirPortDataType": [{
				"_name": "spairport_airport_card_type_airport_extreme",
				"spairport_airport_interfaces": [{
					"_name": "en0",
					"spairport_current_network_information": {
						"_name": "Family Wi-Fi",
						"spairport_security_mode": "` + mode + `"
					}
				}]
			}]
		}`))
	}
	for _, seed := range [][]byte{
		[]byte(`{"SPAirPortDataType":[{"spairport_airport_interfaces":[]}]}`),
		[]byte("{"),
		{},
		[]byte(strings.Repeat("x", 1<<16)),
		{0xff, 0xfe, 0x80},
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		_ = ParseMacOS(data)
	})
}
