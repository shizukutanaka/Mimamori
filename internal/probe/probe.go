package probe

import (
	"context"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Open struct {
	IP   netip.Addr
	Port int
}

type target struct {
	ip   netip.Addr
	port int
}

func Scan(ctx context.Context, hosts []netip.Addr, ports []int, timeout time.Duration, workers int) []Open {
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan target)
	results := make(chan Open)
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			dialer := net.Dialer{Timeout: timeout}
			for job := range jobs {
				if ctx.Err() != nil {
					continue
				}
				address := net.JoinHostPort(job.ip.String(), strconv.Itoa(job.port))
				connection, err := dialer.DialContext(ctx, "tcp", address)
				if err != nil {
					continue
				}
				_ = connection.Close()
				select {
				case results <- Open{IP: job.ip, Port: job.port}:
				case <-ctx.Done():
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, host := range hosts {
			for _, port := range ports {
				select {
				case jobs <- target{ip: host, port: port}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	go func() {
		wait.Wait()
		close(results)
	}()

	opened := make([]Open, 0)
	for result := range results {
		opened = append(opened, result)
	}
	sort.Slice(opened, func(i, j int) bool {
		if opened[i].IP != opened[j].IP {
			return opened[i].IP.Less(opened[j].IP)
		}
		return opened[i].Port < opened[j].Port
	})
	return opened
}
