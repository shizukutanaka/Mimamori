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
	getBSSList := dll.NewProc("WlanGetNetworkBssList")
	for _, proc := range []*syscall.LazyProc{openHandle, enumInterfaces, queryInterface, freeMemory, closeHandle} {
		if err := proc.Find(); err != nil {
			return Unknown, false, nil
		}
	}
	bssListAvailable := getBSSList.Find() == nil

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
	failed := false
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
			failed = true
			continue
		}
		if data == nil || uintptr(size) < unsafe.Sizeof(wlanConnectionAttributes{}) ||
			data.State != wlanInterfaceStateConnected {
			continue
		}

		candidate := wlanSecurity(data.Security.AuthAlgorithm, data.Security.CipherAlgorithm)
		if candidate == Modern && bssListAvailable {
			var bssList *byte
			result, _, _ = getBSSList.Call(
				uintptr(handle),
				uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
				uintptr(unsafe.Pointer(&data.Association.SSID)),
				uintptr(data.Association.BSSType),
				1,
				0,
				uintptr(unsafe.Pointer(&bssList)),
			)
			if result == 0 && bssList != nil {
				total := *(*uint32)(unsafe.Pointer(bssList))
				if total >= wlanBSSListEntriesOffset && total <= 1<<24 {
					ssidLength := data.Association.SSID.Length
					if ssidLength > 32 {
						ssidLength = 32
					}
					if bssAdvertisesWPA1(
						unsafe.Slice(bssList, total),
						data.Association.BSSID,
						data.Association.SSID.SSID[:ssidLength],
					) {
						candidate = WPA1orTKIP
					}
				}
			}
			if bssList != nil {
				freeMemory.Call(uintptr(unsafe.Pointer(bssList)))
			}
		}
		if securityRank(candidate) > securityRank(best) {
			best = candidate
		}
	}
	best, handled := reduceWLAN(best, failed)
	return best, handled, nil
}
