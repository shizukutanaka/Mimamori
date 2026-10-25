package main

import (
	"context"
	"net/netip"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/docker"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedDocker(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{2375, 2376}, func(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
		if port == 2376 {
			return docker.UnauthenticatedTLS(ctx, ip, port, timeout)
		}
		return docker.Unauthenticated(ctx, ip, port, timeout)
	})
}
