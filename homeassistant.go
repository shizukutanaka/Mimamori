package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/homeassistant"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markOpenHomeAssistantOnboarding(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{8123}, homeassistant.UserStepOpen)
}
