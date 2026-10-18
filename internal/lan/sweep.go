package lan

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sort"
)

func SweepHosts(prefixes []netip.Prefix, gateway netip.Addr) []netip.Addr {
	if !gateway.IsValid() || !gateway.Is4() || !gateway.IsPrivate() {
		return nil
	}

	selfAddresses := make(map[netip.Addr]struct{}, len(prefixes))
	for _, prefix := range prefixes {
		if prefix.IsValid() && prefix.Addr().Is4() {
			selfAddresses[prefix.Addr()] = struct{}{}
		}
	}

	hosts := make(map[netip.Addr]struct{})
	for _, prefix := range prefixes {
		if !prefix.IsValid() || !prefix.Addr().Is4() || !prefix.Addr().IsPrivate() ||
			!prefix.Masked().Contains(gateway) {
			continue
		}

		originalSubnet := prefix.Masked()
		subnet := originalSubnet
		capped := subnet.Bits() < 24
		if capped {
			subnet = netip.PrefixFrom(prefix.Addr(), 24).Masked()
			if !subnet.Contains(gateway) {
				continue
			}
		} else if subnet.Bits() > 30 {
			continue
		}

		network := binary.BigEndian.Uint32(subnet.Addr().AsSlice())
		mask := ^uint32(0) << (32 - subnet.Bits())
		broadcast := network | ^mask
		originalNetwork := binary.BigEndian.Uint32(originalSubnet.Addr().AsSlice())
		originalMask := ^uint32(0) << (32 - originalSubnet.Bits())
		originalBroadcast := originalNetwork | ^originalMask
		for value := network; ; value++ {
			skip := false
			if capped {
				skip = value == originalNetwork || value == originalBroadcast
			} else {
				skip = value == network || value == broadcast
			}
			address := netip.AddrFrom4([4]byte{
				byte(value >> 24),
				byte(value >> 16),
				byte(value >> 8),
				byte(value),
			})
			if !skip {
				if _, isSelf := selfAddresses[address]; !isSelf {
					hosts[address] = struct{}{}
				}
			}
			if value == broadcast {
				break
			}
		}
	}

	result := make([]netip.Addr, 0, len(hosts))
	for host := range hosts {
		result = append(result, host)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Less(result[j])
	})
	return result
}

func LocalPrefixes() ([]netip.Prefix, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list network interfaces: %w", err)
	}

	prefixSet := make(map[netip.Prefix]struct{})
	var prefixes []netip.Prefix
	var addressErrors []error
	for _, iface := range interfaces {
		if !sweepableInterface(iface.Flags) {
			continue
		}
		addresses, err := iface.Addrs()
		if err != nil {
			addressErrors = append(addressErrors, fmt.Errorf("%s: %w", iface.Name, err))
			continue
		}
		for _, address := range addresses {
			ipNet, ok := address.(*net.IPNet)
			if !ok || ipNet == nil {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}
			ones, bits := ipNet.Mask.Size()
			if bits != 32 || ones < 0 {
				continue
			}
			addr, ok := netip.AddrFromSlice(ip)
			if !ok {
				continue
			}
			prefix := netip.PrefixFrom(addr, ones)
			if !prefix.IsValid() {
				continue
			}
			if _, exists := prefixSet[prefix]; exists {
				continue
			}
			prefixSet[prefix] = struct{}{}
			prefixes = append(prefixes, prefix)
		}
	}
	sort.Slice(prefixes, func(i, j int) bool {
		if prefixes[i].Addr() != prefixes[j].Addr() {
			return prefixes[i].Addr().Less(prefixes[j].Addr())
		}
		return prefixes[i].Bits() < prefixes[j].Bits()
	})
	return prefixes, errors.Join(addressErrors...)
}

// sweepableInterface reports whether iface is a broadcast LAN interface (Ethernet/Wi-Fi), not loopback or a VPN tunnel.
func sweepableInterface(flags net.Flags) bool {
	return flags&net.FlagUp != 0 && flags&net.FlagLoopback == 0 &&
		flags&net.FlagBroadcast != 0 && flags&net.FlagPointToPoint == 0
}
