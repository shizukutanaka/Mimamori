package i18n

import (
	"strings"
	"testing"
)

func TestSyncthingCatalogCoverage(t *testing.T) {
	languages := Supported()
	if len(languages) != 11 {
		t.Fatalf("Supported() returned %d languages, want 11", len(languages))
	}
	for _, lang := range languages {
		message, ok := messages[lang]["syncthing_noauth"]
		if !ok || message == "" {
			t.Errorf("%s catalog is missing %q", lang, "syncthing_noauth")
		}
		fix, ok := messages[lang]["syncthing_noauth_fix"]
		if !ok || fix == "" {
			t.Errorf("%s catalog is missing %q", lang, "syncthing_noauth_fix")
		}
		if count := strings.Count(message, "%s"); count != 1 {
			t.Errorf("%s syncthing_noauth has %d %%s placeholders, want 1", lang, count)
		}
		if strings.Contains(fix, "%s") {
			t.Errorf("%s syncthing_noauth_fix has an unexpected %%s placeholder", lang)
		}
		if !strings.Contains(message, "Syncthing") || !strings.Contains(fix, "Syncthing") {
			t.Errorf("%s Syncthing name is missing or translated", lang)
		}
		if !strings.Contains(message, "8384") {
			t.Errorf("%s finding is missing port 8384", lang)
		}
	}
	if got := T("en", "syncthing_noauth", "192.168.1.5"); got != "Syncthing's web GUI (port 8384) on 192.168.1.5 has no password. Anyone on your network can add their own device and copy or overwrite files on this machine." {
		t.Errorf("English Syncthing message = %q", got)
	}
	if got := T("en", "syncthing_noauth_fix"); got != "Set a GUI Authentication User and Password in Syncthing's Settings → GUI." {
		t.Errorf("English Syncthing fix = %q", got)
	}
	if got := T("ja", "syncthing_noauth", "台所"); got != "「台所」の Syncthing の管理画面（8384番）は、合言葉が設定されていません。家のネットワークの誰でも同期の相手を足して、この機器のファイルを抜き取ったり書き換えたりできる状態です。" {
		t.Errorf("Japanese Syncthing message = %q", got)
	}
	if got := T("ja", "syncthing_noauth_fix"); got != "Syncthing の管理画面の「設定」→「GUI」で、GUI 認証のユーザー名と合言葉を設定してください。" {
		t.Errorf("Japanese Syncthing fix = %q", got)
	}
}
