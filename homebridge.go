package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/homebridge"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedHomebridge(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{8581}, homebridge.Unauthenticated)
}
