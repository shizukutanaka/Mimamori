package i18n

import (
	"strings"
	"testing"
)

func TestWLEDCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["wled_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "wled_noauth")
		}
		fix, ok := messages[lang]["wled_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "wled_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s wled_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s wled_noauth_fix has an unexpected %%s placeholder", lang)
		}
		if !strings.Contains(message, "WLED") || !strings.Contains(fix, "WLED") {
			t.Errorf("%s WLED name is missing or translated", lang)
		}
		for _, label := range []string{"Config", "Security & Updates", "Settings PIN", "Save"} {
			if !strings.Contains(fix, label) {
				t.Errorf("%s fix is missing untranslated UI label %q", lang, label)
			}
		}
	}
	if got := T("en", "wled_noauth", "Living room"); got != "The WLED light controller on Living room has no settings PIN, so anyone on your network can change its settings and replace its firmware." {
		t.Errorf("English WLED message = %q", got)
	}
	if got := T("en", "wled_noauth_fix"); got != "Open the device's WLED web page, go to Config → Security & Updates, enter a 4-digit Settings PIN, and press Save." {
		t.Errorf("English WLED fix = %q", got)
	}
	if got := T("ja", "wled_noauth", "居間"); got != "「居間」の WLED（LED の照明）に設定の PIN がなく、家の中の誰でも設定を変えたり、ファームウェアを書き換えたりできる状態です。" {
		t.Errorf("Japanese WLED message = %q", got)
	}
	if got := T("ja", "wled_noauth_fix"); got != "WLED の画面を開き、「Config」→「Security & Updates」の「Settings PIN」に4桁の番号を入れて「Save」を押してください。" {
		t.Errorf("Japanese WLED fix = %q", got)
	}
}
