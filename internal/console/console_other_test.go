//go:build !windows

package console

import "testing"

func TestShouldPause(t *testing.T) {
	if ShouldPause() {
		t.Fatal("ShouldPause() = true on a non-Windows platform")
	}
}
