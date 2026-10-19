package main

import (
	"context"
	"net/netip"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/shelly"
	"github.com/shizukutanaka/Mimamori/internal/tasmota"
)

func markUnauthenticatedWeb(ctx context.Context, opened []probe.Open) []probe.Open {
	return markUnauthenticatedWebWithPort(ctx, opened, 80)
}

func markUnauthenticatedWebWithPort(ctx context.Context, opened []probe.Open, probePort int) []probe.Open {
	candidates := make([]int, 0, len(opened))
	for index, candidate := range opened {
		if candidate.Port == 80 {
			candidates = append(candidates, index)
		}
	}
	workers := 16
	if len(candidates) == 0 {
		return opened
	}
	if len(candidates) < workers {
		workers = len(candidates)
	}
	jobs := make(chan int)
	type result struct {
		index   int
		service string
	}
	results := make(chan result, len(candidates))
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			for index := range jobs {
				candidate := opened[index]
				if service := classifyUnauthenticatedWeb(ctx, candidate.IP, probePort); service != "" {
					results <- result{index: index, service: service}
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, index := range candidates {
			select {
			case jobs <- index:
			case <-ctx.Done():
				return
			}
		}
	}()

	wait.Wait()
	close(results)
	services := make([]string, len(opened))
	for result := range results {
		services[result.index] = result.service
	}
	filtered := make([]probe.Open, 0, len(opened))
	for index, candidate := range opened {
		if candidate.Port != 80 {
			filtered = append(filtered, candidate)
			continue
		}
		if services[index] != "" {
			candidate.Service = services[index]
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func classifyUnauthenticatedWeb(ctx context.Context, ip netip.Addr, port int) string {
	if tasmota.Unauthenticated(ctx, ip, port, 3*time.Second) {
		return "tasmota"
	}
	if shelly.Unauthenticated(ctx, ip, port, 3*time.Second) {
		return "shelly"
	}
	return ""
}
