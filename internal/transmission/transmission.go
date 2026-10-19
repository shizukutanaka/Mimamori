package transmission

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseConflict(status int, headerID string, body []byte) bool {
	headerID = strings.TrimSpace(headerID)
	if status != http.StatusConflict || headerID == "" {
		return false
	}
	return bytes.Contains(body, []byte("X-Transmission-Session-Id: "+headerID))
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/transmission/rpc", nil)
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
	return ParseConflict(response.StatusCode, response.Header.Get("X-Transmission-Session-Id"), body)
}
