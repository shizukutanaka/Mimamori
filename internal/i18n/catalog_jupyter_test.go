package i18n

import (
	"strings"
	"testing"
)

func TestJupyterCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["jupyter_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "jupyter_noauth")
		}
		fix, ok := messages[lang]["jupyter_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "jupyter_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s jupyter_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if !strings.Contains(message, "Jupyter") || !strings.Contains(message, "8888") {
			t.Errorf("%s jupyter_noauth must contain Jupyter and 8888", lang)
		}
		if !strings.Contains(fix, "token") {
			t.Errorf("%s jupyter_noauth_fix must contain token", lang)
		}
	}
	if got := T("en", "jupyter_noauth", "192.168.1.5"); got != "Jupyter (port 8888) on 192.168.1.5 is open to your network with no token or password. Anyone at home can run any program on this computer." {
		t.Errorf("English Jupyter message = %q", got)
	}
	if got := T("en", "jupyter_noauth_fix"); got != "Start Jupyter with a token or password (don't launch it with --IdentityProvider.token='' or --NotebookApp.token='')." {
		t.Errorf("English Jupyter fix = %q", got)
	}
	if got := T("ja", "jupyter_noauth", "居間"); got != "「居間」の Jupyter（8888番）は、トークンもパスワードもなしで家のネットワークに開いています。家の誰でも、この PC で好きなプログラムを動かせる状態です。" {
		t.Errorf("Japanese Jupyter message = %q", got)
	}
	if got := T("ja", "jupyter_noauth_fix"); got != "Jupyter を起動するときに token かパスワードを設定してください（`--IdentityProvider.token=''` や `--NotebookApp.token=''` を付けて起動しない）。" {
		t.Errorf("Japanese Jupyter fix = %q", got)
	}
}
