package console

import "testing"

func TestQuickEditMode(t *testing.T) {
	tests := []struct {
		name   string
		mode   uint32
		want   uint32
		change bool
	}{
		{
			name:   "clears quick edit and enables extended flags",
			mode:   0x0001 | enableQuickEditMode,
			want:   0x0001 | enableExtendedFlags,
			change: true,
		},
		{
			name:   "does not change mode without quick edit",
			mode:   0x0001,
			want:   0x0001,
			change: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, change := quickEditMode(test.mode)
			if got != test.want || change != test.change {
				t.Fatalf("quickEditMode(%#x) = (%#x, %v), want (%#x, %v)", test.mode, got, change, test.want, test.change)
			}
		})
	}
}
