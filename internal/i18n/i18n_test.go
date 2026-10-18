package i18n

import (
	"strings"
	"testing"
)

func TestCatalogKeysAndFormatVerbs(t *testing.T) {
	for key, english := range messages["en"] {
		japanese, ok := messages["ja"][key]
		if !ok {
			t.Errorf("Japanese catalog missing key %q", key)
			continue
		}
		if strings.Count(english, "%") != strings.Count(japanese, "%") {
			t.Errorf("format verb count differs for %q: en=%q ja=%q", key, english, japanese)
		}
	}
	for key := range messages["ja"] {
		if _, ok := messages["en"][key]; !ok {
			t.Errorf("English catalog missing key %q", key)
		}
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"Japanese", map[string]string{"LANG": "ja_JP.UTF-8"}, "ja"},
		{"C locale", map[string]string{"LANG": "C"}, "en"},
		{"English", map[string]string{"LANG": "en_US.UTF-8"}, "en"},
		{"LC_ALL overrides LANG", map[string]string{"LC_ALL": "ja_JP.UTF-8", "LANG": "en_US.UTF-8"}, "ja"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Detect(func(key string) string { return test.env[key] })
			if got != test.want {
				t.Fatalf("Detect() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTranslationFallback(t *testing.T) {
	if got := T("unknown", "safe"); got != messages["en"]["safe"] {
		t.Fatalf("T() = %q", got)
	}
	if got := T("en", "missing_key"); got != "missing_key" {
		t.Fatalf("T() = %q", got)
	}
}
