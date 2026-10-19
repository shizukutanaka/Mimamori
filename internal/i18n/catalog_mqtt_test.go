package i18n

import "testing"

func TestMQTTCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"mqtt_anonymous", "mqtt_anonymous_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}
