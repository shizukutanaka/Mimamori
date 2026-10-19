package main

import (
	"context"
	"io"
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func TestMarkUnauthenticatedMQTT(t *testing.T) {
	anonymous := startMQTTTestBroker(t, "127.0.0.1", 0x00)
	authRequired := startMQTTTestBroker(t, "127.0.0.2", 0x05)

	tests := []struct {
		name   string
		opened []probe.Open
		want   []probe.Open
	}{
		{
			name:   "anonymous broker is kept",
			opened: []probe.Open{anonymous},
			want:   []probe.Open{anonymous},
		},
		{
			name:   "authentication-required broker is dropped",
			opened: []probe.Open{authRequired},
		},
		{
			name: "other ports pass through in order",
			opened: []probe.Open{
				{IP: anonymous.IP, Port: 21},
				anonymous,
				{IP: authRequired.IP, Port: 23},
				authRequired,
				{IP: netip.MustParseAddr("127.0.0.3"), Port: 5900},
			},
			want: []probe.Open{
				{IP: anonymous.IP, Port: 21},
				anonymous,
				{IP: authRequired.IP, Port: 23},
				{IP: netip.MustParseAddr("127.0.0.3"), Port: 5900},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := markUnauthenticatedMQTT(context.Background(), test.opened)
			if len(got) != len(test.want) {
				t.Fatalf("markUnauthenticatedMQTT() returned %v, want %v", got, test.want)
			}
			for index := range test.want {
				if got[index] != test.want[index] {
					t.Fatalf("markUnauthenticatedMQTT() returned %v, want %v", got, test.want)
				}
			}
		})
	}
}

func startMQTTTestBroker(t *testing.T, address string, connackCode byte) probe.Open {
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
			header := make([]byte, 2)
			if _, err := io.ReadFull(connection, header); err == nil {
				body := make([]byte, int(header[1]))
				if _, err := io.ReadFull(connection, body); err == nil {
					_, _ = connection.Write([]byte{0x20, 0x02, 0x00, connackCode})
				}
			}
			_ = connection.Close()
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		select {
		case <-done:
		case <-time.After(4 * time.Second):
			t.Error("MQTT fake broker did not stop")
		}
	})
	port := listener.Addr().(*net.TCPAddr).Port
	return probe.Open{IP: netip.MustParseAddr(address), Port: port}
}

func TestMQTTAnonymousFinding(t *testing.T) {
	ip := netip.MustParseAddr("192.168.1.10")
	for _, lang := range []string{"en", "ja"} {
		t.Run(lang, func(t *testing.T) {
			findings := makeFindings(
				[]probe.Open{{IP: ip, Port: 1883}},
				nil,
				netip.Addr{},
				map[netip.Addr]ssdp.Description{},
				map[netip.Addr]string{ip: "Kitchen hub"},
				nil,
				nil,
				lang,
			)
			if len(findings) != 1 {
				t.Fatalf("makeFindings() returned %d findings, want 1: %#v", len(findings), findings)
			}
			finding := findings[0]
			if finding.ID != "mqtt_anonymous" || finding.Level != report.Caution || finding.Points != 8 {
				t.Fatalf("MQTT finding = %#v", finding)
			}
			if finding.Message != i18n.T(lang, "mqtt_anonymous", "Kitchen hub") {
				t.Fatalf("MQTT message = %q, want %q", finding.Message, i18n.T(lang, "mqtt_anonymous", "Kitchen hub"))
			}
		})
	}
}
