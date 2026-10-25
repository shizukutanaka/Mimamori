package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/nodered"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedNodeRED(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{1880}, nodered.Unauthenticated)
}
