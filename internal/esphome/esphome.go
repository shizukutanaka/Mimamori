package esphome

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseDevices(body []byte) bool {
	var devices map[string]json.RawMessage
	if err := json.Unmarshal(body, &devices); err != nil || devices == nil {
		return false
	}
	for _, key := range []string{"configured", "importable"} {
		raw, ok := devices[key]
		if !ok {
			return false
		}
		var entries []json.RawMessage
		if err := json.Unmarshal(raw, &entries); err != nil || entries == nil {
			return false
		}
	}
	return true
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/devices", nil)
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
	return response.StatusCode == http.StatusOK && ParseDevices(body)
}
