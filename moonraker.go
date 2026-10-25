package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/moonraker"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedMoonraker(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{7125}, moonraker.Unauthenticated)
}
