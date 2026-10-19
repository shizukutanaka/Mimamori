package shelly

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseInfo(body []byte) bool {
	var info map[string]json.RawMessage
	if err := json.Unmarshal(body, &info); err != nil || info == nil {
		return false
	}
	if generation, ok := info["gen"]; ok {
		var gen float64
		var model string
		authEnabled, hasAuth := info["auth_en"]
		if json.Unmarshal(generation, &gen) != nil || gen < 2 ||
			json.Unmarshal(info["model"], &model) != nil || model == "" ||
			!hasAuth || !isFalse(authEnabled) {
			return false
		}
		return true
	}
	var deviceType, mac string
	auth, hasAuth := info["auth"]
	if json.Unmarshal(info["type"], &deviceType) != nil || deviceType == "" ||
		json.Unmarshal(info["mac"], &mac) != nil || mac == "" ||
		!hasAuth || !isFalse(auth) {
		return false
	}
	return true
}

func isFalse(value json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(value), []byte("false"))
}

func requestURL(ip netip.Addr, port int) string {
	address := ip.String()
	if port != 80 {
		address = netip.AddrPortFrom(ip, uint16(port)).String()
	}
	return "http://" + address + "/shelly"
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL(ip, port), nil)
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
	if response.StatusCode != http.StatusOK {
		return false
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil || len(body) > maxResponseBytes {
		return false
	}
	return ParseInfo(body)
}
