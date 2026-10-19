//go:build !windows

package baseline

import (
	"path/filepath"
	"testing"
)

func TestSaveWithDirectorySync(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "devices.json")
	if err := Save(path, emptyStore()); err != nil {
		t.Fatal(err)
	}
}
