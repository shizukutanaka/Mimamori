package docker

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseVersion(body []byte) bool {
	var version map[string]json.RawMessage
	if err := json.Unmarshal(body, &version); err != nil || version == nil {
		return false
	}
	for _, key := range []string{"Version", "ApiVersion", "Os"} {
		raw, ok := version[key]
		if !ok {
			return false
		}
		var value string
		if err := json.Unmarshal(raw, &value); err != nil || value == "" {
			return false
		}
	}
	return true
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	return probeVersion(ctx, ip, port, timeout, "http", &http.Transport{Proxy: nil, DisableKeepAlives: true})
}

func UnauthenticatedTLS(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	return probeVersion(ctx, ip, port, timeout, "https", &http.Transport{
		Proxy:             nil,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12},
	})
}

func probeVersion(ctx context.Context, ip netip.Addr, port int, timeout time.Duration, scheme string, transport *http.Transport) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, scheme+"://"+address+"/version", nil)
	if err != nil {
		return false
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes))
	response.Body.Close()
	if err != nil {
		return false
	}
	if response.StatusCode != http.StatusOK || !ParseVersion(body) {
		return false
	}

	request, err = http.NewRequestWithContext(ctx, http.MethodOptions, scheme+"://"+address+"/containers/create", nil)
	if err != nil {
		return false
	}
	response, err = client.Do(request)
	if err != nil {
		return false
	}
	_, readErr := io.Copy(io.Discard, response.Body)
	closeErr := response.Body.Close()
	if readErr != nil || closeErr != nil {
		return false
	}
	return response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusForbidden
}
