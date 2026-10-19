package tasmota

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParsePage(body []byte) bool {
	return bytes.Contains(body, []byte("Theo Arends")) &&
		bytes.Contains(body, []byte(">Tasmota ")) &&
		bytes.Contains(body, []byte("action='up'"))
}

func requestURL(ip netip.Addr, port int) string {
	if port == 80 {
		return "http://" + ip.String() + "/"
	}
	return "http://" + netip.AddrPortFrom(ip, uint16(port)).String() + "/"
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
	return ParsePage(body)
}
