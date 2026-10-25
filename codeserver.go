package main

import (
	"context"

	"github.com/shizukutanaka/Mimamori/internal/codeserver"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedCodeServer(ctx context.Context, opened []probe.Open) []probe.Open {
	return keepIfConfirmed(ctx, opened, []int{8443}, codeserver.Unauthenticated)
}
