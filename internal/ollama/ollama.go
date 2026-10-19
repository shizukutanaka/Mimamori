package ollama

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

func ParseRoot(status int, body []byte) bool {
	return status == http.StatusOK && bytes.Equal(bytes.TrimSpace(body), []byte("Ollama is running"))
}

func ParseVersion(status int, body []byte) bool {
	if status != http.StatusOK {
		return false
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(body, &object); err != nil || object == nil {
		return false
	}
	var version string
	if json.Unmarshal(object["version"], &version) != nil || version == "" {
		return false
	}
	return true
}

func Open(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	baseURL := "http://" + address
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{Proxy: nil, DisableKeepAlives: true},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	rootStatus, rootBody, ok := get(client, ctx, baseURL+"/")
	if !ok || !ParseRoot(rootStatus, rootBody) {
		return false
	}
	versionStatus, versionBody, ok := get(client, ctx, baseURL+"/api/version")
	return ok && ParseVersion(versionStatus, versionBody)
}

func get(client *http.Client, ctx context.Context, url string) (int, []byte, bool) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, false
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, false
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil || len(body) > maxResponseBytes {
		return 0, nil, false
	}
	return response.StatusCode, body, true
}
