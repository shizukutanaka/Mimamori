package nodered

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseSettings(body []byte) bool {
	var settings map[string]json.RawMessage
	if err := json.Unmarshal(body, &settings); err != nil || settings == nil {
		return false
	}
	for _, key := range []string{"httpNodeRoot", "version"} {
		raw, ok := settings[key]
		if !ok {
			return false
		}
		var value string
		if err := json.Unmarshal(raw, &value); err != nil || value == "" {
			return false
		}
	}
	userRaw, hasUser := settings["user"]
	if !hasUser {
		return true
	}
	var user map[string]json.RawMessage
	if err := json.Unmarshal(userRaw, &user); err != nil || user == nil {
		return false
	}
	permissions, hasPermissions := user["permissions"]
	return hasPermissions && hasFlowWritePermission(permissions)
}

func hasFlowWritePermission(permissions json.RawMessage) bool {
	hasWriteScope := func(scope string) bool {
		switch scope {
		case "*", "flows.write", "write", "*.write":
			return true
		default:
			return false
		}
	}

	var scope string
	if err := json.Unmarshal(permissions, &scope); err == nil {
		return hasWriteScope(scope)
	}

	var scopes []json.RawMessage
	if err := json.Unmarshal(permissions, &scopes); err != nil || scopes == nil {
		return false
	}
	for _, rawScope := range scopes {
		if err := json.Unmarshal(rawScope, &scope); err == nil && hasWriteScope(scope) {
			return true
		}
	}
	return false
}

// Unauthenticated reports a Node-RED editor whose GET /settings returns 200 with httpNodeRoot and version;
// with adminAuth configured the same request gets 401, and ParseSettings drops a read-only anonymous user.
func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/settings", nil)
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
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes))
	if err != nil {
		return false
	}
	return response.StatusCode == http.StatusOK && ParseSettings(body)
}
