//go:build windows

package baseline

import (
	"errors"
	"os"
	"syscall"
	"time"
)

func renameReplace(tempPath, path string) error {
	return retryRename(func() error {
		return os.Rename(tempPath, path)
	}, isEphemeralRenameError, time.Sleep, 2*time.Second)
}

func isEphemeralRenameError(err error) bool {
	var errno syscall.Errno
	return errors.As(err, &errno) && (errno == syscall.Errno(5) || errno == syscall.Errno(32))
}
