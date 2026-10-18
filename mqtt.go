package main

import (
	"context"
	"net/netip"
	"sort"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/mqtt"
	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func unauthenticatedMQTT(ctx context.Context, hosts []netip.Addr) []probe.Open {
	openHosts := probe.Scan(ctx, hosts, []int{1883}, 700*time.Millisecond, 128)
	if len(openHosts) == 0 {
		return nil
	}

	workers := 16
	if len(openHosts) < workers {
		workers = len(openHosts)
	}
	jobs := make(chan probe.Open)
	results := make(chan probe.Open, len(openHosts))
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			for candidate := range jobs {
				if mqtt.Anonymous(ctx, candidate.IP, candidate.Port, 2*time.Second) {
					results <- probe.Open{IP: candidate.IP, Port: 1883}
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, candidate := range openHosts {
			select {
			case jobs <- candidate:
			case <-ctx.Done():
				return
			}
		}
	}()

	wait.Wait()
	close(results)
	opened := make([]probe.Open, 0, len(results))
	for result := range results {
		opened = append(opened, result)
	}
	sort.Slice(opened, func(i, j int) bool {
		return opened[i].IP.Less(opened[j].IP)
	})
	return opened
}
