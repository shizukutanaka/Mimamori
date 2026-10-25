package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/rtsp"
)

func markUnauthenticatedRTSP(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{554, 8554}, rtsp.Unauthenticated)
}
