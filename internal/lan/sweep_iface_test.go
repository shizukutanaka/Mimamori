package lan

import (
	"net"
	"testing"
)

func TestSweepableInterface(t *testing.T) {
	tests := []struct {
		name  string
		flags net.Flags
		want  bool
	}{
		{
			name:  "up broadcast LAN",
			flags: net.FlagUp | net.FlagBroadcast | net.FlagMulticast | net.FlagRunning,
			want:  true,
		},
		{
			name:  "down broadcast interface",
			flags: net.FlagBroadcast,
			want:  false,
		},
		{
			name:  "loopback",
			flags: net.FlagUp | net.FlagLoopback,
			want:  false,
		},
		{
			name:  "point-to-point VPN tunnel",
			flags: net.FlagUp | net.FlagPointToPoint | net.FlagMulticast | net.FlagRunning,
			want:  false,
		},
		{
			name:  "interface without broadcast",
			flags: net.FlagUp | net.FlagRunning,
			want:  false,
		},
		{
			name:  "broadcast point-to-point interface",
			flags: net.FlagUp | net.FlagBroadcast | net.FlagPointToPoint,
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := sweepableInterface(test.flags); got != test.want {
				t.Fatalf("sweepableInterface(%v) = %v, want %v", test.flags, got, test.want)
			}
		})
	}
}
