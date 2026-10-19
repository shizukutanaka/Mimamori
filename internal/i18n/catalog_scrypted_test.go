package i18n

import (
	"strings"
	"testing"
)

func TestScryptedCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["scrypted_onboarding"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "scrypted_onboarding")
		}
		fix, ok := messages[lang]["scrypted_onboarding_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "scrypted_onboarding_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s scrypted_onboarding has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s scrypted_onboarding_fix has an unexpected %%s placeholder", lang)
		}
		if !strings.Contains(message, "Scrypted") || !strings.Contains(fix, "Scrypted") {
			t.Errorf("%s Scrypted name is missing or translated", lang)
		}
		if !strings.Contains(message, "11080") {
			t.Errorf("%s finding is missing port 11080", lang)
		}
		if !strings.Contains(fix, "10443") {
			t.Errorf("%s fix is missing port 10443", lang)
		}
	}
	if got := T("en", "scrypted_onboarding", "192.168.1.5"); got != "Scrypted (port 11080) on 192.168.1.5 hasn't had its first account created yet. Anyone on your network can create the admin account first, then install plugins and run commands." {
		t.Errorf("English Scrypted message = %q", got)
	}
	if got := T("en", "scrypted_onboarding_fix"); got != "Open Scrypted (https://<this device>:10443) and create the admin account now. If you don't use it, turn it off." {
		t.Errorf("English Scrypted fix = %q", got)
	}
	if got := T("ja", "scrypted_onboarding", "192.168.1.5"); got != "「192.168.1.5」の Scrypted（11080番）は、最初のアカウントがまだ作られていません。家のネットワークの誰でも、先に管理者のアカウントを作り、プラグインを入れて命令を実行できる状態です。" {
		t.Errorf("Japanese Scrypted message = %q", got)
	}
	if got := T("ja", "scrypted_onboarding_fix"); got != "Scrypted（https://この機器:10443）を開いて、いますぐ管理者のアカウントを作ってください。使っていなければ、止めてください。" {
		t.Errorf("Japanese Scrypted fix = %q", got)
	}
}
