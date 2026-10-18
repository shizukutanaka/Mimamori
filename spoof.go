package main

import (
	"context"
	"net/netip"
	"sort"
	"strconv"
	"strings"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/lan"
	"github.com/shizukutanaka/Mimamori/internal/report"
)

func gatewayNeighbors(
	ctx context.Context,
	neighbors []lan.Neighbor,
	gateway netip.Addr,
	readARP func(context.Context) ([]lan.Neighbor, error),
) []lan.Neighbor {
	if !gateway.IsValid() {
		return neighbors
	}
	for _, neighbor := range neighbors {
		if neighbor.IP == gateway {
			return neighbors
		}
	}
	refreshed, err := readARP(ctx)
	if err != nil {
		return neighbors
	}
	return mergeNeighbors(neighbors, refreshed)
}

func gatewaySpoofFindings(neighbors []lan.Neighbor, gateway netip.Addr, lang string) []report.Finding {
	if !gateway.IsValid() {
		return nil
	}

	var gatewayMAC [6]byte
	foundGateway := false
	for _, neighbor := range neighbors {
		if neighbor.IP != gateway {
			continue
		}
		var ok bool
		gatewayMAC, ok = normalizeSpoofMAC(neighbor.MAC)
		if !ok {
			return nil
		}
		foundGateway = true
		break
	}
	if !foundGateway {
		return nil
	}

	otherIPs := make(map[netip.Addr]struct{})
	for _, neighbor := range neighbors {
		if neighbor.IP == gateway || !neighbor.IP.Is4() || !neighbor.IP.IsPrivate() {
			continue
		}
		mac, ok := normalizeSpoofMAC(neighbor.MAC)
		if ok && mac == gatewayMAC {
			otherIPs[neighbor.IP] = struct{}{}
		}
	}
	if len(otherIPs) == 0 {
		return nil
	}

	sortedIPs := make([]netip.Addr, 0, len(otherIPs))
	for ip := range otherIPs {
		sortedIPs = append(sortedIPs, ip)
	}
	sort.Slice(sortedIPs, func(i, j int) bool {
		return sortedIPs[i].Less(sortedIPs[j])
	})
	firstOtherIP := sortedIPs[0].String()
	return []report.Finding{{
		ID: "gateway_spoof", Level: report.Caution, Points: 8, IP: firstOtherIP,
		Message: i18n.T(lang, "gateway_spoof", firstOtherIP),
		Fix:     i18n.T(lang, "gateway_spoof_fix"),
	}}
}

func normalizeSpoofMAC(mac string) ([6]byte, bool) {
	var normalized [6]byte
	parts := strings.Split(strings.ToLower(strings.ReplaceAll(strings.TrimSpace(mac), "-", ":")), ":")
	if len(parts) != len(normalized) {
		return normalized, false
	}
	for i, part := range parts {
		if len(part) == 0 || len(part) > 2 {
			return normalized, false
		}
		if len(part) == 1 {
			part = "0" + part
		}
		value, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return normalized, false
		}
		normalized[i] = byte(value)
	}

	allZero, broadcast := true, true
	for _, octet := range normalized {
		allZero = allZero && octet == 0
		broadcast = broadcast && octet == 0xff
	}
	if allZero || broadcast || normalized[0]&1 != 0 {
		return [6]byte{}, false
	}
	return normalized, true
}
