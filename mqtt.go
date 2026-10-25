package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/mqtt"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedMQTT(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{1883}, mqtt.Anonymous)
}
