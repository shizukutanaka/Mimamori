//go:build windows

package console

import (
	"syscall"
	"unsafe"
)

func DisableQuickEdit() func() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getStdHandle := kernel32.NewProc("GetStdHandle")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	handle, _, _ := getStdHandle.Call(uintptr(^uint32(9)))
	if handle == 0 || handle == ^uintptr(0) {
		return func() {}
	}

	var mode uint32
	ok, _, _ := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if ok == 0 {
		return func() {}
	}

	newMode, change := quickEditMode(mode)
	if !change {
		return func() {}
	}
	ok, _, _ = setConsoleMode.Call(handle, uintptr(newMode))
	if ok == 0 {
		return func() {}
	}

	return func() {
		_, _, _ = setConsoleMode.Call(handle, uintptr(mode))
	}
}
