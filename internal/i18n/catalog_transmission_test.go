package i18n

import (
	"strings"
	"testing"
)

func TestTransmissionCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["transmission_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "transmission_noauth")
		}
		fix, ok := messages[lang]["transmission_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "transmission_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s transmission_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s transmission_noauth_fix has an unexpected %%s placeholder", lang)
		}
		combined := message + fix
		for _, required := range []string{"Transmission", "9091", "USER", "PASS"} {
			if !strings.Contains(combined, required) {
				t.Errorf("%s catalog is missing %q", lang, required)
			}
		}
	}
	if got := T("en", "transmission_noauth", "192.168.1.5"); got != "Transmission's remote control (port 9091) on 192.168.1.5 accepts commands from your network without a password. Anyone at home can change where files are saved and set a program to run when a download finishes." {
		t.Errorf("English Transmission message = %q", got)
	}
	if got := T("en", "transmission_noauth_fix"); got != "Turn on RPC authentication with a username and password in Transmission (for the linuxserver.io Docker image, set USER and PASS)." {
		t.Errorf("English Transmission fix = %q", got)
	}
	if got := T("ja", "transmission_noauth", "居間"); got != "「居間」の Transmission の操作窓口（9091番）は、合言葉なしで家のネットワークに開いています。家の誰でも、ファイルの保存先を変えたり、取り込みが終わったときに動かすプログラムを設定したりできる状態です。" {
		t.Errorf("Japanese Transmission message = %q", got)
	}
	if got := T("ja", "transmission_noauth_fix"); got != "Transmission の設定で、操作窓口のユーザー名と合言葉を設定してください（linuxserver.io の Docker なら USER と PASS）。" {
		t.Errorf("Japanese Transmission fix = %q", got)
	}
}
