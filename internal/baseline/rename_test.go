package baseline

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRetryRenameSucceedsAfterEphemeralFailures(t *testing.T) {
	ephemeral := errors.New("temporary rename failure")
	attempts := 0
	var waits []time.Duration
	err := retryRename(
		func() error {
			attempts++
			if attempts <= 2 {
				return ephemeral
			}
			return nil
		},
		func(err error) bool { return err == ephemeral },
		func(wait time.Duration) { waits = append(waits, wait) },
		time.Second,
	)
	if err != nil {
		t.Fatalf("retryRename() error = %v", err)
	}
	if attempts != 3 {
		t.Fatalf("rename attempts = %d, want 3", attempts)
	}
	if want := []time.Duration{time.Millisecond, 2 * time.Millisecond}; !reflect.DeepEqual(waits, want) {
		t.Fatalf("sleep durations = %v, want %v", waits, want)
	}
}

func TestRetryRenameStopsOnNonEphemeralError(t *testing.T) {
	permanent := errors.New("permanent rename failure")
	attempts := 0
	var waits []time.Duration
	err := retryRename(
		func() error {
			attempts++
			return permanent
		},
		func(error) bool { return false },
		func(wait time.Duration) { waits = append(waits, wait) },
		time.Second,
	)
	if err != permanent {
		t.Fatalf("retryRename() error = %v, want %v", err, permanent)
	}
	if attempts != 1 || len(waits) != 0 {
		t.Fatalf("attempts = %d, sleeps = %v; want one attempt and no sleeps", attempts, waits)
	}
}

func TestRetryRenameReturnsLastErrorAtLimit(t *testing.T) {
	attempts := 0
	var last error
	var waits []time.Duration
	err := retryRename(
		func() error {
			attempts++
			last = errors.New("temporary rename failure")
			return last
		},
		func(error) bool { return true },
		func(wait time.Duration) { waits = append(waits, wait) },
		4*time.Millisecond,
	)
	if err != last {
		t.Fatalf("retryRename() error = %v, want last error %v", err, last)
	}
	if attempts != 4 {
		t.Fatalf("rename attempts = %d, want 4", attempts)
	}
	if want := []time.Duration{time.Millisecond, 2 * time.Millisecond, time.Millisecond}; !reflect.DeepEqual(waits, want) {
		t.Fatalf("sleep durations = %v, want %v", waits, want)
	}
}

func TestRetryRenameSucceedsAfterFinalSleep(t *testing.T) {
	ephemeral := errors.New("temporary rename failure")
	attempts := 0
	err := retryRename(
		func() error {
			attempts++
			if attempts <= 3 {
				return ephemeral
			}
			return nil
		},
		func(err error) bool { return err == ephemeral },
		func(time.Duration) {},
		4*time.Millisecond,
	)
	if err != nil {
		t.Fatalf("retryRename() error = %v, want nil", err)
	}
	if attempts != 4 {
		t.Fatalf("rename attempts = %d, want 4", attempts)
	}
}

func TestRetryRenameCapsBackoff(t *testing.T) {
	var waits []time.Duration
	err := retryRename(
		func() error { return errors.New("temporary rename failure") },
		func(error) bool { return true },
		func(wait time.Duration) { waits = append(waits, wait) },
		400*time.Millisecond,
	)
	if err == nil {
		t.Fatal("retryRename() succeeded, want limit error")
	}
	want := []time.Duration{
		time.Millisecond, 2 * time.Millisecond, 4 * time.Millisecond, 8 * time.Millisecond,
		16 * time.Millisecond, 32 * time.Millisecond, 64 * time.Millisecond,
		100 * time.Millisecond, 100 * time.Millisecond, 73 * time.Millisecond,
	}
	if !reflect.DeepEqual(waits, want) {
		t.Fatalf("sleep durations = %v, want %v", waits, want)
	}
	for _, wait := range waits {
		if wait > 100*time.Millisecond {
			t.Fatalf("sleep duration = %v, exceeds 100 ms cap", wait)
		}
	}
}
