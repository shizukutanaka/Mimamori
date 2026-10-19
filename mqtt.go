package main

import (
	"context"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/mqtt"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func markUnauthenticatedMQTT(ctx context.Context, opened []probe.Open) []probe.Open {
	candidates := make([]int, 0, len(opened))
	for index, candidate := range opened {
		if candidate.Port == 1883 {
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
				if mqtt.Anonymous(ctx, candidate.IP, candidate.Port, 2*time.Second) {
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
	anonymous := make([]bool, len(opened))
	for index := range results {
		anonymous[index] = true
	}
	filtered := make([]probe.Open, 0, len(opened))
	for index, candidate := range opened {
		if candidate.Port != 1883 || anonymous[index] {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}
