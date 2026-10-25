package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/esphome"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedESPHome(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{6052}, esphome.Unauthenticated)
}
