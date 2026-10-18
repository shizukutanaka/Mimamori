package i18n

import (
	"syscall"
	"unsafe"
)

func systemLocale() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	var buffer [85]uint16
	languageID, _, _ := kernel32.NewProc("GetUserDefaultUILanguage").Call()
	if languageID != 0 {
		result, _, _ := kernel32.NewProc("LCIDToLocaleName").Call(
			uintptr(uint32(languageID)),
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			0,
		)
		if result != 0 {
			return syscall.UTF16ToString(buffer[:])
		}
	}
	result, _, _ := kernel32.NewProc("GetUserDefaultLocaleName").Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	if result == 0 {
		return ""
	}
	return syscall.UTF16ToString(buffer[:])
}
