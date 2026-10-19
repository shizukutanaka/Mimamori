package i18n

import (
	"strings"
	"testing"
)

func TestHomebridgeCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["homebridge_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "homebridge_noauth")
		}
		fix, ok := messages[lang]["homebridge_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "homebridge_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s homebridge_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s homebridge_noauth_fix has an unexpected %%s placeholder", lang)
		}
		if !strings.Contains(message, "Homebridge") || !strings.Contains(fix, "Homebridge") {
			t.Errorf("%s Homebridge name is missing or translated", lang)
		}
		for _, setting := range []string{`config.json`, `"platform": "config"`, `"auth": "none"`, `"form"`} {
			if !strings.Contains(fix, setting) {
				t.Errorf("%s fix is missing untranslated setting %q", lang, setting)
			}
		}
	}
	if got := T("en", "homebridge_noauth", "192.168.1.5"); got != "The Homebridge UI on 192.168.1.5 has login turned off, so anyone on your network can sign in as an administrator, install plugins, and run commands." {
		t.Errorf("English Homebridge message = %q", got)
	}
	if got := T("en", "homebridge_noauth_fix"); got != "In Homebridge's config.json, remove \"auth\": \"none\" from the \"platform\": \"config\" section (or change it to \"form\"), then restart Homebridge." {
		t.Errorf("English Homebridge fix = %q", got)
	}
	if got := T("ja", "homebridge_noauth", "192.168.1.5"); got != "「192.168.1.5」の Homebridge の管理画面でログインが切られていて、家の中の誰でも管理者として入り、プラグインを入れたり命令を実行したりできる状態です。" {
		t.Errorf("Japanese Homebridge message = %q", got)
	}
	if got := T("ja", "homebridge_noauth_fix"); got != "Homebridge の config.json で、\"platform\": \"config\" の部分にある \"auth\": \"none\" を消すか \"form\" に変えて、Homebridge を再起動してください。" {
		t.Errorf("Japanese Homebridge fix = %q", got)
	}
}
