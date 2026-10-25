package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/frigate"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedFrigate(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{5000}, frigate.Unauthenticated)
}
