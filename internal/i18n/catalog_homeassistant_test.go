package i18n

import "testing"

func TestHomeAssistantOnboardingCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"homeassistant_onboarding", "homeassistant_onboarding_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
	if got := T("en", "homeassistant_onboarding", "Home Server"); got != "Home Assistant (port 8123) on Home Server hasn't finished its first-time setup. Anyone on your network can create the owner account first and control all your smart-home devices." {
		t.Errorf("English Home Assistant message = %q", got)
	}
	if got := T("ja", "homeassistant_onboarding_fix"); got != "Home Assistant を開いて、最初の設定（持ち主のアカウントの作成）を済ませてください。使っていないなら、止めてください。" {
		t.Errorf("Japanese Home Assistant fix = %q", got)
	}
}
