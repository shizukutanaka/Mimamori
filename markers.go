package main

import (
	"context"
	"sort"
	"sync"

	"github.com/shizukutanaka/Mimamori/internal/probe"
)

type portMarker struct {
	ports []int
	mark  func(context.Context, []probe.Open) []probe.Open
}

var serviceMarkers = []portMarker{
	{ports: []int{23, 2323, 21}, mark: markUnchanged},
	{ports: []int{445}, mark: markSMB1},
	{ports: []int{554, 8554}, mark: markUnauthenticatedRTSP},
	{ports: []int{1883}, mark: markUnauthenticatedMQTT},
	{ports: []int{1984}, mark: markUnauthenticatedGo2RTC},
	{ports: []int{5000}, mark: markUnauthenticatedFrigate},
	{ports: []int{1880}, mark: markUnauthenticatedNodeRED},
	{ports: []int{8581}, mark: markUnauthenticatedHomebridge},
	{ports: []int{8384}, mark: markUnauthenticatedSyncthing},
	{ports: []int{9091}, mark: markUnauthenticatedTransmission},
	{ports: []int{8888}, mark: markUnauthenticatedJupyter},
	{ports: []int{11434}, mark: markOpenOllama},
	{ports: []int{8443}, mark: markUnauthenticatedCodeServer},
	{ports: []int{6052}, mark: markUnauthenticatedESPHome},
	{ports: []int{2375, 2376}, mark: markUnauthenticatedDocker},
	{ports: []int{8123}, mark: markOpenHomeAssistantOnboarding},
	{ports: []int{8080}, mark: markUnauthenticatedZigbee2MQTT},
	{ports: []int{7125}, mark: markUnauthenticatedMoonraker},
	{ports: []int{80}, mark: markUnauthenticatedWeb},
	{ports: []int{5900}, mark: markUnauthenticatedVNC},
	{ports: []int{5555}, mark: markUnauthenticatedADB},
}

func markUnchanged(_ context.Context, opened []probe.Open) []probe.Open {
	return opened
}

func runMarkers(ctx context.Context, opened []probe.Open, markers []portMarker) []probe.Open {
	owners := make(map[int]int)
	for markerIndex, marker := range markers {
		for _, port := range marker.ports {
			if _, exists := owners[port]; exists {
				panic("multiple port markers claim the same port")
			}
			owners[port] = markerIndex
		}
	}
	if len(opened) == 0 {
		return opened
	}

	unowned := make([]probe.Open, 0, len(opened))
	partitions := make([][]probe.Open, len(markers))
	for _, candidate := range opened {
		if markerIndex, owned := owners[candidate.Port]; owned {
			partitions[markerIndex] = append(partitions[markerIndex], candidate)
		} else {
			unowned = append(unowned, candidate)
		}
	}

	results := make([][]probe.Open, len(markers))
	var wait sync.WaitGroup
	for markerIndex, marker := range markers {
		if len(partitions[markerIndex]) == 0 {
			continue
		}
		wait.Add(1)
		go func(index int, current portMarker, candidates []probe.Open) {
			defer wait.Done()
			results[index] = current.mark(ctx, candidates)
		}(markerIndex, marker, partitions[markerIndex])
	}
	wait.Wait()

	marked := make([]probe.Open, 0, len(opened))
	marked = append(marked, unowned...)
	for _, result := range results {
		marked = append(marked, result...)
	}
	sort.SliceStable(marked, func(i, j int) bool {
		if marked[i].IP != marked[j].IP {
			return marked[i].IP.Less(marked[j].IP)
		}
		return originalPort(marked[i].Port) < originalPort(marked[j].Port)
	})
	return marked
}

func originalPort(port int) int {
	if port < 0 {
		return -port
	}
	return port
}
