package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/go2rtc"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedGo2RTC(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{1984}, go2rtc.Unauthenticated)
}
