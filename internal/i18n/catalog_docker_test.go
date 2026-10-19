package i18n

import "testing"

func TestDockerCatalogCoverage(t *testing.T) {
	for _, lang := range Supported() {
		for _, key := range []string{"docker_noauth", "docker_noauth_fix"} {
			if value, ok := messages[lang][key]; !ok || value == "" {
				t.Errorf("%s catalog is missing %q", lang, key)
			}
		}
	}
	if got := T("en", "docker_noauth", "Home Server"); got != "The Docker remote API (port 2375) on Home Server has no password. Anyone on your network can run programs on it and take over the whole device." {
		t.Errorf("English Docker message = %q", got)
	}
	if got := T("ja", "docker_noauth_fix"); got != "Docker の 2375番（暗号化なし）の公開をやめてください。離れた場所から使うときは、TLS（2376番）か SSH を使ってください。" {
		t.Errorf("Japanese Docker fix = %q", got)
	}
}
