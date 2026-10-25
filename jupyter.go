package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/jupyter"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedJupyter(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{8888}, jupyter.Unauthenticated)
}
