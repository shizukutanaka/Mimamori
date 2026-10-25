package i18n

import "testing"

func TestFrigateCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"frigate_noauth", "frigate_noauth_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
	if got := T("en", "frigate_noauth", "Home Assistant"); got != "The Frigate NVR (port 5000) on Home Assistant has no password. Anyone on your network can watch your cameras and change its settings as an administrator." {
		t.Errorf("English Frigate message = %q", got)
	}
	if got := T("ja", "frigate_noauth_fix"); got != "Frigate を 0.14 以降に更新し、5000番の公開（Docker の ports の 5000:5000）を外して、ログインが必要な 8971番を使ってください。" {
		t.Errorf("Japanese Frigate fix = %q", got)
	}
}
