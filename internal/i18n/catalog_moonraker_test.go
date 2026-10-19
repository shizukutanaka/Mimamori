package i18n

import (
	"strings"
	"testing"
)

func TestMoonrakerCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	requiredFixTokens := []string{
		"Moonraker",
		"moonraker.conf",
		"[authorization]",
		"trusted_clients",
		"192.168.0.0/16",
		"force_logins: True",
	}
	for _, lang := range languages {
		message, ok := messages[lang]["moonraker_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "moonraker_noauth")
		}
		fix, ok := messages[lang]["moonraker_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "moonraker_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s moonraker_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s moonraker_noauth_fix has an unexpected %%s placeholder", lang)
		}
		for _, token := range requiredFixTokens {
			if !strings.Contains(fix, token) {
				t.Errorf("%s moonraker_noauth_fix is missing %q", lang, token)
			}
		}
	}
	if got := T("en", "moonraker_noauth", "Office"); got != "Moonraker (the 3D printer control service) on Office accepts commands from this PC without a password. With its common default settings, anyone on your network can heat the printer and start prints." {
		t.Errorf("English Moonraker message = %q", got)
	}
	if got := T("en", "moonraker_noauth_fix"); got != "In moonraker.conf, under [authorization], remove whole network ranges such as 192.168.0.0/16 from trusted_clients and list only your own devices, or create a user and set force_logins: True. Then restart Moonraker." {
		t.Errorf("English Moonraker fix = %q", got)
	}
	if got := T("ja", "moonraker_noauth", "居間"); got != "「居間」の Moonraker（3D プリンターの操作窓口）が、この PC からパスワードなしで操作できる状態です。よく使われる初めの設定のままだと、家の中の誰でもプリンターを加熱したり、印刷を始めたりできます。" {
		t.Errorf("Japanese Moonraker message = %q", got)
	}
	if got := T("ja", "moonraker_noauth_fix"); got != "moonraker.conf の [authorization] にある trusted_clients から 192.168.0.0/16 のような範囲を消して自分の機器だけを書くか、ユーザーを作って force_logins: True を設定し、Moonraker を再起動してください。" {
		t.Errorf("Japanese Moonraker fix = %q", got)
	}
}
