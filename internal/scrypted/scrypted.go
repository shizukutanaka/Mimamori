package scrypted

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseLogin(body []byte) bool {
	var login map[string]json.RawMessage
	if err := json.Unmarshal(body, &login); err != nil || login == nil {
		return false
	}
	hasLoginRaw, ok := login["hasLogin"]
	if !ok {
		return false
	}
	var hasLogin any
	if err := json.Unmarshal(hasLoginRaw, &hasLogin); err != nil {
		return false
	}
	hasAccount, ok := hasLogin.(bool)
	if !ok || hasAccount {
		return false
	}
	addressesRaw, ok := login["addresses"]
	if !ok {
		return false
	}
	var addresses any
	if err := json.Unmarshal(addressesRaw, &addresses); err != nil {
		return false
	}
	if _, ok := addresses.([]any); !ok {
		return false
	}
	hostnameRaw, ok := login["hostname"]
	if !ok {
		return false
	}
	var hostname any
	if err := json.Unmarshal(hostnameRaw, &hostname); err != nil {
		return false
	}
	_, ok = hostname.(string)
	return ok
}

func Unclaimed(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/login", nil)
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
	return ParseLogin(body)
}
