package baseline

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestObserve(t *testing.T) {
	store := emptyStore()
	now := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	if !store.Observe("aa:bb:cc:dd:ee:ff", "192.168.1.2", "Camera", now) {
		t.Fatal("first observation should be new")
	}
	if store.Observe("AA:BB:CC:DD:EE:FF", "192.168.1.3", "Camera", now.Add(time.Hour)) {
		t.Fatal("second observation should not be new")
	}
	seen := store.Devices["aa:bb:cc:dd:ee:ff"]
	if seen.LastIP != "192.168.1.3" || !seen.FirstSeen.Equal(now) || seen.Name != "Camera" {
		t.Fatalf("stored device = %#v", seen)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "devices.json")
	store := emptyStore()
	store.Observe("aa:bb:cc:dd:ee:ff", "192.168.1.2", "Camera", time.Now().UTC())
	if err := Save(path, store); err != nil {
		t.Fatal(err)
	}
	loaded, existed, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !existed || len(loaded.Devices) != 1 || loaded.Version != 1 {
		t.Fatalf("Load() = %#v, existed=%v", loaded, existed)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("file mode = %o, want 600", info.Mode().Perm())
	}
}

func TestLoadMissingFile(t *testing.T) {
	store, existed, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil || existed || store.Devices == nil {
		t.Fatalf("Load() = %#v, existed=%v, err=%v", store, existed, err)
	}
}
