package i18n

import "testing"

func TestUsageCatalogCoverage(t *testing.T) {
	keys := []string{
		"usage_header", "flag_json", "flag_lang", "flag_state", "flag_wait", "flag_version",
		"unexpected_argument", "invalid_lang", "wait_negative",
	}
	for _, lang := range Supported() {
		for _, key := range keys {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}
