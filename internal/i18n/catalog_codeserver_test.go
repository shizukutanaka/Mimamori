package i18n

import (
	"strings"
	"testing"
)

func TestCodeServerCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["codeserver_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "codeserver_noauth")
		}
		fix, ok := messages[lang]["codeserver_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "codeserver_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s codeserver_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if !strings.Contains(message, "code-server") || !strings.Contains(message, "8443") {
			t.Errorf("%s codeserver_noauth must contain code-server and 8443", lang)
		}
		if !strings.Contains(fix, "PASSWORD") {
			t.Errorf("%s codeserver_noauth_fix must contain PASSWORD", lang)
		}
	}
	if got := T("en", "codeserver_noauth", "192.168.1.5"); got != "code-server (VS Code in the browser, port 8443) on 192.168.1.5 is open to your network with no password. Anyone at home can open its terminal and run any command." {
		t.Errorf("English code-server message = %q", got)
	}
	if got := T("en", "codeserver_noauth_fix"); got != "Set the PASSWORD or HASHED_PASSWORD environment variable for the container and restart it (don't start it with --auth none)." {
		t.Errorf("English code-server fix = %q", got)
	}
	if got := T("ja", "codeserver_noauth", "書斎"); got != "「書斎」の code-server（ブラウザで使う VS Code、8443番）は、合言葉なしで家のネットワークに開いています。家の誰でも、ブラウザから端末を開き、好きな命令を実行できる状態です。" {
		t.Errorf("Japanese code-server message = %q", got)
	}
	if got := T("ja", "codeserver_noauth_fix"); got != "Docker の環境変数 PASSWORD か HASHED_PASSWORD を設定して、起動し直してください（--auth none で起動しない）。" {
		t.Errorf("Japanese code-server fix = %q", got)
	}
}
