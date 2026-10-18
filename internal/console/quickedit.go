package console

const (
	enableQuickEditMode uint32 = 0x0040
	enableExtendedFlags uint32 = 0x0080
)

func quickEditMode(mode uint32) (uint32, bool) {
	if mode&enableQuickEditMode == 0 {
		return mode, false
	}
	return (mode | enableExtendedFlags) &^ enableQuickEditMode, true
}
