package wifi

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrWiFiUnreadable = errors.New("a connected Wi-Fi interface was found but its security could not be read")

// wirelessConnected reports whether any interface under root (normally /sys/class/net) is a Wi-Fi interface whose operstate is "up".
func wirelessConnected(root string) bool {
	interfaces, err := os.ReadDir(root)
	if err != nil {
		return false
	}
	for _, iface := range interfaces {
		ifacePath := filepath.Join(root, iface.Name())
		if _, err := os.Stat(filepath.Join(ifacePath, "phy80211")); err != nil {
			if _, err := os.Stat(filepath.Join(ifacePath, "wireless")); err != nil {
				continue
			}
		}
		operstate, err := os.ReadFile(filepath.Join(ifacePath, "operstate"))
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(operstate)) == "up" {
			return true
		}
	}
	return false
}

func unreadableIfConnected(root string, security Security, err error) (Security, error) {
	if err == nil && security == Unknown && wirelessConnected(root) {
		return Unknown, ErrWiFiUnreadable
	}
	return security, err
}
