package main

import (
	"context"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/tasmota"
)

func markUnauthenticatedTasmota(ctx context.Context, opened []probe.Open) []probe.Open {
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
	results := make(chan int, len(candidates))
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			for index := range jobs {
				candidate := opened[index]
				if tasmota.Unauthenticated(ctx, candidate.IP, candidate.Port, 3*time.Second) {
					results <- index
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
	unauthenticated := make([]bool, len(opened))
	for index := range results {
		unauthenticated[index] = true
	}
	filtered := make([]probe.Open, 0, len(opened))
	for index, candidate := range opened {
		if candidate.Port != 80 || unauthenticated[index] {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}
