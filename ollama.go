package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/ollama"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markOpenOllama(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{11434}, ollama.Open)
}
