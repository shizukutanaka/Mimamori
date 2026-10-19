package i18n

import "testing"

func TestHistoryCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		if value, ok := messages[lang]["history_unavailable"]; !ok || value == "" {
			t.Errorf("%s catalog is missing %q", lang, "history_unavailable")
		}
	}
}
