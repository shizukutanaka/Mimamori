package lan

import "testing"

func TestParseLinuxRouteMasks(t *testing.T) {
	header := "Iface\tDestination\tGateway\tFlags\tRefCnt\tUse\tMetric\tMask\tMTU\tWindow\tIRTT\n"
	tests := []struct {
		name   string
		text   string
		want   string
		wantOK bool
	}{
		{
			name: "OpenVPN def1 prefers the real default route",
			text: header +
				"tun0\t00000000\t0100080A\t0003\t0\t0\t0\t00000080\t0\t0\t0\n" +
				"tun0\t00000080\t0100080A\t0003\t0\t0\t0\t00000080\t0\t0\t0\n" +
				"wlan0\t00000000\t0101A8C0\t0003\t0\t0\t600\t00000000\t0\t0\t0\n",
			want:   "192.168.1.1",
			wantOK: true,
		},
		{
			name: "only /1 routes are not a default route",
			text: header +
				"tun0\t00000000\t0100080A\t0003\t0\t0\t0\t00000080\t0\t0\t0\n" +
				"tun0\t00000080\t0100080A\t0003\t0\t0\t0\t00000080\t0\t0\t0\n",
			wantOK: false,
		},
		{
			name: "plain default route",
			text: header +
				"wlan0\t00000000\t0101A8C0\t0003\t0\t0\t600\t00000000\t0\t0\t0\n",
			want:   "192.168.1.1",
			wantOK: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := ParseLinuxRoute(test.text)
			if ok != test.wantOK {
				t.Fatalf("ParseLinuxRoute() ok = %v, want %v (gateway %v)", ok, test.wantOK, got)
			}
			if ok && got.String() != test.want {
				t.Fatalf("ParseLinuxRoute() = %v, want %s", got, test.want)
			}
		})
	}
}
