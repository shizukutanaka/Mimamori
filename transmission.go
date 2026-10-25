package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/transmission"
)

func markUnauthenticatedTransmission(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{9091}, transmission.Unauthenticated)
}
