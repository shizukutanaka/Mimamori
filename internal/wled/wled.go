package wled

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
	var version string
	if err := json.Unmarshal(info["ver"], &version); err != nil || version == "" {
		return false
	}
	if !isJSONNumber(info["vid"]) {
		return false
	}
	var leds map[string]json.RawMessage
	if err := json.Unmarshal(info["leds"], &leds); err != nil || leds == nil {
		return false
	}
	return true
}

func isJSONNumber(value json.RawMessage) bool {
	var decoded any
	decoder := json.NewDecoder(bytes.NewReader(value))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return false
	}
	_, ok := decoded.(json.Number)
	return ok
}

func requestURL(ip netip.Addr, port int, path string) string {
	address := ip.String()
	if port != 80 {
		address = netip.AddrPortFrom(ip, uint16(port)).String()
	}
	return "http://" + address + path
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{Proxy: nil, DisableKeepAlives: true},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	infoRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL(ip, port, "/json/info"), nil)
	if err != nil {
		return false
	}
	infoResponse, err := client.Do(infoRequest)
	if err != nil {
		return false
	}
	infoBody, ok := readResponseBody(infoResponse)
	if !ok || infoResponse.StatusCode != http.StatusOK || !ParseInfo(infoBody) {
		return false
	}

	updateRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL(ip, port, "/update"), nil)
	if err != nil {
		return false
	}
	updateResponse, err := client.Do(updateRequest)
	if err != nil {
		return false
	}
	_, ok = readResponseBody(updateResponse)
	return ok && updateResponse.StatusCode == http.StatusOK
}

func readResponseBody(response *http.Response) ([]byte, bool) {
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil || len(body) > maxResponseBytes {
		return nil, false
	}
	return body, true
}
