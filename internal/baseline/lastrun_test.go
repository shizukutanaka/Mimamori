package baseline

import (
	"testing"
	"time"
)

func TestLastRun(t *testing.T) {
	first := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	latest := first.Add(3 * time.Hour)

	tests := []struct {
		name    string
		store   Store
		want    time.Time
		wantSet bool
	}{
		{name: "empty"},
		{
			name:    "single device",
			store:   Store{Devices: map[string]Seen{"one": {LastSeen: first}}},
			want:    first,
			wantSet: true,
		},
		{
			name: "maximum across devices",
			store: Store{Devices: map[string]Seen{
				"one":   {LastSeen: first},
				"two":   {LastSeen: latest},
				"three": {LastSeen: first.Add(time.Hour)},
			}},
			want:    latest,
			wantSet: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := test.store.LastRun()
			if ok != test.wantSet || !got.Equal(test.want) {
				t.Fatalf("LastRun() = (%v, %v), want (%v, %v)", got, ok, test.want, test.wantSet)
			}
		})
	}
}
