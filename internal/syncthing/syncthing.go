package syncthing

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
	"unicode"
)

const maxResponseBytes = 1 << 20

const metadataPrefix = "var metadata = "

func ParseMeta(body []byte, headerID string) bool {
	body = bytes.TrimRightFunc(body, unicode.IsSpace)
	if !bytes.HasPrefix(body, []byte(metadataPrefix)) || len(body) == 0 || body[len(body)-1] != ';' {
		return false
	}
	var metadata struct {
		DeviceID *string `json:"deviceID"`
	}
	if err := json.Unmarshal(body[len(metadataPrefix):len(body)-1], &metadata); err != nil {
		return false
	}
	return metadata.DeviceID != nil && *metadata.DeviceID != "" && *metadata.DeviceID == headerID
}

// Unauthenticated reports a Syncthing GUI whose GET /meta.js returns 200 naming this device (the body's
// deviceID must match the X-Syncthing-ID header). /meta.js is outside the no-auth paths, so it returns
// 401/403 once a GUI user and password are set.
func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/meta.js", nil)
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
	return ParseMeta(body, response.Header.Get("X-Syncthing-ID"))
}
