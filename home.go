package main

import (
	"net/netip"

	"github.com/shizukutanaka/Mimamori/internal/lan"
)

func homeNeighbors(neighbors []lan.Neighbor, prefixes []netip.Prefix, gateway netip.Addr) []lan.Neighbor {
	if !gateway.IsValid() || !gateway.Is4() {
		return neighbors
	}

	var homePrefixes []netip.Prefix
	for _, prefix := range prefixes {
		if prefix.Contains(gateway) {
			homePrefixes = append(homePrefixes, prefix)
		}
	}
	if len(homePrefixes) == 0 {
		return neighbors
	}

	filtered := make([]lan.Neighbor, 0, len(neighbors))
	for _, neighbor := range neighbors {
		for _, prefix := range homePrefixes {
			if prefix.Contains(neighbor.IP) {
				filtered = append(filtered, neighbor)
				break
			}
		}
	}
	return filtered
}
