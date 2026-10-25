package codeserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseHealth(status int, body []byte) bool {
	if status != http.StatusOK {
		return false
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(body, &object); err != nil || object == nil {
		return false
	}
	var healthStatus string
	if json.Unmarshal(object["status"], &healthStatus) != nil || (healthStatus != "alive" && healthStatus != "expired") {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(object["lastHeartbeat"]))
	decoder.UseNumber()
	var heartbeat any
	if err := decoder.Decode(&heartbeat); err != nil {
		return false
	}
	_, ok := heartbeat.(json.Number)
	return ok
}

func LoginRedirectsHome(status int, location string) bool {
	if status != http.StatusFound || location == "" {
		return false
	}
	base, err := url.Parse("http://h/login")
	if err != nil {
		return false
	}
	reference, err := url.Parse(location)
	if err != nil {
		return false
	}
	resolved := base.ResolveReference(reference)
	return (resolved.Host == "" || resolved.Host == base.Host) && resolved.Path == "/"
}

// Unauthenticated reports a code-server whose /healthz answers and whose /login is absent or sends the
// browser back to /. With the default `auth: password` the login page is served, so this means `auth: none`.
func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
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
	healthStatus, healthBody, _, ok := get(client, ctx, baseURL+"/healthz")
	if !ok || !ParseHealth(healthStatus, healthBody) {
		return false
	}
	loginStatus, _, location, ok := get(client, ctx, baseURL+"/login")
	return ok && LoginRedirectsHome(loginStatus, location)
}

func get(client *http.Client, ctx context.Context, requestURL string) (int, []byte, string, bool) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return 0, nil, "", false
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, "", false
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil || len(body) > maxResponseBytes {
		return 0, nil, "", false
	}
	return response.StatusCode, body, response.Header.Get("Location"), true
}
