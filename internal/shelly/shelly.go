package shelly

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseInfo(body []byte) bool {
	var info map[string]json.RawMessage
	if err := json.Unmarshal(body, &info); err != nil || info == nil {
		return false
	}
	if generation, ok := info["gen"]; ok {
		var gen float64
		var model string
		authEnabled, hasAuth := info["auth_en"]
		if json.Unmarshal(generation, &gen) != nil || gen < 2 ||
			json.Unmarshal(info["model"], &model) != nil || model == "" ||
			!hasAuth || !isFalse(authEnabled) {
			return false
		}
		return true
	}
	var deviceType, mac string
	auth, hasAuth := info["auth"]
	if json.Unmarshal(info["type"], &deviceType) != nil || deviceType == "" ||
		json.Unmarshal(info["mac"], &mac) != nil || mac == "" ||
		!hasAuth || !isFalse(auth) {
		return false
	}
	return true
}

func isFalse(value json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(value), []byte("false"))
}

func requestURL(ip netip.Addr, port int) string {
	address := ip.String()
	if port != 80 {
		address = netip.AddrPortFrom(ip, uint16(port)).String()
	}
	return "http://" + address + "/shelly"
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	return unauthenticated(ctx, ip, port, 443, timeout)
}

func unauthenticated(ctx context.Context, ip netip.Addr, port, httpsPort int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || httpsPort < 1 || httpsPort > 65535 || timeout <= 0 {
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
	if response.StatusCode == http.StatusOK {
		return parseResponse(response)
	}
	locationHeader := response.Header.Get("Location")
	response.Body.Close()
	if !isHTTPSRedirect(response.StatusCode) || locationHeader == "" {
		return false
	}

	location, err := url.Parse(locationHeader)
	if err != nil {
		return false
	}
	location = request.URL.ResolveReference(location)
	if !strings.EqualFold(location.Scheme, "https") {
		return false
	}
	redirectIP, err := netip.ParseAddr(location.Hostname())
	if err != nil || redirectIP != ip {
		return false
	}
	if !redirectPortMatches(location.Port(), httpsPort) {
		return false
	}

	httpsRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, httpsRequestURL(ip, httpsPort), nil)
	if err != nil {
		return false
	}
	httpsClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy:             nil,
			DisableKeepAlives: true,
			// Shelly uses its own internal PKI; only the public /shelly info is read.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	httpsResponse, err := httpsClient.Do(httpsRequest)
	if err != nil {
		return false
	}
	return parseResponse(httpsResponse)
}

func httpsRequestURL(ip netip.Addr, port int) string {
	address := ip.String()
	if port != 443 {
		address = netip.AddrPortFrom(ip, uint16(port)).String()
	}
	return "https://" + address + "/shelly"
}

func isHTTPSRedirect(status int) bool {
	switch status {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		return true
	default:
		return false
	}
}

func redirectPortMatches(locationPort string, httpsPort int) bool {
	if httpsPort == 443 {
		return locationPort == "" || locationPort == "443"
	}
	return locationPort == strconv.Itoa(httpsPort)
}

func parseResponse(response *http.Response) bool {
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return false
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil || len(body) > maxResponseBytes {
		return false
	}
	return ParseInfo(body)
}
