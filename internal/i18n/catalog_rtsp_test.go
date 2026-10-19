package i18n

import "testing"

func TestRTSPCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"rtsp_noauth", "rtsp_noauth_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}
