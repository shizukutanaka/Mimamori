package main

import (
	"context"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/smb"
)

const smb1Port = -445

func markSMB1(ctx context.Context, opened []probe.Open) []probe.Open {
	if ctx == nil {
		return opened
	}

	candidates := make([]int, 0)
	for index, candidate := range opened {
		if candidate.Port == 445 {
			candidates = append(candidates, index)
		}
	}
	if len(candidates) == 0 {
		return opened
	}

	workers := 16
	if len(candidates) < workers {
		workers = len(candidates)
	}
	type result struct {
		index   int
		enabled bool
	}
	jobs := make(chan int)
	results := make(chan result, len(candidates))
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			for index := range jobs {
				candidate := opened[index]
				results <- result{
					index:   index,
					enabled: smb.SMB1Enabled(ctx, candidate.IP, candidate.Port, 2*time.Second),
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
	enabled := make([]bool, len(opened))
	for result := range results {
		enabled[result.index] = result.enabled
	}

	marked := make([]probe.Open, 0, len(opened))
	for index, candidate := range opened {
		if candidate.Port != 445 {
			marked = append(marked, candidate)
		} else if enabled[index] {
			candidate.Port = smb1Port
			marked = append(marked, candidate)
		}
	}
	return marked
}
