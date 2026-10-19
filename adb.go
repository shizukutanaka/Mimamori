package main

import (
	"context"
	"net/netip"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/adb"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

const adbNoAuthPort = -5555

func markUnauthenticatedADB(ctx context.Context, opened []probe.Open) []probe.Open {
	if ctx == nil {
		return opened
	}

	candidates := make([]probe.Open, 0)
	for _, candidate := range opened {
		if candidate.Port == 5555 {
			candidates = append(candidates, candidate)
		}
	}
	if len(candidates) == 0 {
		return opened
	}

	workers := 16
	if len(candidates) < workers {
		workers = len(candidates)
	}
	jobs := make(chan probe.Open)
	results := make(chan netip.Addr, len(candidates))
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			for candidate := range jobs {
				if adb.NoAuth(ctx, candidate.IP, candidate.Port, 2*time.Second) {
					results <- candidate.IP
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, candidate := range candidates {
			select {
			case jobs <- candidate:
			case <-ctx.Done():
				return
			}
		}
	}()

	wait.Wait()
	close(results)
	unauthenticated := make(map[netip.Addr]struct{})
	for ip := range results {
		unauthenticated[ip] = struct{}{}
	}

	marked := make([]probe.Open, len(opened))
	copy(marked, opened)
	for i := range marked {
		if marked[i].Port == 5555 {
			if _, ok := unauthenticated[marked[i].IP]; ok {
				marked[i].Port = adbNoAuthPort
			}
		}
	}
	return marked
}
