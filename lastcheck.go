package main

import (
	"time"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
)

func daysSinceLastCheck(now, last time.Time) int {
	days := int(now.Sub(last).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

func lastCheckLine(result output, days *int, lang string) (string, bool) {
	if days == nil || result.FirstRun {
		return "", false
	}
	switch {
	case *days <= 0:
		return i18n.T(lang, "last_check_today"), true
	case *days == 1:
		return i18n.T(lang, "last_check_day"), true
	default:
		return i18n.T(lang, "last_check_days", *days), true
	}
}
