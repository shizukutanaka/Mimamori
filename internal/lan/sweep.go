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

		subnet := prefix.Masked()
		if subnet.Bits() < 24 {
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
		for value := network + 1; value < broadcast; value++ {
			address := netip.AddrFrom4([4]byte{
				byte(value >> 24),
				byte(value >> 16),
				byte(value >> 8),
				byte(value),
			})
			if _, isSelf := selfAddresses[address]; isSelf {
				continue
			}
			hosts[address] = struct{}{}
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
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
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
