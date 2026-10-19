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

func homeNeighbors(neighbors []lan.Neighbor, prefixes []netip.Prefix, gateway netip.Addr) []lan.Neighbor {
	if !gateway.IsValid() || !gateway.Is4() {
		return neighbors
	}

	gatewayInNeighbors := false
	for _, neighbor := range neighbors {
		if neighbor.IP == gateway {
			gatewayInNeighbors = true
			break
		}
	}
	if !gatewayInNeighbors {
		return neighbors
	}

	matchingPrefixes := homePrefixes(prefixes, gateway)
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
