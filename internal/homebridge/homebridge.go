package homebridge

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseSettings(body []byte) bool {
	var settings map[string]json.RawMessage
	if err := json.Unmarshal(body, &settings); err != nil || settings == nil {
		return false
	}
	var env map[string]json.RawMessage
	if err := json.Unmarshal(settings["env"], &env); err != nil || env == nil {
		return false
	}
	var packageName string
	if err := json.Unmarshal(env["packageName"], &packageName); err != nil || packageName != "homebridge-config-ui-x" {
		return false
	}
	formAuthRaw, ok := settings["formAuth"]
	if !ok {
		return false
	}
	var formAuth any
	if err := json.Unmarshal(formAuthRaw, &formAuth); err != nil {
		return false
	}
	enabled, ok := formAuth.(bool)
	return ok && !enabled
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/api/auth/settings", nil)
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
	return ParseSettings(body)
}
