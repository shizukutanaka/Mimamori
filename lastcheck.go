package main

import "time"

func daysSinceLastCheck(now, last time.Time) int {
	days := int(now.Sub(last).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}
