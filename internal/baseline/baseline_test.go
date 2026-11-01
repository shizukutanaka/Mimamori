package baseline

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
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
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Fatalf("file mode = %o, want 600", info.Mode().Perm())
	}
}

func TestSaveSyncsParentsOfCreatedDirectories(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "a", "b", "c")
	path := filepath.Join(directory, "history.json")

	previousSyncDir := syncDirFunc
	var synced []string
	syncDirFunc = func(path string) {
		synced = append(synced, path)
	}
	t.Cleanup(func() { syncDirFunc = previousSyncDir })

	store := emptyStore()
	if err := Save(path, store); err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(root, "a", "b"),
		filepath.Join(root, "a"),
		root,
		directory,
	}
	if !reflect.DeepEqual(synced, want) {
		t.Fatalf("first Save() synced %v, want %v", synced, want)
	}

	synced = nil
	if err := Save(path, store); err != nil {
		t.Fatal(err)
	}
	want = []string{directory}
	if !reflect.DeepEqual(synced, want) {
		t.Fatalf("second Save() synced %v, want %v", synced, want)
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

func TestQuarantinePreservesPreviousBackup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devices.json")
	firstContents := []byte("first corrupt history")
	if err := os.WriteFile(path, firstContents, 0600); err != nil {
		t.Fatal(err)
	}
	firstBackup, err := Quarantine(path)
	if err != nil {
		t.Fatal(err)
	}
	if firstBackup != path+".corrupt" {
		t.Fatalf("first Quarantine() path = %q, want %q", firstBackup, path+".corrupt")
	}
	if err := os.WriteFile(path, []byte("second corrupt history"), 0600); err != nil {
		t.Fatal(err)
	}
	secondBackup, err := Quarantine(path)
	if err != nil {
		t.Fatal(err)
	}
	if secondBackup != path+".corrupt.1" {
		t.Fatalf("second Quarantine() path = %q, want %q", secondBackup, path+".corrupt.1")
	}

	got, err := os.ReadFile(firstBackup)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(firstContents) {
		t.Fatalf("first backup contents = %q, want %q", got, firstContents)
	}
	got, err = os.ReadFile(secondBackup)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second corrupt history" {
		t.Fatalf("second backup contents = %q", got)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("original history still exists or stat failed: %v", err)
	}
}

func TestQuarantineFailsWhenBackupNamesAreExhausted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devices.json")
	original := []byte("corrupt history")
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatal(err)
	}
	for index := 0; index <= 99; index++ {
		backupPath := path + ".corrupt"
		if index > 0 {
			backupPath += "." + strconv.Itoa(index)
		}
		if err := os.WriteFile(backupPath, []byte("existing backup"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	backupPath, err := Quarantine(path)
	if err == nil {
		t.Fatal("Quarantine() succeeded with all backup names occupied")
	}
	if backupPath != path+".corrupt.99" {
		t.Fatalf("Quarantine() path = %q, want final candidate %q", backupPath, path+".corrupt.99")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("original contents = %q, want %q", got, original)
	}
}

func TestSaveKeepsSymlinkAndUpdatesTarget(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "real", "devices.json")
	if err := Save(target, Store{Version: 1, Devices: map[string]Seen{}}); err != nil {
		t.Fatalf("Save(target) error = %v", err)
	}
	link := filepath.Join(root, "devices.json")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	seen := Seen{FirstSeen: time.Unix(1, 0).UTC(), LastSeen: time.Unix(2, 0).UTC(), LastIP: "192.168.1.10"}
	if err := Save(link, Store{Version: 1, Devices: map[string]Seen{"aa:bb:cc:dd:ee:ff": seen}}); err != nil {
		t.Fatalf("Save(link) error = %v", err)
	}
	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("Lstat(link) error = %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Save(link) replaced the symlink with a regular file")
	}
	loaded, _, err := Load(target)
	if err != nil {
		t.Fatalf("Load(target) error = %v", err)
	}
	if _, ok := loaded.Devices["aa:bb:cc:dd:ee:ff"]; !ok {
		t.Fatalf("Save(link) did not update the symlink target: %+v", loaded.Devices)
	}
}

func TestQuarantineKeepsSymlinkAndMovesTarget(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "real", "devices.json")
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "devices.json")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	backup, err := Quarantine(link)
	if err != nil {
		t.Fatalf("Quarantine(link) error = %v", err)
	}
	if want := target + ".corrupt"; backup != want {
		t.Fatalf("Quarantine(link) backup = %q, want %q", backup, want)
	}
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Quarantine(link) removed or replaced the symlink: info=%v err=%v", info, err)
	}
	if _, err := os.Stat(target); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Quarantine(link) left the corrupt target in place: %v", err)
	}
	if err := Save(link, Store{Version: 1, Devices: map[string]Seen{}}); err != nil {
		t.Fatalf("Save(link) after quarantine error = %v", err)
	}
	if _, _, err := Load(target); err != nil {
		t.Fatalf("Load(target) after re-save error = %v", err)
	}
}

func TestObserveKeepsKnownNameWhenOnlyIPIsKnown(t *testing.T) {
	var store Store
	now := time.Now().UTC()
	store.Observe("aa:bb:cc:dd:ee:ff", "192.168.1.2", "Camera", now)
	store.Observe("aa:bb:cc:dd:ee:ff", "192.168.1.9", "192.168.1.9", now.Add(time.Hour))
	seen := store.Devices["aa:bb:cc:dd:ee:ff"]
	if seen.Name != "Camera" {
		t.Fatalf("Name = %q, want the earlier real name kept when this run only knows the IP", seen.Name)
	}
	if seen.LastIP != "192.168.1.9" {
		t.Fatalf("LastIP = %q, want the new IP", seen.LastIP)
	}
	store.Observe("aa:bb:cc:dd:ee:ff", "192.168.1.9", "Camera 2", now.Add(2*time.Hour))
	if seen := store.Devices["aa:bb:cc:dd:ee:ff"]; seen.Name != "Camera 2" {
		t.Fatalf("Name = %q, want a newly resolved real name to replace the old one", seen.Name)
	}
}
