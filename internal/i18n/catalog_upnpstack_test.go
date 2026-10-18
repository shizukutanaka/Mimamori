package i18n

import "testing"

func TestUPnPStackCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"upnp_stack_vulnerable", "upnp_stack_vulnerable_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}
