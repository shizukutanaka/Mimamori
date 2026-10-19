package i18n

import "testing"

func TestDatabaseCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"forward_database", "forward_database_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}
