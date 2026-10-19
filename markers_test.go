package main

import (
	"context"
	"math/rand"
	"net/netip"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func TestRunMarkersMatchesSequentialFold(t *testing.T) {
	markers := []portMarker{
		testPortMarker([]int{100, 101}, func(candidate probe.Open) (probe.Open, bool) {
			if candidate.Port == 100 && candidate.IP.As4()[3]%2 == 0 {
				return candidate, false
			}
			if candidate.Port == 101 {
				candidate.Port = -candidate.Port
			} else {
				candidate.Service = "first"
			}
			return candidate, true
		}),
		testPortMarker([]int{200}, func(candidate probe.Open) (probe.Open, bool) {
			if candidate.IP.As4()[3]%3 == 0 {
				return candidate, false
			}
			candidate.Service = "second"
			return candidate, true
		}),
		testPortMarker([]int{300, 301}, func(candidate probe.Open) (probe.Open, bool) {
			if candidate.Port == 300 {
				candidate.Port = -candidate.Port
			} else {
				candidate.Service = "third"
			}
			return candidate, true
		}),
		testPortMarker([]int{400}, func(candidate probe.Open) (probe.Open, bool) {
			candidate.Service = "fourth"
			return candidate, candidate.IP.As4()[3]%2 == 1
		}),
	}
	ips := []netip.Addr{
		netip.MustParseAddr("192.168.1.10"),
		netip.MustParseAddr("192.168.1.11"),
		netip.MustParseAddr("192.168.1.12"),
		netip.MustParseAddr("192.168.1.13"),
	}
	ports := []int{21, 80, 100, 101, 200, 300, 301, 400, 5900}
	random := rand.New(rand.NewSource(20260415))

	for iteration := range 500 {
		opened := make([]probe.Open, 0, len(ips)*len(ports))
		for _, ip := range ips {
			for _, port := range ports {
				if random.Intn(3) == 0 {
					opened = append(opened, probe.Open{IP: ip, Port: port})
				}
			}
		}
		if len(opened) == 0 {
			opened = append(opened, probe.Open{IP: ips[0], Port: 21})
		}

		got := runMarkers(context.Background(), opened, markers)
		want := sequentialFold(context.Background(), opened, markers)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("iteration %d: runMarkers() = %#v, sequential fold = %#v", iteration, got, want)
		}
	}
}

func TestRunMarkersConcurrent(t *testing.T) {
	markers := make([]portMarker, 5)
	opened := make([]probe.Open, len(markers))
	var calls atomic.Int32
	for i := range markers {
		port := i + 1
		markers[i] = portMarker{
			ports: []int{port},
			mark: func(_ context.Context, candidates []probe.Open) []probe.Open {
				calls.Add(1)
				time.Sleep(200 * time.Millisecond)
				return candidates
			},
		}
		opened[i] = probe.Open{IP: netip.MustParseAddr("192.168.1.10"), Port: port}
	}

	started := time.Now()
	runMarkers(context.Background(), opened, markers)
	if elapsed := time.Since(started); elapsed >= 600*time.Millisecond {
		t.Fatalf("runMarkers() took %s, want less than 600ms", elapsed)
	}
	if got := calls.Load(); got != int32(len(markers)) {
		t.Fatalf("runMarkers() called %d markers, want %d", got, len(markers))
	}
}

func TestServiceMarkers(t *testing.T) {
	markerPorts := make(map[int]struct{})
	for _, marker := range serviceMarkers {
		for _, port := range marker.ports {
			if _, exists := markerPorts[port]; exists {
				t.Fatalf("serviceMarkers contains duplicate port %d", port)
			}
			markerPorts[port] = struct{}{}
		}
	}

	scannedPorts := make(map[int]struct{}, len(scanPorts))
	for _, port := range scanPorts {
		if _, exists := scannedPorts[port]; exists {
			t.Fatalf("scanPorts contains duplicate port %d", port)
		}
		scannedPorts[port] = struct{}{}
	}
	if !reflect.DeepEqual(markerPorts, scannedPorts) {
		t.Fatalf("service marker ports = %v, scan ports = %v", markerPorts, scannedPorts)
	}
}

func TestRunMarkersPanicsOnDuplicatePortClaims(t *testing.T) {
	markers := []portMarker{{ports: []int{1}}, {ports: []int{1}}}
	defer func() {
		if recover() == nil {
			t.Fatal("runMarkers() did not panic for duplicate port claims")
		}
	}()
	runMarkers(context.Background(), nil, markers)
}

func testPortMarker(ports []int, transform func(probe.Open) (probe.Open, bool)) portMarker {
	owned := make(map[int]struct{}, len(ports))
	for _, port := range ports {
		owned[port] = struct{}{}
	}
	return portMarker{
		ports: ports,
		mark: func(_ context.Context, opened []probe.Open) []probe.Open {
			marked := make([]probe.Open, 0, len(opened))
			for _, candidate := range opened {
				if _, isOwned := owned[candidate.Port]; !isOwned {
					marked = append(marked, candidate)
					continue
				}
				changed, keep := transform(candidate)
				if keep {
					marked = append(marked, changed)
				}
			}
			return marked
		},
	}
}

func sequentialFold(ctx context.Context, opened []probe.Open, markers []portMarker) []probe.Open {
	result := opened
	for _, marker := range markers {
		result = marker.mark(ctx, result)
	}
	return result
}
