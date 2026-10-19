package frigate

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseStats(body []byte) bool {
	var stats map[string]json.RawMessage
	if err := json.Unmarshal(body, &stats); err != nil || stats == nil {
		return false
	}

	serviceJSON, ok := stats["service"]
	if !ok {
		return false
	}
	var service map[string]json.RawMessage
	if err := json.Unmarshal(serviceJSON, &service); err != nil || service == nil {
		return false
	}

	versionJSON, ok := service["version"]
	if !ok {
		return false
	}
	var version string
	if err := json.Unmarshal(versionJSON, &version); err != nil || version == "" {
		return false
	}
	uptimeJSON, ok := service["uptime"]
	if !ok {
		return false
	}
	var uptime float64
	if err := json.Unmarshal(uptimeJSON, &uptime); err != nil {
		return false
	}
	storageJSON, ok := service["storage"]
	if !ok {
		return false
	}
	var storage map[string]json.RawMessage
	if err := json.Unmarshal(storageJSON, &storage); err != nil || storage == nil {
		return false
	}
	if _, ok := service["latest_version"]; !ok {
		return false
	}

	detectorsJSON, ok := stats["detectors"]
	if !ok {
		return false
	}
	var detectors map[string]json.RawMessage
	if err := json.Unmarshal(detectorsJSON, &detectors); err != nil || detectors == nil {
		return false
	}

	return true
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/api/stats", nil)
	if err != nil {
		return false
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{Proxy: nil, DisableKeepAlives: true},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes))
	if err != nil {
		return false
	}
	return response.StatusCode == http.StatusOK && ParseStats(body)
}
