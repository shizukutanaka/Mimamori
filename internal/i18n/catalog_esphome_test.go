package i18n

import "testing"

func TestESPHomeCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"esphome_noauth", "esphome_noauth_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
	if got := T("en", "esphome_noauth", "Office ESPHome"); got != "Anyone can use the ESPHome dashboard (Office ESPHome) without a password (port 6052)" {
		t.Errorf("English ESPHome message = %q", got)
	}
	if got := T("ja", "esphome_noauth_fix"); got != "ESPHome にユーザー名とパスワードを設定して（Docker なら環境変数 ESPHOME_USERNAME と ESPHOME_PASSWORD。古い版は USERNAME と PASSWORD）、再起動してください。" {
		t.Errorf("Japanese ESPHome fix = %q", got)
	}
}
