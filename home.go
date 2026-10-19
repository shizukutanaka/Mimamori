package main

import (
	"net/netip"

	"github.com/shizukutanaka/Mimamori/internal/lan"
)

func homePrefixes(prefixes []netip.Prefix, gateway netip.Addr) []netip.Prefix {
	if !gateway.IsValid() || !gateway.Is4() {
		return nil
	}

	mostSpecificBits := -1
	for _, prefix := range prefixes {
		if prefix.Contains(gateway) && prefix.Bits() > mostSpecificBits {
			mostSpecificBits = prefix.Bits()
		}
	}
	if mostSpecificBits < 0 {
		return nil
	}

	var matchingPrefixes []netip.Prefix
	for _, prefix := range prefixes {
		if prefix.Contains(gateway) && prefix.Bits() == mostSpecificBits {
			matchingPrefixes = append(matchingPrefixes, prefix)
		}
	}
	return matchingPrefixes
}

func scanScope(neighbors []lan.Neighbor, prefixes []netip.Prefix, gateway netip.Addr) []netip.Prefix {
	if !gatewayInNeighbors(neighbors, gateway) {
		return nil
	}
	return homePrefixes(prefixes, gateway)
}

func scanSweepHosts(prefixes, scope []netip.Prefix, gateway netip.Addr) []netip.Addr {
	if len(scope) > 0 {
		prefixes = scope
	}
	return lan.SweepHosts(prefixes, gateway)
}

func gatewayInNeighbors(neighbors []lan.Neighbor, gateway netip.Addr) bool {
	if !gateway.IsValid() || !gateway.Is4() {
		return false
	}
	for _, neighbor := range neighbors {
		if neighbor.IP == gateway {
			return true
		}
	}
	return false
}

func homeNeighbors(neighbors []lan.Neighbor, prefixes []netip.Prefix, gateway netip.Addr) []lan.Neighbor {
	matchingPrefixes := scanScope(neighbors, prefixes, gateway)
	if len(matchingPrefixes) == 0 {
		return neighbors
	}

	filtered := make([]lan.Neighbor, 0, len(neighbors))
	for _, neighbor := range neighbors {
		for _, prefix := range matchingPrefixes {
			if prefix.Contains(neighbor.IP) {
				filtered = append(filtered, neighbor)
				break
			}
		}
	}
	return filtered
}
