package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/syncthing"
)

func markUnauthenticatedSyncthing(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{8384}, syncthing.Unauthenticated)
}
