package baseline

import (
	"path/filepath"
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
		{
			name: "last check newer than devices",
			store: Store{
				LastCheck: latest.Add(time.Hour),
				Devices:   map[string]Seen{"one": {LastSeen: latest}},
			},
			want:    latest.Add(time.Hour),
			wantSet: true,
		},
		{
			name:    "last check without devices",
			store:   Store{LastCheck: latest},
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

func TestLastCheckSaveLoad(t *testing.T) {
	lastCheck := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	path := filepath.Join(t.TempDir(), "devices.json")
	if err := Save(path, Store{Version: 1, Devices: map[string]Seen{}, LastCheck: lastCheck}); err != nil {
		t.Fatal(err)
	}
	loaded, existed, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !existed {
		t.Fatal("Load() existed = false, want true")
	}
	if !loaded.LastCheck.Equal(lastCheck) {
		t.Fatalf("LastCheck = %v, want %v", loaded.LastCheck, lastCheck)
	}
}
