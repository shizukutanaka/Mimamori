package i18n

import (
	"strings"
	"testing"
)

func TestSMB1CatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"smb1_enabled", "smb1_enabled_fix"} {
			if message := messages[lang][key]; message == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
		if count := strings.Count(messages[lang]["smb1_enabled"], "%s"); count != 1 {
			t.Errorf("%s smb1_enabled has %d %%s placeholders, want 1", lang, count)
		}
	}
}
