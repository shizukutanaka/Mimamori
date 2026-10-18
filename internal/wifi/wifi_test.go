package wifi

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurrentReturnsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	security, err := Current(ctx)
	if security != Unknown {
		t.Fatalf("Current() security = %v, want Unknown", security)
	}
	if err != context.Canceled {
		t.Fatalf("Current() error = %v, want %v", err, context.Canceled)
	}
}

func TestParseMacOS(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want Security
	}{
		{name: "open", mode: "spairport_security_mode_none", want: Open},
		{name: "WEP", mode: "spairport_security_mode_wep", want: WEP},
		{name: "WPA personal", mode: "spairport_security_mode_wpa_personal", want: WPA1orTKIP},
		{name: "WPA enterprise", mode: "spairport_security_mode_wpa_enterprise", want: WPA1orTKIP},
		{name: "WPA2 personal", mode: "spairport_security_mode_wpa2_personal", want: Modern},
		{name: "WPA3 personal", mode: "spairport_security_mode_wpa3_personal", want: Modern},
		{name: "unknown", mode: "spairport_security_mode_unknown", want: Unknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := []byte(`{
				"SPAirPortDataType": [{
					"_name": "spairport_airport_card_type_airport_extreme",
					"spairport_airport_interfaces": [{
						"_name": "en0",
						"spairport_current_network_information": {
							"_name": "Family Wi-Fi",
							"spairport_security_mode": "` + test.mode + `"
						}
					}]
				}]
			}`)
			if got := ParseMacOS(fixture); got != test.want {
				t.Fatalf("ParseMacOS() = %v, want %v", got, test.want)
			}
		})
	}
	if got := ParseMacOS([]byte(`{"SPAirPortDataType":[{"spairport_airport_interfaces":[]}]}`)); got != Unknown {
		t.Fatalf("ParseMacOS(no current network) = %v, want Unknown", got)
	}
	if got := ParseMacOS([]byte("{")); got != Unknown {
		t.Fatalf("ParseMacOS(invalid JSON) = %v, want Unknown", got)
	}
}

func TestParseNmcli(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Security
	}{
		{name: "open active network", text: "no:WPA2\nyes:\n", want: Open},
		{name: "WEP", text: "yes:WEP\n", want: WEP},
		{name: "WPA1", text: "yes:WPA1\n", want: WPA1orTKIP},
		{name: "WPA2", text: "no:WEP\nyes:WPA2\n", want: Modern},
		{name: "WPA3", text: "yes:WPA3\n", want: Modern},
		{name: "WPA1 mixed with WPA2 remains weak", text: "yes:WPA1 WPA2\n", want: WPA1orTKIP},
		{name: "WPA2 mixed with WPA3 is modern", text: "yes:WPA2 WPA3\n", want: Modern},
		{name: "localized active value is not accepted", text: "はい:WPA2\n", want: Unknown},
		{name: "no active network", text: "no:WPA2\n", want: Unknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseNmcli(test.text); got != test.want {
				t.Fatalf("ParseNmcli() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestParseWPACli(t *testing.T) {
	fixture := func(ssid, state, keyManagement, pairwiseCipher, groupCipher string) string {
		return strings.Join([]string{
			"bssid=aa:bb:cc:dd:ee:ff",
			"freq=2412",
			"ssid=" + ssid,
			"id=0",
			"mode=station",
			"pairwise_cipher=" + pairwiseCipher,
			"group_cipher=" + groupCipher,
			"key_mgmt=" + keyManagement,
			"wpa_state=" + state,
			"address=aa:bb:cc:dd:ee:ff",
			"ip_address=192.168.1.42",
		}, "\n")
	}
	tests := []struct {
		name string
		text string
		want Security
	}{
		{name: "completed WPA2 PSK with CCMP", text: fixture("Home Wi-Fi", "COMPLETED", "WPA2-PSK", "CCMP", "CCMP"), want: Modern},
		{name: "SAE", text: fixture("Home Wi-Fi", "COMPLETED", "SAE", "CCMP", "CCMP"), want: Modern},
		{name: "OWE", text: fixture("Home Wi-Fi", "COMPLETED", "OWE", "CCMP", "CCMP"), want: Modern},
		{name: "WPA2 enterprise", text: fixture("Enterprise", "COMPLETED", "WPA2/IEEE 802.1X/EAP", "CCMP", "CCMP"), want: Modern},
		{name: "group TKIP is weak", text: fixture("Home Wi-Fi", "COMPLETED", "WPA2-PSK", "CCMP", "TKIP"), want: WPA1orTKIP},
		{name: "WPA2 mixed with WPA1 is weak", text: fixture("Home Wi-Fi", "COMPLETED", "WPA2-PSK+WPA-PSK", "CCMP", "CCMP"), want: WPA1orTKIP},
		{name: "WPA PSK with TKIP", text: fixture("Home Wi-Fi", "COMPLETED", "WPA-PSK", "TKIP", "TKIP"), want: WPA1orTKIP},
		{name: "WPA enterprise is weak", text: fixture("Enterprise", "COMPLETED", "WPA/IEEE 802.1X/EAP", "CCMP", "CCMP"), want: WPA1orTKIP},
		{name: "NONE is open", text: fixture("Guest", "COMPLETED", "NONE", "NONE", "NONE"), want: Open},
		{name: "NONE with WEP cipher", text: fixture("Legacy", "COMPLETED", "NONE", "WEP-104", "NONE"), want: WEP},
		{name: "dynamic WEP", text: fixture("Legacy", "COMPLETED", "IEEE 802.1X (no WPA)", "WEP-40", "WEP-40"), want: WEP},
		{name: "scanning", text: fixture("Home Wi-Fi", "SCANNING", "WPA2-PSK", "CCMP", "CCMP"), want: Unknown},
		{name: "disconnected", text: fixture("Home Wi-Fi", "DISCONNECTED", "WPA2-PSK", "CCMP", "CCMP"), want: Unknown},
		{name: "empty status", text: "", want: Unknown},
		{name: "empty key management", text: fixture("Home Wi-Fi", "COMPLETED", "", "CCMP", "CCMP"), want: Unknown},
		{name: "SSID values do not affect classification", text: fixture("WEP TKIP NONE", "COMPLETED", "WPA2-PSK", "CCMP", "CCMP"), want: Modern},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseWPACli(test.text); got != test.want {
				t.Fatalf("ParseWPACli() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestParseNetsh(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Security
	}{
		{
			name: "English open network",
			text: `Name                   : Wi-Fi
State                  : connected
Authentication         : Open
Cipher                 : None`,
			want: Open,
		},
		{
			name: "English WEP network",
			text: `Name                   : Wi-Fi
State                  : connected
Authentication         : Open
Cipher                 : WEP`,
			want: WEP,
		},
		{
			name: "SSID named WEP does not override WPA2",
			text: `Name                   : Wi-Fi
State                  : connected
SSID                   : WEP
Authentication         : WPA2-Personal
Cipher                 : CCMP`,
			want: Modern,
		},
		{
			name: "SSID named TKIP does not override WPA2",
			text: `Name                   : Wi-Fi
State                  : connected
SSID                   : TKIP
Authentication         : WPA2-Personal
Cipher                 : CCMP`,
			want: Modern,
		},
		{
			name: "TKIP network",
			text: `Name                   : Wi-Fi
State                  : connected
Authentication         : WPA-Personal
Cipher                 : TKIP`,
			want: WPA1orTKIP,
		},
		{
			name: "WPA-Personal with CCMP remains weak",
			text: `Name                   : Wi-Fi
State                  : connected
Authentication         : WPA-Personal
Cipher                 : CCMP`,
			want: WPA1orTKIP,
		},
		{
			name: "WPA-Enterprise with CCMP remains weak",
			text: `Name                   : Wi-Fi
State                  : connected
Authentication         : WPA-Enterprise
Cipher                 : CCMP`,
			want: WPA1orTKIP,
		},
		{
			name: "Japanese labels with WPA2 and CCMP values",
			text: `名前                   : Wi-Fi
状態                   : 接続済み
認証                   : WPA2-Personal
暗号                   : CCMP`,
			want: Modern,
		},
		{
			name: "WPA2 with TKIP is weak",
			text: `Name                   : Wi-Fi
State                  : connected
Authentication         : WPA2-Personal
Cipher                 : TKIP`,
			want: WPA1orTKIP,
		},
		{
			name: "disconnected interface",
			text: `Name                   : Wi-Fi
State                  : disconnected
Authentication         : Open
Cipher                 : None`,
			want: Unknown,
		},
		{
			name: "disconnected WEP does not override connected WPA2",
			text: `Name                   : Wi-Fi 1
State                  : disconnected
Authentication         : Open
Cipher                 : WEP

Name                   : Wi-Fi 2
State                  : connected
Authentication         : WPA2-Personal
Cipher                 : CCMP`,
			want: Modern,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseNetsh(test.text); got != test.want {
				t.Fatalf("ParseNetsh() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestNetshNeedsLocation(t *testing.T) {
	english := strings.ReplaceAll(`There is 1 interface on the system:

Network shell commands need location permission to access WLAN information. Turn on Location services on the Location page in Privacy & security settings.

Here is the URI for the Location page in the Settings app:

ms-settings:privacy-location

To open the Location page in the Settings app, hold down the Ctrl key and select the link, or run the following command:

start ms-settings:privacy-location

Or, to open the Location page from the Run dialog box, press Windows logo key + R, and then copy and paste the URI above.

Function WlanQueryInterface returns error 5:

Access is denied.`, "\n", "\r\n")
	japanese := strings.ReplaceAll(`システム上に 1 つのインターフェイスがあります:

WLAN 情報にアクセスするには位置情報のアクセス許可が必要です。プライバシーとセキュリティの設定で位置情報サービスをオンにしてください。

設定アプリの位置情報ページを開く URI:

ms-settings:privacy-location

Function WlanQueryInterface はエラー 5 を返しました:

アクセスが拒否されました。`, "\n", "\r\n")
	for _, test := range []struct {
		name string
		text string
	}{
		{name: "English transcript", text: english},
		{name: "Japanese transcript", text: japanese},
	} {
		t.Run(test.name, func(t *testing.T) {
			if !NetshNeedsLocation(test.text) {
				t.Fatal("NetshNeedsLocation() = false, want true")
			}
			if got := ParseNetsh(test.text); got != Unknown {
				t.Fatalf("ParseNetsh(location transcript) = %v, want Unknown", got)
			}
		})
	}

	connected := `Name                   : Wi-Fi
State                  : connected
Authentication         : WPA2-Personal
Cipher                 : CCMP`
	if NetshNeedsLocation(connected) {
		t.Fatal("NetshNeedsLocation(connected interface) = true, want false")
	}

	openSSID := strings.ReplaceAll(`Name                   : Wi-Fi
State                  : connected
SSID                   : ms-settings:privacy-location
BSSID                  : 00:11:22:33:44:55
Authentication         : Open
Cipher                 : None`, "\n", "\r\n")
	if NetshNeedsLocation(openSSID) {
		t.Fatal("NetshNeedsLocation(open SSID) = true, want false")
	}
	if got := ParseNetsh(openSSID); got != Open {
		t.Fatalf("ParseNetsh(open SSID) = %v, want Open", got)
	}
}

func TestWorstSecurity(t *testing.T) {
	tests := []struct {
		name    string
		results []Security
		want    Security
	}{
		{name: "empty"},
		{name: "all unknown", results: []Security{Unknown, Unknown}},
		{name: "unknown and modern", results: []Security{Unknown, Modern}, want: Modern},
		{name: "modern and open", results: []Security{Modern, Open}, want: Open},
		{name: "WEP and WPA1", results: []Security{WEP, WPA1orTKIP}, want: WEP},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := worstSecurity(test.results); got != test.want {
				t.Fatalf("worstSecurity(%v) = %v, want %v", test.results, got, test.want)
			}
		})
	}
}

func TestSupplicantInterfaces(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "run", "wpa_supplicant")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"wlan0", "wlan1", "p2p-dev-wlan0"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}

	symlink := filepath.Join(root, "var", "run", "wpa_supplicant")
	if err := os.MkdirAll(filepath.Dir(symlink), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(dir, symlink); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(root, "missing", "wpa_supplicant")

	got := supplicantInterfaces([]string{dir, symlink, missing})
	want := []supplicantInterface{
		{dir: dir, name: "wlan0"},
		{dir: dir, name: "wlan1"},
	}
	if len(got) != len(want) {
		t.Fatalf("supplicantInterfaces() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("supplicantInterfaces()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}
