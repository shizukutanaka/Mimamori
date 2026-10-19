package homeassistant

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseOnboarding(body []byte) bool {
	var steps []json.RawMessage
	if err := json.Unmarshal(body, &steps); err != nil || len(steps) == 0 {
		return false
	}
	userStepOpen := false
	for _, rawStep := range steps {
		var step map[string]json.RawMessage
		if err := json.Unmarshal(rawStep, &step); err != nil || step == nil {
			return false
		}
		stepRaw, ok := step["step"]
		if !ok {
			return false
		}
		var name string
		if err := json.Unmarshal(stepRaw, &name); err != nil || name == "" {
			return false
		}
		doneRaw, ok := step["done"]
		if !ok {
			return false
		}
		var value any
		if err := json.Unmarshal(doneRaw, &value); err != nil {
			return false
		}
		done, ok := value.(bool)
		if !ok {
			return false
		}
		if name == "user" && !done {
			userStepOpen = true
		}
	}
	return userStepOpen
}

func UserStepOpen(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/api/onboarding", nil)
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
	return response.StatusCode == http.StatusOK && ParseOnboarding(body)
}
