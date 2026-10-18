package baseline

import (
	"testing"
	"time"
)

func TestPruneRemovesOnlyOldPrunableDevices(t *testing.T) {
	now := time.Date(2025, 1, 31, 12, 0, 0, 0, time.UTC)
	store := Store{Devices: map[string]Seen{
		"02:00:00:00:00:01": {LastSeen: now.Add(-31 * 24 * time.Hour)},
		"02:00:00:00:00:02": {LastSeen: now.Add(-time.Hour)},
		"02:00:00:00:00:03": {LastSeen: now.Add(-30 * 24 * time.Hour)},
		"00:11:22:33:44:55": {LastSeen: now.Add(-90 * 24 * time.Hour)},
	}}

	pruned := store.Prune(now, 30*24*time.Hour, func(mac string) bool {
		return mac == "02:00:00:00:00:01" || mac == "02:00:00:00:00:02" || mac == "02:00:00:00:00:03"
	})
	if pruned != 1 {
		t.Fatalf("Prune() = %d, want 1", pruned)
	}
	if _, ok := store.Devices["02:00:00:00:00:01"]; ok {
		t.Fatal("old prunable device was not removed")
	}
	for _, mac := range []string{"02:00:00:00:00:02", "02:00:00:00:00:03", "00:11:22:33:44:55"} {
		if _, ok := store.Devices[mac]; !ok {
			t.Errorf("device %s was removed", mac)
		}
	}
}
