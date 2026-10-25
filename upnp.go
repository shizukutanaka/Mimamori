package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func newUPnPHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func fetchPortMappings(
	ctx context.Context,
	client *http.Client,
	gateway netip.Addr,
	responses []ssdp.Response,
	descriptions map[netip.Addr]ssdp.Description,
) ([]ssdp.PortMapping, bool, error) {
	routerIP, selected := routerUPnPIP(gateway, responses, descriptions)
	if gateway.IsValid() {
		routerIP = gateway
		description, ok := descriptions[routerIP]
		if !ok {
			return nil, false, nil
		}
		isIGD := description.Device.IsIGD() ||
			(selected && routerIP == gateway)
		if !isIGD {
			return nil, false, nil
		}
	} else if !selected {
		return nil, false, nil
	}
	description, ok := descriptions[routerIP]
	if !ok {
		return nil, false, fmt.Errorf("router advertises UPnP but its port forwarding list cannot be queried")
	}
	type serviceEndpoint struct {
		service    ssdp.Service
		controlURL string
	}
	endpoints := make([]serviceEndpoint, 0)
	var errs []error
	for _, service := range description.Device.WANConnections() {
		controlURL, ok := ssdp.ControlURL(description.Location, description.URLBase, service.ControlURL, routerIP)
		if !ok {
			errs = append(errs, fmt.Errorf("%s: control URL is not on the router", service.ServiceType))
			continue
		}
		endpoints = append(endpoints, serviceEndpoint{service: service, controlURL: controlURL})
	}
	if len(endpoints) == 0 {
		return nil, false, fmt.Errorf("router advertises UPnP but its port forwarding list cannot be queried")
	}
	queryContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var mappings []ssdp.PortMapping
	for _, endpoint := range endpoints {
		serviceMappings, err := ssdp.PortMappings(queryContext, client, endpoint.controlURL, endpoint.service.ServiceType)
		mappings = append(mappings, serviceMappings...)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", endpoint.service.ServiceType, err))
		}
	}
	return mappings, true, errors.Join(errs...)
}

var descriptionFetchBudget = 5 * time.Second

func fetchDescriptions(ctx context.Context, client *http.Client, responses []ssdp.Response, gateway netip.Addr) (map[netip.Addr]ssdp.Description, []error) {
	ctx, cancel := context.WithTimeout(ctx, descriptionFetchBudget)
	defer cancel()

	selected := make(map[netip.Addr]ssdp.Response)
	for _, response := range responses {
		if !ssdp.AllowedLocation(response.Location, response.IP) {
			continue
		}
		current, exists := selected[response.IP]
		if !exists || descriptionPriority(response) > descriptionPriority(current) {
			selected[response.IP] = response
		}
	}

	type request struct {
		ip       netip.Addr
		location string
	}
	ordered := make([]request, 0, len(selected))
	for ip, response := range selected {
		ordered = append(ordered, request{ip: ip, location: response.Location})
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].ip.Less(ordered[j].ip)
	})
	for i := range ordered {
		if ordered[i].ip == gateway && i > 0 {
			gatewayRequest := ordered[i]
			copy(ordered[1:i+1], ordered[:i])
			ordered[0] = gatewayRequest
			break
		}
	}

	requests := make(chan request)
	results := make(chan descriptionResult)
	var workers sync.WaitGroup
	workerCount := 8
	if len(selected) < workerCount {
		workerCount = len(selected)
	}
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			for item := range requests {
				device, err := ssdp.FetchDescription(ctx, client, item.location)
				results <- descriptionResult{ip: item.ip, device: device, err: err}
			}
		}()
	}
	skipped := 0
	go func() {
	feed:
		for index, item := range ordered {
			if ctx.Err() != nil {
				skipped = len(ordered) - index
				break
			}
			select {
			case requests <- item:
			case <-ctx.Done():
				skipped = len(ordered) - index
				break feed
			}
		}
		close(requests)
		workers.Wait()
		close(results)
	}()

	devices := make(map[netip.Addr]ssdp.Description, len(selected))
	var errs []error
	for result := range results {
		if result.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", result.ip, result.err))
			continue
		}
		devices[result.ip] = result.device
	}
	if skipped > 0 {
		errs = append(errs, fmt.Errorf("description fetch time limit reached; %d devices not read", skipped))
	}
	return devices, errs
}

func descriptionPriority(response ssdp.Response) int {
	switch {
	case isUPnPIGDResponse(response.ST):
		return 3
	case ssdp.IsGateway(response):
		return 2
	case response.ST == "upnp:rootdevice":
		return 1
	default:
		return 0
	}
}

func isUPnPIGDResponse(st string) bool {
	st = strings.ToLower(st)
	return strings.HasPrefix(st, "urn:schemas-upnp-org:") &&
		(strings.Contains(st, ":device:internetgatewaydevice:") ||
			strings.Contains(st, ":device:wandevice:") ||
			strings.Contains(st, ":device:wanconnectiondevice:") ||
			strings.Contains(st, ":service:wanipconnection:") ||
			strings.Contains(st, ":service:wanpppconnection:") ||
			strings.Contains(st, ":service:wancommoninterfaceconfig:") ||
			strings.Contains(st, ":service:layer3forwarding:"))
}

func portMappingWarning(err error) string {
	return stripUnsafeRunes(strings.ReplaceAll(err.Error(), "\n", "; "))
}

func riskyPortForward(mapping ssdp.PortMapping) bool {
	if !strings.EqualFold(mapping.Protocol, "TCP") {
		return false
	}
	remoteHost := strings.TrimSpace(mapping.RemoteHost)
	if remoteHost != "" && remoteHost != "*" && remoteHost != "0.0.0.0" {
		return false
	}
	switch mapping.InternalPort {
	case 21, 22, 23, 135, 139, 445, 554, 8554, 1883, 1880, 1984, 2323, 2375, 3389, 5555, 5800, 5900, 1433, 3306, 5432, 6379, 9200, 27017, 34567, 37777:
		return true
	default:
		return false
	}
}

func routerUPnPIP(gateway netip.Addr, responses []ssdp.Response, descriptions map[netip.Addr]ssdp.Description) (netip.Addr, bool) {
	isIGDDevice := func(device ssdp.Device) bool {
		return device.IsIGD()
	}
	isIGDResponse := func(response ssdp.Response) bool {
		if ssdp.IsGateway(response) {
			return true
		}
		if description, ok := descriptions[response.IP]; ok && isIGDDevice(description.Device) {
			return true
		}
		return false
	}
	if gateway.IsValid() {
		if description, ok := descriptions[gateway]; ok && isIGDDevice(description.Device) {
			return gateway, true
		}
	}
	for _, response := range responses {
		if gateway.IsValid() && response.IP != gateway {
			continue
		}
		if isIGDResponse(response) {
			return response.IP, true
		}
	}
	return netip.Addr{}, false
}
