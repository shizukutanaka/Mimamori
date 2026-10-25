package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/zigbee2mqtt"
)

func markUnauthenticatedZigbee2MQTT(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{8080}, zigbee2mqtt.Unauthenticated)
}
