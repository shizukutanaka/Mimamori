package i18n

import (
	"strings"
	"testing"
)

func TestOllamaCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["ollama_open"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "ollama_open")
		}
		fix, ok := messages[lang]["ollama_open_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "ollama_open_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s ollama_open has %d %%s placeholders, want 1", lang, count)
		}
		if !strings.Contains(message, "Ollama") || !strings.Contains(message, "11434") {
			t.Errorf("%s ollama_open must contain Ollama and 11434", lang)
		}
		if !strings.Contains(fix, "OLLAMA_HOST") {
			t.Errorf("%s ollama_open_fix must contain OLLAMA_HOST", lang)
		}
	}
	if got := T("en", "ollama_open", "192.168.1.5"); got != "Ollama (local AI, port 11434) on 192.168.1.5 is open to your network, and Ollama has no password. Anyone at home can run models on this computer's GPU, delete models, and use cloud models on its ollama.com account if it is signed in." {
		t.Errorf("English Ollama message = %q", got)
	}
	if got := T("en", "ollama_open_fix"); got != "If no other device needs it, set OLLAMA_HOST to 127.0.0.1:11434 (for Docker, publish with -p 127.0.0.1:11434:11434)." {
		t.Errorf("English Ollama fix = %q", got)
	}
	if got := T("ja", "ollama_open", "居間"); got != "「居間」の Ollama（ローカルの AI、11434番）は、合言葉のしくみがないまま家のネットワークに開いています。家の誰でも、この PC の GPU で AI を動かしたり、モデルを消したりできます。ollama.com にサインインしていれば、その枠のクラウドのモデルも使えます。" {
		t.Errorf("Japanese Ollama message = %q", got)
	}
	if got := T("ja", "ollama_open_fix"); got != "ほかの機器から使わないなら、OLLAMA_HOST を 127.0.0.1:11434 にしてください（Docker なら -p 127.0.0.1:11434:11434）。" {
		t.Errorf("Japanese Ollama fix = %q", got)
	}
}
