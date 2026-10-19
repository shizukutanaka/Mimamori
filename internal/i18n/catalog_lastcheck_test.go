package i18n

import "testing"

func TestLastCheckCatalogCoverage(t *testing.T) {
	keys := []string{"last_check_today", "last_check_day", "last_check_days"}
	for _, lang := range Supported() {
		for _, key := range keys {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}
