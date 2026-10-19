package probe

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Open struct {
	IP      netip.Addr
	Port    int
	Service string
}

type target struct {
	ip   netip.Addr
	port int
}

type dialFunc func(ctx context.Context, address string) (net.Conn, error)

type scanResult struct {
	target   target
	open     bool
	failed   bool
	timedOut bool
}

func Scan(ctx context.Context, hosts []netip.Addr, ports []int, timeout time.Duration, workers int) []Open {
	dialer := net.Dialer{Timeout: timeout}
	dial := func(ctx context.Context, address string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp", address)
	}
	return scan(ctx, hosts, ports, timeout, workers, dial)
}

// Alive dials each host once on port and returns the hosts that accepted or actively refused the connection.
func Alive(ctx context.Context, hosts []netip.Addr, port int, timeout time.Duration, workers int) map[netip.Addr]struct{} {
	dialer := net.Dialer{Timeout: timeout}
	dial := func(ctx context.Context, address string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp", address)
	}
	return alive(ctx, hosts, port, timeout, workers, dial)
}

func scan(ctx context.Context, hosts []netip.Addr, ports []int, timeout time.Duration, workers int, dial dialFunc) []Open {
	if workers < 1 {
		workers = 1
	}
	targets := make([]target, 0, len(hosts)*len(ports))
	for _, host := range hosts {
		for _, port := range ports {
			targets = append(targets, target{ip: host, port: port})
		}
	}

	firstPass := scanPass(ctx, targets, workers, dial)
	aliveHosts := make(map[netip.Addr]struct{})
	for _, result := range firstPass {
		if result.open || (result.failed && !result.timedOut) {
			aliveHosts[result.target.ip] = struct{}{}
		}
	}

	opened := make([]Open, 0)
	seen := make(map[target]struct{})
	appendOpens := func(results []scanResult) {
		for _, result := range results {
			if !result.open {
				continue
			}
			if _, exists := seen[result.target]; exists {
				continue
			}
			seen[result.target] = struct{}{}
			opened = append(opened, Open{IP: result.target.ip, Port: result.target.port})
		}
	}
	appendOpens(firstPass)

	if ctx.Err() == nil {
		retries := make([]target, 0)
		for _, result := range firstPass {
			if !result.timedOut {
				continue
			}
			if _, alive := aliveHosts[result.target.ip]; alive {
				retries = append(retries, result.target)
			}
		}
		appendOpens(scanPass(ctx, retries, workers, dial))
	}

	sort.Slice(opened, func(i, j int) bool {
		if opened[i].IP != opened[j].IP {
			return opened[i].IP.Less(opened[j].IP)
		}
		return opened[i].Port < opened[j].Port
	})
	return opened
}

func alive(ctx context.Context, hosts []netip.Addr, port int, _ time.Duration, workers int, dial dialFunc) map[netip.Addr]struct{} {
	if workers < 1 {
		workers = 1
	}
	targets := make([]target, 0, len(hosts))
	for _, host := range hosts {
		targets = append(targets, target{ip: host, port: port})
	}

	aliveHosts := make(map[netip.Addr]struct{})
	for _, result := range scanPass(ctx, targets, workers, dial) {
		if result.open || (result.failed && !result.timedOut) {
			aliveHosts[result.target.ip] = struct{}{}
		}
	}
	return aliveHosts
}

func scanPass(ctx context.Context, targets []target, workers int, dial dialFunc) []scanResult {
	results := make([]scanResult, len(targets))
	for index, job := range targets {
		results[index].target = job
	}
	jobs := make(chan int)
	var wait sync.WaitGroup
	wait.Add(workers)
	for range workers {
		go func() {
			defer wait.Done()
			for index := range jobs {
				if ctx.Err() != nil {
					continue
				}
				job := targets[index]
				address := net.JoinHostPort(job.ip.String(), strconv.Itoa(job.port))
				connection, err := dial(ctx, address)
				result := scanResult{target: job}
				if err != nil {
					result.failed = true
					var netErr net.Error
					result.timedOut = errors.As(err, &netErr) && netErr.Timeout()
				} else {
					if connection != nil {
						_ = connection.Close()
					}
					result.open = true
				}
				results[index] = result
			}
		}()
	}
	go func() {
		defer close(jobs)
		for index := range targets {
			select {
			case jobs <- index:
			case <-ctx.Done():
				return
			}
		}
	}()
	wait.Wait()
	return results
}
