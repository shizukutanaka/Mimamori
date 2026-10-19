//go:build !windows

package console

func DisableQuickEdit() func() {
	return func() {}
}
