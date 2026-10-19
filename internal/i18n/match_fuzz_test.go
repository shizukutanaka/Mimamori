package i18n

import (
	"strings"
	"testing"
)

func FuzzMatch(f *testing.F) {
	for _, tag := range localeFuzzSeeds() {
		f.Add(tag)
	}

	supported := supportedLanguages()
	f.Fuzz(func(t *testing.T, tag string) {
		lang, ok := Match(tag)
		if ok {
			if _, exists := supported[lang]; !exists {
				t.Fatalf("Match(%q) returned unsupported language %q", tag, lang)
			}
		}
	})
}

func FuzzDetect(f *testing.F) {
	for _, value := range localeFuzzSeeds() {
		f.Add(value, "", "", "")
		f.Add("", value, "", "")
		f.Add("", "", value, "")
		f.Add("", "", "", value)
	}

	supported := supportedLanguages()
	f.Fuzz(func(t *testing.T, language, lcAll, lcMessages, lang string) {
		env := map[string]string{
			"LANGUAGE":    language,
			"LC_ALL":      lcAll,
			"LC_MESSAGES": lcMessages,
			"LANG":        lang,
		}
		got := Detect(func(key string) string { return env[key] })
		if _, ok := supported[got]; !ok {
			t.Fatalf("Detect() returned unsupported language %q", got)
		}
	})
}

func localeFuzzSeeds() []string {
	return []string{
		"",
		"C",
		"POSIX",
		"ja_JP.UTF-8",
		"zh-Hant-TW",
		"zh_Hans_HK",
		"pt_BR@latin",
		":::",
		string([]byte{0xff, 0xfe, 0x80}),
		strings.Repeat("x", 10*1024),
	}
}

func supportedLanguages() map[string]struct{} {
	supported := make(map[string]struct{})
	for _, lang := range Supported() {
		supported[lang] = struct{}{}
	}
	return supported
}
