package i18n

import (
	"strings"
	"testing"
)

func TestVNCCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"vnc_noauth", "vnc_noauth_fix"} {
			if message := messages[lang][key]; message == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
		if count := strings.Count(messages[lang]["vnc_noauth"], "%s"); count != 1 {
			t.Errorf("%s vnc_noauth has %d %%s placeholders, want 1", lang, count)
		}
	}
}
