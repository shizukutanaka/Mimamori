package pihole

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

const maxResponseBytes = 1 << 20

func ParseAuth(body []byte) bool {
	var auth map[string]json.RawMessage
	if err := json.Unmarshal(body, &auth); err != nil || auth == nil {
		return false
	}
	sessionRaw, ok := auth["session"]
	if !ok {
		return false
	}
	var session map[string]json.RawMessage
	if err := json.Unmarshal(sessionRaw, &session); err != nil || session == nil {
		return false
	}
	validRaw, ok := session["valid"]
	if !ok || !bytes.Equal(bytes.TrimSpace(validRaw), []byte("true")) {
		return false
	}
	sidRaw, ok := session["sid"]
	if !ok || !bytes.Equal(bytes.TrimSpace(sidRaw), []byte("null")) {
		return false
	}
	validityRaw, ok := session["validity"]
	if !ok || !isNegativeOneJSONNumber(validityRaw) {
		return false
	}
	messageRaw, ok := session["message"]
	if !ok {
		return false
	}
	var message any
	if err := json.Unmarshal(messageRaw, &message); err != nil {
		return false
	}
	value, ok := message.(string)
	return ok && value == "no password set"
}

func isNegativeOneJSONNumber(raw json.RawMessage) bool {
	var decoded any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return false
	}
	number, ok := decoded.(json.Number)
	if !ok {
		return false
	}
	numberText := number.String()
	if !strings.HasPrefix(numberText, "-") {
		return false
	}
	numberText = numberText[1:]
	exponent := int64(0)
	if index := strings.IndexAny(numberText, "eE"); index >= 0 {
		parsed, err := strconv.ParseInt(numberText[index+1:], 10, 64)
		if err != nil {
			return false
		}
		exponent = parsed
		numberText = numberText[:index]
		maxExponent := int64(len(numberText) + 1)
		if exponent > maxExponent || exponent < -maxExponent {
			return false
		}
	}
	fractionalDigits := 0
	if index := strings.IndexByte(numberText, '.'); index >= 0 {
		fractionalDigits = len(numberText) - index - 1
		numberText = numberText[:index] + numberText[index+1:]
	}
	digits := strings.TrimLeft(numberText, "0")
	if digits == "" {
		return false
	}
	scale := exponent - int64(fractionalDigits)
	trailingZeros := len(digits) - len(strings.TrimRight(digits, "0"))
	if trailingZeros > 0 {
		digits = digits[:len(digits)-trailingZeros]
		scale += int64(trailingZeros)
	}
	return digits == "1" && scale == 0
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}
	address := netip.AddrPortFrom(ip, uint16(port)).String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/api/auth", nil)
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
	return ParseAuth(body)
}
