package i18n

import (
	"syscall"
	"unsafe"
)

func systemLocale() string {
	procedure := syscall.NewLazyDLL("kernel32.dll").NewProc("GetUserDefaultLocaleName")
	var buffer [85]uint16
	result, _, _ := procedure.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	if result == 0 {
		return ""
	}
	return syscall.UTF16ToString(buffer[:])
}
