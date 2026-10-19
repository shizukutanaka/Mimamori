//go:build darwin

package i18n

import (
	"context"
	"os/exec"
	"time"
)

func systemLocale() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "/usr/bin/defaults", "read", "-g", "AppleLanguages").Output()
	if err != nil {
		return ""
	}
	return parseAppleLanguages(string(output))
}
