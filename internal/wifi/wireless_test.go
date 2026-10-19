package wifi

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWirelessConnected(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		marker        string
		operstate     string
		missingRoot   bool
		want          bool
	}{
		{name: "phy80211 up", interfaceName: "wlan0", marker: "phy80211", operstate: "up\n", want: true},
		{name: "wireless up", interfaceName: "wlan0", marker: "wireless", operstate: "up", want: true},
		{name: "dormant", interfaceName: "wlan0", marker: "phy80211", operstate: "dormant"},
		{name: "down", interfaceName: "wlan0", marker: "phy80211", operstate: "down"},
		{name: "wired interface", interfaceName: "eth0", operstate: "up"},
		{name: "missing root", missingRoot: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if test.missingRoot {
				root = filepath.Join(root, "missing")
			} else {
				makeWirelessFixture(t, root, test.interfaceName, test.marker, test.operstate)
			}
			if got := wirelessConnected(root); got != test.want {
				t.Fatalf("wirelessConnected() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnreadableIfConnected(t *testing.T) {
	connectedRoot := t.TempDir()
	makeWirelessFixture(t, connectedRoot, "wlan0", "phy80211", "up\n")
	wiredRoot := t.TempDir()
	makeWirelessFixture(t, wiredRoot, "eth0", "", "up")

	tests := []struct {
		name         string
		root         string
		security     Security
		err          error
		wantSecurity Security
		wantErr      error
	}{
		{name: "connected but unknown", root: connectedRoot, security: Unknown, wantSecurity: Unknown, wantErr: ErrWiFiUnreadable},
		{name: "connected modern", root: connectedRoot, security: Modern, wantSecurity: Modern},
		{name: "connected open", root: connectedRoot, security: Open, wantSecurity: Open},
		{name: "connected with error", root: connectedRoot, security: Unknown, err: context.Canceled, wantSecurity: Unknown, wantErr: context.Canceled},
		{name: "wired only", root: wiredRoot, security: Unknown, wantSecurity: Unknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotSecurity, gotErr := unreadableIfConnected(test.root, test.security, test.err)
			if gotSecurity != test.wantSecurity {
				t.Errorf("unreadableIfConnected() security = %v, want %v", gotSecurity, test.wantSecurity)
			}
			if !errors.Is(gotErr, test.wantErr) {
				t.Errorf("unreadableIfConnected() error = %v, want %v", gotErr, test.wantErr)
			}
		})
	}
}

func makeWirelessFixture(t *testing.T, root, interfaceName, marker, operstate string) {
	t.Helper()
	interfacePath := filepath.Join(root, interfaceName)
	if err := os.MkdirAll(interfacePath, 0o755); err != nil {
		t.Fatal(err)
	}
	if marker != "" {
		if err := os.MkdirAll(filepath.Join(interfacePath, marker), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(interfacePath, "operstate"), []byte(operstate), 0o644); err != nil {
		t.Fatal(err)
	}
}
