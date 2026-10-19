package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
)

func TestDaysSinceLastCheck(t *testing.T) {
	now := time.Date(2025, 4, 15, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		last time.Time
		want int
	}{
		{name: "today", last: now, want: 0},
		{name: "one day", last: now.Add(-24 * time.Hour), want: 1},
		{name: "twelve days", last: now.Add(-12 * 24 * time.Hour), want: 12},
		{name: "clock skew", last: now.Add(time.Hour), want: 0},
		{
			name: "calendar day in JST",
			last: time.Date(2025, 4, 14, 23, 55, 0, 0, time.FixedZone("JST", 9*3600)),
			want: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := daysSinceLastCheck(now, test.last, time.UTC); got != test.want {
				t.Fatalf("daysSinceLastCheck() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestDaysSinceLastCheckCalendarDates(t *testing.T) {
	loc := time.FixedZone("JST", 9*3600)
	nextDay := time.Date(2025, 4, 15, 0, 5, 0, 0, loc)
	tests := []struct {
		name string
		now  time.Time
		last time.Time
		loc  *time.Location
		want int
	}{
		{
			name: "crosses JST midnight",
			now:  nextDay,
			last: time.Date(2025, 4, 14, 23, 55, 0, 0, loc),
			loc:  loc,
			want: 1,
		},
		{
			name: "same date despite elapsed time",
			now:  time.Date(2025, 4, 14, 23, 55, 0, 0, loc),
			last: time.Date(2025, 4, 14, 0, 5, 0, 0, loc),
			loc:  loc,
			want: 0,
		},
		{
			name: "two calendar boundaries in 47 hours",
			now:  time.Date(2025, 4, 3, 0, 10, 0, 0, loc),
			last: time.Date(2025, 4, 1, 0, 30, 0, 0, loc),
			loc:  loc,
			want: 2,
		},
		{
			name: "clock skew",
			now:  time.Date(2025, 4, 14, 23, 55, 0, 0, loc),
			last: nextDay,
			loc:  loc,
			want: 0,
		},
		{
			name: "UTC dates differ from JST",
			now:  nextDay,
			last: time.Date(2025, 4, 14, 23, 55, 0, 0, loc),
			loc:  time.UTC,
			want: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := daysSinceLastCheck(test.now, test.last, test.loc); got != test.want {
				t.Fatalf("daysSinceLastCheck() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestRenderTextLastCheck(t *testing.T) {
	tests := []struct {
		name string
		lang string
		days int
		want string
	}{
		{name: "English today", lang: "en", days: 0, want: "Last check: today"},
		{name: "English singular", lang: "en", days: 1, want: "Last check: 1 day ago"},
		{name: "English plural", lang: "en", days: 12, want: "Last check: 12 days ago"},
		{name: "Japanese today", lang: "ja", days: 0, want: "前回の確認: 今日"},
		{name: "Japanese singular", lang: "ja", days: 1, want: "前回の確認: 1日前"},
		{name: "Japanese plural", lang: "ja", days: 12, want: "前回の確認: 12日前"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			days := test.days
			var buffer bytes.Buffer
			renderText(&buffer, output{
				Score: 100, Status: "safe",
			}, test.lang, &days)

			lines := strings.Split(strings.TrimSuffix(buffer.String(), "\n"), "\n")
			countLine := i18n.T(test.lang, "device_counts", 0, 0)
			countIndex := -1
			for index, line := range lines {
				if line == countLine {
					countIndex = index
					break
				}
			}
			if countIndex < 0 || countIndex+1 >= len(lines) {
				t.Fatalf("device-count line not found in output: %q", buffer.String())
			}
			if lines[countIndex+1] != test.want {
				t.Fatalf("line after device count = %q, want %q", lines[countIndex+1], test.want)
			}
		})
	}
}

func TestRenderTextOmitsLastCheckOnFirstRun(t *testing.T) {
	days := 12
	var buffer bytes.Buffer
	renderText(&buffer, output{
		Score: 100, Status: "safe", FirstRun: true,
	}, "ja", &days)

	if strings.Contains(buffer.String(), "前回の確認:") {
		t.Fatalf("first-run output contains a last-check line: %q", buffer.String())
	}
}

func TestWriteOutputOmitsLastCheckFromJSON(t *testing.T) {
	days := 12
	var buffer bytes.Buffer
	if err := writeOutput(&buffer, output{Score: 100, Status: "safe"}, true, "en", &days); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buffer.String(), "last_check") {
		t.Fatalf("JSON output contains last-check data: %q", buffer.String())
	}
}
