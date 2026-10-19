package i18n

import (
	"strings"
	"testing"
)

func TestShellyCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["shelly_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "shelly_noauth")
		}
		fix, ok := messages[lang]["shelly_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "shelly_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s shelly_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s shelly_noauth_fix has an unexpected %%s placeholder", lang)
		}
		if !strings.Contains(message, "Shelly") || !strings.Contains(fix, "Shelly") {
			t.Errorf("%s Shelly name is missing or translated", lang)
		}
	}
	if got := T("en", "shelly_noauth", "Kitchen"); got != "The Shelly smart device on Kitchen has no password, so anyone on your network can switch it on or off and replace its firmware." {
		t.Errorf("English Shelly message = %q", got)
	}
	if got := T("en", "shelly_noauth_fix"); got != "In the Shelly app or the device's web page, open the authentication settings, turn on password protection, and set a password." {
		t.Errorf("English Shelly fix = %q", got)
	}
	if got := T("ja", "shelly_noauth", "台所"); got != "「台所」の Shelly の機器にパスワードがなく、家の中の誰でも電源を入れたり切ったり、ファームウェアを書き換えたりできる状態です。" {
		t.Errorf("Japanese Shelly message = %q", got)
	}
	if got := T("ja", "shelly_noauth_fix"); got != "Shelly のアプリか機器の画面で、認証（パスワードの保護）の設定を開いてオンにし、パスワードを設定してください。" {
		t.Errorf("Japanese Shelly fix = %q", got)
	}
}
