package console

import (
	"syscall"
	"unsafe"
)

func ShouldPause() bool {
	processes := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleProcessList")
	var processIDs [2]uint32
	count, _, _ := processes.Call(uintptr(unsafe.Pointer(&processIDs[0])), uintptr(len(processIDs)))
	return count == 1
}
