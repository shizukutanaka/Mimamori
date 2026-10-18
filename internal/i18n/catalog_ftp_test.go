package i18n

import (
	"strings"
	"testing"
)

func TestFTPCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"ftp_anonymous", "ftp_anonymous_fix"} {
			if message := messages[lang][key]; message == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
		if count := strings.Count(messages[lang]["ftp_anonymous"], "%s"); count != 1 {
			t.Errorf("%s ftp_anonymous has %d %%s placeholders, want 1", lang, count)
		}
	}
}
