//go:build !windows

package baseline

import "os"

func renameReplace(tempPath, path string) error {
	return os.Rename(tempPath, path)
}
