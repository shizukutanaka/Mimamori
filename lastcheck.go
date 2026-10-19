package main

import (
	"time"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
)

func daysSinceLastCheck(now, last time.Time, loc *time.Location) int {
	nowYear, nowMonth, nowDay := now.In(loc).Date()
	lastYear, lastMonth, lastDay := last.In(loc).Date()
	nowDate := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(lastYear, lastMonth, lastDay, 0, 0, 0, 0, time.UTC)
	days := int(nowDate.Sub(lastDate).Hours() / 24)
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
