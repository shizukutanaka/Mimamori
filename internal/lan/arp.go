package lan

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Neighbor struct {
	IP  netip.Addr
	MAC string
}

var (
	ipv4Pattern = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	macPattern  = regexp.MustCompile(`(?i)\b([0-9a-f]{1,2}[:-]){5}[0-9a-f]{1,2}\b`)
)

func ParseARP(text string) []Neighbor {
	neighbors := make(map[netip.Addr]string)
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(strings.ToLower(line), "incomplete") {
			continue
		}
		if fields := strings.Fields(line); len(fields) >= 6 && strings.HasPrefix(fields[2], "0x") {
			if flags, err := strconv.ParseUint(fields[2][2:], 16, 32); err != nil || flags&0x2 == 0 {
				continue
			}
		}

		ipText := ipv4Pattern.FindString(line)
		macText := macPattern.FindString(line)
		if ipText == "" || macText == "" {
			continue
		}
		ip, err := netip.ParseAddr(ipText)
		if err != nil || !ip.Is4() || !ip.IsPrivate() {
			continue
		}

		mac, ok := normalizeMAC(macText)
		if !ok {
			continue
		}
		if _, exists := neighbors[ip]; !exists {
			neighbors[ip] = mac
		}
	}

	result := make([]Neighbor, 0, len(neighbors))
	for ip, mac := range neighbors {
		result = append(result, Neighbor{IP: ip, MAC: mac})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].IP.Less(result[j].IP)
	})
	return result
}

func normalizeMAC(mac string) (string, bool) {
	octets := strings.FieldsFunc(mac, func(r rune) bool { return r == ':' || r == '-' })
	if len(octets) != 6 {
		return "", false
	}
	values := make([]byte, len(octets))
	for i, octet := range octets {
		value, err := strconv.ParseUint(octet, 16, 8)
		if err != nil {
			return "", false
		}
		values[i] = byte(value)
	}
	if values[0]&1 != 0 {
		return "", false
	}
	if values[0] == 0 && values[1] == 0 && values[2] == 0 && values[3] == 0 && values[4] == 0 && values[5] == 0 {
		return "", false
	}
	if values[0] == 0xff && values[1] == 0xff && values[2] == 0xff && values[3] == 0xff && values[4] == 0xff && values[5] == 0xff {
		return "", false
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", values[0], values[1], values[2], values[3], values[4], values[5]), true
}

func IsRandomizedMAC(mac string) bool {
	octets := strings.FieldsFunc(mac, func(r rune) bool { return r == ':' || r == '-' })
	if len(octets) != 6 {
		return false
	}
	first, err := strconv.ParseUint(octets[0], 16, 8)
	return err == nil && byte(first)&2 != 0
}

const commandTimeout = 15 * time.Second

func ReadARP(ctx context.Context) ([]Neighbor, error) {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/net/arp")
		if err != nil {
			return nil, err
		}
		return withoutAddresses(ParseARP(string(data)), LocalAddresses()), nil
	}

	name, args := "arp", []string{"-an"}
	if runtime.GOOS == "windows" {
		args = []string{"-a"}
	}
	commandContext, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()
	output, err := exec.CommandContext(commandContext, name, args...).Output()
	if err != nil {
		return nil, fmt.Errorf("arp table: %w", err)
	}
	return withoutAddresses(ParseARP(string(output)), LocalAddresses()), nil
}

// LocalAddresses returns the IP addresses assigned to this computer's interfaces.
func LocalAddresses() map[netip.Addr]struct{} {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	local := make(map[netip.Addr]struct{}, len(addresses))
	for _, address := range addresses {
		var ip net.IP
		switch address := address.(type) {
		case *net.IPNet:
			ip = address.IP
		case *net.IPAddr:
			ip = address.IP
		default:
			continue
		}

		if addr, ok := netip.AddrFromSlice(ip); ok {
			local[addr.Unmap()] = struct{}{}
		}
	}
	return local
}

func withoutAddresses(neighbors []Neighbor, excluded map[netip.Addr]struct{}) []Neighbor {
	if len(excluded) == 0 {
		return neighbors
	}

	filtered := make([]Neighbor, 0, len(neighbors))
	for _, neighbor := range neighbors {
		if _, skip := excluded[neighbor.IP]; !skip {
			filtered = append(filtered, neighbor)
		}
	}
	return filtered
}
