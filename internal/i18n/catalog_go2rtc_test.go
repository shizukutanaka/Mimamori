package i18n

import "testing"

func TestGo2RTCCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"go2rtc_noauth", "go2rtc_noauth_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}
