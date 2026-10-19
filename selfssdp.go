package main

import (
	"net/netip"

	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func withoutLocalResponses(responses []ssdp.Response, local map[netip.Addr]struct{}) []ssdp.Response {
	if len(local) == 0 {
		return responses
	}
	filtered := make([]ssdp.Response, 0, len(responses))
	for _, response := range responses {
		if _, isLocal := local[response.IP]; !isLocal {
			filtered = append(filtered, response)
		}
	}
	return filtered
}
