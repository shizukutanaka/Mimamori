package main

import (
	"context"
	"io"
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func TestMarkUnauthenticatedRTSP(t *testing.T) {
	noAuth := startRTSPTestServer(t, "127.0.0.1", "RTSP/1.0 200 OK\r\nCSeq: 1\r\nContent-Type: application/sdp\r\nContent-Length: 4\r\n\r\nv=0\n")
	authRequired := startRTSPTestServer(t, "127.0.0.2", "RTSP/1.0 401 Unauthorized\r\nCSeq: 1\r\nWWW-Authenticate: Basic\r\n\r\n")

	tests := []struct {
		name   string
		opened []probe.Open
		want   []probe.Open
	}{
		{
			name:   "unauthenticated server is kept",
			opened: []probe.Open{noAuth},
			want:   []probe.Open{noAuth},
		},
		{
			name:   "authentication-required server is dropped",
			opened: []probe.Open{authRequired},
		},
		{
			name: "other ports pass through in order",
			opened: []probe.Open{
				{IP: noAuth.IP, Port: 23},
				noAuth,
				{IP: authRequired.IP, Port: 21},
				authRequired,
				{IP: netip.MustParseAddr("127.0.0.3"), Port: 5900},
			},
			want: []probe.Open{
				{IP: noAuth.IP, Port: 23},
				noAuth,
				{IP: authRequired.IP, Port: 21},
				{IP: netip.MustParseAddr("127.0.0.3"), Port: 5900},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := markUnauthenticatedRTSP(context.Background(), test.opened)
			if len(got) != len(test.want) {
				t.Fatalf("markUnauthenticatedRTSP() returned %v, want %v", got, test.want)
			}
			for index := range test.want {
				if got[index] != test.want[index] {
					t.Fatalf("markUnauthenticatedRTSP() returned %v, want %v", got, test.want)
				}
			}
		})
	}
}

func startRTSPTestServer(t *testing.T, address, response string) probe.Open {
	t.Helper()
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "0"))
	if err != nil {
		t.Skipf("cannot listen on %s: %v", address, err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			_ = connection.SetDeadline(time.Now().Add(3 * time.Second))
			request := make([]byte, 2048)
			n, _ := connection.Read(request)
			if n > 0 {
				_, _ = io.WriteString(connection, response)
				_ = connection.Close()
				continue
			}
			_ = connection.Close()
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		select {
		case <-done:
		case <-time.After(4 * time.Second):
			t.Error("RTSP fake server did not stop")
		}
	})
	port := listener.Addr().(*net.TCPAddr).Port
	return probe.Open{IP: netip.MustParseAddr(address), Port: port}
}

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
