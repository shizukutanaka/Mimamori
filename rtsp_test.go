package main

import (
	"net/netip"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func TestRTSPNoAuthFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.23")
	device := "Front Door"
	tests := []struct {
		lang    string
		message string
		fix     string
	}{
		{
			lang:    "en",
			message: "The camera feed on Front Door can be viewed without a password (RTSP).",
			fix:     "Turn on RTSP authentication (a password) in the camera's settings, or turn RTSP off if you don't use it, and update the firmware.",
		},
		{
			lang:    "ja",
			message: "Front Door のカメラ映像が、パスワードなしで見られる状態です（RTSP）。",
			fix:     "カメラの設定で RTSP の認証（パスワード）を有効にするか、使っていなければ RTSP を無効にし、ファームウェアを更新してください。",
		},
	}

	for _, test := range tests {
		t.Run(test.lang, func(t *testing.T) {
			findings := makeFindings(
				[]probe.Open{{IP: ip, Port: 554}},
				nil,
				netip.Addr{},
				nil,
				map[netip.Addr]string{ip: device},
				nil,
				nil,
				test.lang,
			)
			if len(findings) != 1 {
				t.Fatalf("makeFindings() returned %d findings, want one: %#v", len(findings), findings)
			}
			finding := findings[0]
			if finding.ID != "rtsp_noauth" || finding.Level != report.Caution || finding.Points != 8 {
				t.Fatalf("RTSP finding = %#v, want rtsp_noauth/Caution/8", finding)
			}
			if finding.IP != ip.String() || finding.Device != device {
				t.Fatalf("RTSP finding target = (%q, %q), want (%q, %q)",
					finding.IP, finding.Device, ip.String(), device)
			}
			if finding.Message != test.message || finding.Fix != test.fix {
				t.Fatalf("RTSP finding text = (%q, %q), want (%q, %q)",
					finding.Message, finding.Fix, test.message, test.fix)
			}
		})
	}
}
