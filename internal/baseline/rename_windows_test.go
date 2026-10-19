//go:build windows

package baseline

import (
	"errors"
	"os"
	"syscall"
	"testing"
)

func TestIsEphemeralRenameError(t *testing.T) {
	for _, code := range []syscall.Errno{5, 32} {
		err := &os.LinkError{Op: "rename", Old: "temp", New: "history", Err: code}
		if !isEphemeralRenameError(err) {
			t.Errorf("isEphemeralRenameError(%v) = false, want true", code)
		}
	}

	for _, err := range []error{
		&os.LinkError{Op: "rename", Old: "temp", New: "history", Err: syscall.Errno(1)},
		errors.New("other error"),
	} {
		if isEphemeralRenameError(err) {
			t.Errorf("isEphemeralRenameError(%v) = true, want false", err)
		}
	}
}
