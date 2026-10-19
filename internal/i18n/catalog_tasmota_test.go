package i18n

import (
	"strings"
	"testing"
)

func TestTasmotaCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	requiredLabels := []string{
		"Tasmota",
		"Configuration",
		"Other",
		"Web Admin Password",
		"Save",
		"Console",
		"WebPassword",
	}
	for _, lang := range languages {
		message, ok := messages[lang]["tasmota_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "tasmota_noauth")
		}
		fix, ok := messages[lang]["tasmota_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "tasmota_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s tasmota_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s tasmota_noauth_fix has an unexpected %%s placeholder", lang)
		}
		for _, label := range requiredLabels {
			if !strings.Contains(fix, label) {
				t.Errorf("%s tasmota_noauth_fix is missing UI label %q", lang, label)
			}
		}
	}
	if got := T("en", "tasmota_noauth", "Kitchen"); got != "The Tasmota smart device on Kitchen has no web password, so anyone on your network can switch it on or off, run commands, and replace its firmware." {
		t.Errorf("English Tasmota message = %q", got)
	}
	if got := T("en", "tasmota_noauth_fix"); got != "Open the device's Tasmota web page, go to Configuration → Other, enter a Web Admin Password, and press Save. Or, in the Console, run WebPassword followed by a password." {
		t.Errorf("English Tasmota fix = %q", got)
	}
	if got := T("ja", "tasmota_noauth", "台所"); got != "「台所」の Tasmota の機器に管理画面のパスワードがなく、家の中の誰でも電源を入れたり切ったり、命令を実行させたり、ファームウェアを書き換えたりできる状態です。" {
		t.Errorf("Japanese Tasmota message = %q", got)
	}
	if got := T("ja", "tasmota_noauth_fix"); got != "機器の Tasmota の画面を開き、「Configuration」→「Other」の「Web Admin Password」にパスワードを入れて「Save」を押してください。または「Console」で WebPassword に続けてパスワードを入力してください。" {
		t.Errorf("Japanese Tasmota fix = %q", got)
	}
}
