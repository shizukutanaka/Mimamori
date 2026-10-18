package wifi

import (
	"syscall"
	"unsafe"
)

func currentWLAN() (Security, bool, error) {
	dll := syscall.NewLazyDLL("wlanapi.dll")
	openHandle := dll.NewProc("WlanOpenHandle")
	enumInterfaces := dll.NewProc("WlanEnumInterfaces")
	queryInterface := dll.NewProc("WlanQueryInterface")
	freeMemory := dll.NewProc("WlanFreeMemory")
	closeHandle := dll.NewProc("WlanCloseHandle")
	for _, proc := range []*syscall.LazyProc{openHandle, enumInterfaces, queryInterface, freeMemory, closeHandle} {
		if err := proc.Find(); err != nil {
			return Unknown, false, nil
		}
	}

	var negotiatedVersion uint32
	var handle syscall.Handle
	result, _, _ := openHandle.Call(
		2,
		0,
		uintptr(unsafe.Pointer(&negotiatedVersion)),
		uintptr(unsafe.Pointer(&handle)),
	)
	if result != 0 {
		return Unknown, false, nil
	}
	defer closeHandle.Call(uintptr(handle), 0)

	var list *wlanInterfaceInfoList
	result, _, _ = enumInterfaces.Call(
		uintptr(handle),
		0,
		uintptr(unsafe.Pointer(&list)),
	)
	if result != 0 {
		if list != nil {
			freeMemory.Call(uintptr(unsafe.Pointer(list)))
		}
		return Unknown, false, nil
	}
	if list == nil {
		return Unknown, true, nil
	}
	defer freeMemory.Call(uintptr(unsafe.Pointer(list)))
	if list.NumberOfItems == 0 {
		return Unknown, true, nil
	}

	interfaces := unsafe.Slice(&list.InterfaceInfo[0], int(list.NumberOfItems))
	best := Unknown
	for i := range interfaces {
		iface := &interfaces[i]
		if iface.State != wlanInterfaceStateConnected {
			continue
		}

		var size uint32
		var data *wlanConnectionAttributes
		result, _, _ = queryInterface.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
			uintptr(wlanIntfOpcodeCurrentConnection),
			0,
			uintptr(unsafe.Pointer(&size)),
			uintptr(unsafe.Pointer(&data)),
			0,
		)
		if data != nil {
			defer freeMemory.Call(uintptr(unsafe.Pointer(data)))
		}
		if result == 5 {
			return Unknown, true, ErrLocationPermission
		}
		if result != 0 {
			return Unknown, false, nil
		}
		if data == nil || uintptr(size) < unsafe.Sizeof(wlanConnectionAttributes{}) ||
			data.State != wlanInterfaceStateConnected {
			continue
		}

		candidate := wlanSecurity(data.Security.AuthAlgorithm, data.Security.CipherAlgorithm)
		if securityRank(candidate) > securityRank(best) {
			best = candidate
		}
	}
	return best, true, nil
}
