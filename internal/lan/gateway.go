package lan

import (
	"context"
	"fmt"
	"net/netip"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var windowsRoutePattern = regexp.MustCompile(`(?m)^\s*0\.0\.0\.0\s+0\.0\.0\.0\s+(\d+\.\d+\.\d+\.\d+)`)

func ParseLinuxRoute(text string) (netip.Addr, bool) {
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[1] != "00000000" || len(fields[2]) != 8 {
			continue
		}
		value, err := strconv.ParseUint(fields[2], 16, 32)
		if err != nil {
			continue
		}
		gateway := netip.AddrFrom4([4]byte{
			byte(value),
			byte(value >> 8),
			byte(value >> 16),
			byte(value >> 24),
		})
		if !gateway.IsUnspecified() {
			return gateway, true
		}
	}
	return netip.Addr{}, false
}

func ParseBSDRoute(text string) (netip.Addr, bool) {
	for _, line := range strings.Split(text, "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(fields) != 2 || !strings.EqualFold(strings.TrimSpace(fields[0]), "gateway") {
			continue
		}
		gateway, err := netip.ParseAddr(strings.TrimSpace(fields[1]))
		if err != nil || !gateway.Is4() {
			return netip.Addr{}, false
		}
		return gateway, true
	}
	return netip.Addr{}, false
}

func ParseWindowsRoute(text string) (netip.Addr, bool) {
	for _, match := range windowsRoutePattern.FindAllStringSubmatch(text, -1) {
		gateway, err := netip.ParseAddr(match[1])
		if err == nil && gateway.Is4() && !gateway.IsUnspecified() {
			return gateway, true
		}
	}
	return netip.Addr{}, false
}

func DefaultGateway(ctx context.Context) (netip.Addr, error) {
	var output []byte
	var err error
	switch runtime.GOOS {
	case "linux":
		output, err = os.ReadFile("/proc/net/route")
		if err == nil {
			if gateway, ok := ParseLinuxRoute(string(output)); ok {
				return gateway, nil
			}
			return netip.Addr{}, fmt.Errorf("default route not found")
		}
	case "windows":
		output, err = exec.CommandContext(ctx, "route", "print", "-4", "0.0.0.0").Output()
		if err == nil {
			if gateway, ok := ParseWindowsRoute(string(output)); ok {
				return gateway, nil
			}
			return netip.Addr{}, fmt.Errorf("default route not found")
		}
	default:
		output, err = exec.CommandContext(ctx, "route", "-n", "get", "default").Output()
		if err == nil {
			if gateway, ok := ParseBSDRoute(string(output)); ok {
				return gateway, nil
			}
			return netip.Addr{}, fmt.Errorf("default route not found")
		}
	}
	return netip.Addr{}, fmt.Errorf("read default route: %w", err)
}
