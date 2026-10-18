//go:build !windows

package console

func ShouldPause() bool {
	return false
}
