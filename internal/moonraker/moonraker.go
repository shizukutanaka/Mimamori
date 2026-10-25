package moonraker

import (
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
	rawResult, ok := info["result"]
	if !ok {
		return false
	}
	var result map[string]json.RawMessage
	if err := json.Unmarshal(rawResult, &result); err != nil || result == nil {
		return false
	}
	return isNonEmptyString(result["moonraker_version"]) && isNonEmptyString(result["api_version_string"])
}

func isNonEmptyString(raw json.RawMessage) bool {
	var value string
	return json.Unmarshal(raw, &value) == nil && value != ""
}

// Unauthenticated reports a Moonraker whose GET /server/info returns 200 with moonraker_version. The
// endpoint requires authorization, so 200 means this PC is in trusted_clients (MainsailOS trusts every
// private range by default); an untrusted client gets 401.
func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/server/info", nil)
	if err != nil {
		return false
	}
	request.Header.Set("Accept", "application/json")
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
