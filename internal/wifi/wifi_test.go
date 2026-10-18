package wifi

import (
	"context"
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
}
