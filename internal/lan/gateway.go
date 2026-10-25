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

var windowsRoutePattern = regexp.MustCompile(`(?m)^\s*0\.0\.0\.0\s+0\.0\.0\.0\s+(\d+\.\d+\.\d+\.\d+)\s+\S+\s+(\d+)\s*$`)

func ParseLinuxRoute(text string) (netip.Addr, bool) {
	var bestGateway netip.Addr
	var bestMetric uint64
	found := false
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 8 || fields[1] != "00000000" || fields[7] != "00000000" || len(fields[2]) != 8 {
			continue
		}
		flags, err := strconv.ParseUint(fields[3], 16, 32)
		if err != nil || flags&0x1 == 0 || flags&0x2 == 0 {
			continue
		}
		metric, err := strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
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
		if gateway.IsUnspecified() {
			continue
		}
		if !found || metric < bestMetric {
			bestGateway = gateway
			bestMetric = metric
			found = true
		}
	}
	return bestGateway, found
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
	var bestGateway netip.Addr
	var bestMetric uint64
	found := false
	for _, match := range windowsRoutePattern.FindAllStringSubmatch(text, -1) {
		gateway, err := netip.ParseAddr(match[1])
		if err != nil || !gateway.Is4() || gateway.IsUnspecified() {
			continue
		}
		metric, err := strconv.ParseUint(match[2], 10, 64)
		if err != nil {
			continue
		}
		if !found || metric < bestMetric {
			bestGateway = gateway
			bestMetric = metric
			found = true
		}
	}
	return bestGateway, found
}

func DefaultGateway(ctx context.Context) (netip.Addr, error) {
	commandContext, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()
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
		output, err = exec.CommandContext(commandContext, "route", "print", "-4", "0.0.0.0").Output()
		if err == nil {
			if gateway, ok := ParseWindowsRoute(string(output)); ok {
				return gateway, nil
			}
			return netip.Addr{}, fmt.Errorf("default route not found")
		}
	default:
		output, err = exec.CommandContext(commandContext, "route", "-n", "get", "default").Output()
		if err == nil {
			if gateway, ok := ParseBSDRoute(string(output)); ok {
				return gateway, nil
			}
			return netip.Addr{}, fmt.Errorf("default route not found")
		}
	}
	return netip.Addr{}, fmt.Errorf("read default route: %w", err)
}
