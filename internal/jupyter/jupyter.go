package jupyter

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

func ParseStatus(status int, body []byte) bool {
	if status != http.StatusOK {
		return false
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(body, &object); err != nil || object == nil {
		return false
	}
	var started, lastActivity string
	if json.Unmarshal(object["started"], &started) != nil || started == "" ||
		json.Unmarshal(object["last_activity"], &lastActivity) != nil || lastActivity == "" {
		return false
	}
	if !isJSONNumber(object["kernels"]) || !isJSONNumber(object["connections"]) {
		return false
	}
	return true
}

func isJSONNumber(raw json.RawMessage) bool {
	var value any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if decoder.Decode(&value) != nil {
		return false
	}
	_, ok := value.(json.Number)
	return ok
}

// Unauthenticated reports a Jupyter server whose GET /api/status returns 200. The handler requires
// authentication, so 200 without a token means both the token and the password are disabled.
func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/api/status", nil)
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
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil || len(body) > maxResponseBytes {
		return false
	}
	return ParseStatus(response.StatusCode, body)
}
