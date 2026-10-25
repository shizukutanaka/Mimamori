package i18n

import (
	"strings"
	"testing"
)

func TestPiHoleCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["pihole_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "pihole_noauth")
		}
		fix, ok := messages[lang]["pihole_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "pihole_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s pihole_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s pihole_noauth_fix has an unexpected %%s placeholder", lang)
		}
		if !strings.Contains(message, "Pi-hole") || !strings.Contains(fix, "Pi-hole") {
			t.Errorf("%s Pi-hole name is missing or translated", lang)
		}
		if !strings.Contains(message, "80") {
			t.Errorf("%s finding is missing port 80", lang)
		}
		if !strings.Contains(fix, "pihole setpassword") {
			t.Errorf("%s fix is missing the pihole setpassword command", lang)
		}
	}
	if got := T("en", "pihole_noauth", "Kitchen"); got != "Pi-hole (port 80) on Kitchen has an empty admin password. Anyone on your network can change its DNS settings and send every device in your home to fake websites." {
		t.Errorf("English Pi-hole message = %q", got)
	}
	if got := T("en", "pihole_noauth_fix"); got != "Set a password in the Pi-hole admin page or with `pihole setpassword`." {
		t.Errorf("English Pi-hole fix = %q", got)
	}
	if got := T("ja", "pihole_noauth", "台所"); got != "「台所」の Pi-hole（80番）は、管理画面のパスワードが空です。家のネットワークの誰でも DNS の設定を書き換えて、家じゅうの機器を偽のサイトへ向けられる状態です。" {
		t.Errorf("Japanese Pi-hole message = %q", got)
	}
	if got := T("ja", "pihole_noauth_fix"); got != "Pi-hole の管理画面か `pihole setpassword` で、パスワードを設定してください。" {
		t.Errorf("Japanese Pi-hole fix = %q", got)
	}
}
