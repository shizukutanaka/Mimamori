package go2rtc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 64 * 1024

func ParseAPIInfo(status int, body []byte) bool {
	if status != http.StatusOK {
		return false
	}
	var info map[string]json.RawMessage
	if err := json.Unmarshal(body, &info); err != nil {
		return false
	}
	version, ok := info["version"]
	if !ok {
		return false
	}
	var versionString string
	if err := json.Unmarshal(version, &versionString); err != nil || versionString == "" {
		return false
	}
	if _, ok := info["host"]; !ok {
		return false
	}
	var rtsp struct {
		Listen *string `json:"listen"`
	}
	raw, ok := info["rtsp"]
	return ok && json.Unmarshal(raw, &rtsp) == nil && rtsp.Listen != nil
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/api", nil)
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
	return ParseAPIInfo(response.StatusCode, body)
}
