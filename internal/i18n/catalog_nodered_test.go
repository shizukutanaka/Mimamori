package i18n

import "testing"

func TestNodeREDCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"nodered_noauth", "nodered_noauth_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
	if got := T("en", "nodered_noauth", "Home Assistant"); got != "The Node-RED editor (port 1880) on Home Assistant has no password. Anyone on your network can change your flows and run commands on that device." {
		t.Errorf("English Node-RED message = %q", got)
	}
	if got := T("ja", "nodered_noauth_fix"); got != "Node-RED の settings.js で adminAuth を設定してパスワードを付け、Node-RED を再起動してください。" {
		t.Errorf("Japanese Node-RED fix = %q", got)
	}
}
