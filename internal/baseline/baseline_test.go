package baseline

import (
	"errors"
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

func TestLoadDistinguishesParseAndReadErrors(t *testing.T) {
	directory := t.TempDir()
	corruptPath := filepath.Join(directory, "corrupt.json")
	if err := os.WriteFile(corruptPath, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	_, existed, err := Load(corruptPath)
	if !existed || !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("corrupt Load() existed=%v, err=%v; want existed and parse error", existed, err)
	}

	_, existed, err = Load(directory)
	if err == nil || existed || errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("directory Load() existed=%v, err=%v; want read error without existence", existed, err)
	}
}

func TestQuarantineOverwritesPreviousBackup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devices.json")
	corrupt := []byte("{invalid")
	if err := os.WriteFile(path, corrupt, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path+".corrupt", []byte("older backup"), 0600); err != nil {
		t.Fatal(err)
	}

	quarantinePath, err := Quarantine(path)
	if err != nil {
		t.Fatal(err)
	}
	if quarantinePath != path+".corrupt" {
		t.Fatalf("Quarantine() path = %q, want %q", quarantinePath, path+".corrupt")
	}
	got, err := os.ReadFile(quarantinePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(corrupt) {
		t.Fatalf("quarantine contents = %q, want %q", got, corrupt)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("original history still exists or stat failed: %v", err)
	}
}
