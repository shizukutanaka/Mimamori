package i18n

import (
	"reflect"
	"regexp"
	"testing"
)

var formatVerbPattern = regexp.MustCompile(`%[-+# 0-9.]*[a-zA-Z%]`)

func TestCatalogKeysAndFormatVerbs(t *testing.T) {
	wantLanguages := []string{"en", "ja", "zh-Hans", "zh-Hant", "es", "pt", "fr", "de", "ru", "ko", "id"}
	supported := Supported()
	if !reflect.DeepEqual(supported, wantLanguages) {
		t.Fatalf("Supported() = %v, want %v", supported, wantLanguages)
	}
	supportedSet := make(map[string]bool, len(supported))
	for _, lang := range supported {
		if supportedSet[lang] {
			t.Errorf("Supported() contains duplicate language %q", lang)
		}
		supportedSet[lang] = true
	}
	if len(messages) != len(supportedSet) {
		t.Fatalf("Supported() has %d languages, messages has %d", len(supportedSet), len(messages))
	}
	for lang := range messages {
		if !supportedSet[lang] {
			t.Errorf("messages contains unsupported language %q", lang)
		}
	}
	for _, lang := range supported {
		if lang == "en" {
			continue
		}
		t.Run(lang, func(t *testing.T) {
			english := messages["en"]
			catalog := messages[lang]
			if len(catalog) != len(english) {
				t.Errorf("catalog has %d keys, English has %d", len(catalog), len(english))
			}
			for key, englishText := range english {
				translated, ok := catalog[key]
				if !ok {
					t.Errorf("catalog missing key %q", key)
					continue
				}
				if got, want := formatVerbPattern.FindAllString(translated, -1), formatVerbPattern.FindAllString(englishText, -1); !reflect.DeepEqual(got, want) {
					t.Errorf("format verbs for %q = %v, want %v", key, got, want)
				}
			}
			for key := range catalog {
				if _, ok := english[key]; !ok {
					t.Errorf("catalog has extra key %q", key)
				}
			}
		})
	}
}

func TestADBCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"adb_open", "adb_open_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}

func TestSupplementalCatalogKeys(t *testing.T) {
	keys := []string{
		"wifi_open", "wifi_open_fix", "wifi_wep", "wifi_wep_fix", "wifi_weak", "wifi_weak_fix",
		"incomplete_notice",
	}
	for _, lang := range Supported() {
		for _, key := range keys {
			if _, ok := messages[lang][key]; !ok {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		tag  string
		want string
		ok   bool
	}{
		{"zh_TW.UTF-8", "zh-Hant", true},
		{"zh-Hant-HK", "zh-Hant", true},
		{"zh_HK", "zh-Hant", true},
		{"zh_CN.UTF-8", "zh-Hans", true},
		{"zh-Hans-HK", "zh-Hans", true},
		{"zh_Hans_TW", "zh-Hans", true},
		{"zh-Hant-CN", "zh-Hant", true},
		{"zh", "zh-Hans", true},
		{"zh-Hans-SG", "zh-Hans", true},
		{"pt_BR.UTF-8", "pt", true},
		{"de-AT", "de", true},
		{"ko_KR", "ko", true},
		{"id_ID", "id", true},
		{"ru_RU.KOI8-R", "ru", true},
		{"ja_JP.eucJP", "ja", true},
		{"C", "", false},
		{"POSIX", "", false},
		{"", "", false},
		{"sr_RS@latin", "", false},
		{"xx", "", false},
		{"ZH-HANT", "zh-Hant", true},
		{"zh-MO", "zh-Hant", true},
		{" en-US ", "en", true},
	}
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			got, ok := Match(test.tag)
			if got != test.want || ok != test.ok {
				t.Fatalf("Match(%q) = (%q, %t), want (%q, %t)", test.tag, got, ok, test.want, test.ok)
			}
		})
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
		{"language priority list", map[string]string{"LANGUAGE": "xx:de_DE:fr", "LANG": "es_ES"}, "de"},
		{"language falls through to locale", map[string]string{"LANGUAGE": "xx:sr_RS", "LANG": "fr_FR.UTF-8"}, "fr"},
		{"LC_ALL C overrides LANG", map[string]string{"LC_ALL": "C", "LANG": "ja"}, "en"},
		{"invalid first locale falls back to English", map[string]string{"LC_MESSAGES": "sr_RS@latin", "LANG": "de"}, "en"},
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
