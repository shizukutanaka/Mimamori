package wifi

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Security int

const (
	Unknown Security = iota
	Open
	WEP
	WPA1orTKIP
	Modern
)

func Current(ctx context.Context) (Security, error) {
	commandContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.CommandContext(commandContext, "system_profiler", "SPAirPortDataType", "-json")
	case "linux":
		command = exec.CommandContext(commandContext, "nmcli", "-t", "-f", "ACTIVE,SECURITY", "dev", "wifi")
		command.Env = withLCAllC(os.Environ())
	case "windows":
		command = exec.CommandContext(commandContext, "netsh", "wlan", "show", "interfaces")
	default:
		return Unknown, nil
	}

	output, err := command.Output()
	if err != nil {
		if commandContext.Err() != nil {
			return Unknown, commandContext.Err()
		}
		return Unknown, nil
	}

	switch runtime.GOOS {
	case "darwin":
		return ParseMacOS(output), nil
	case "linux":
		return ParseNmcli(string(output)), nil
	case "windows":
		return ParseNetsh(string(output)), nil
	default:
		return Unknown, nil
	}
}

func withLCAllC(environment []string) []string {
	result := make([]string, 0, len(environment)+1)
	for _, entry := range environment {
		if !strings.HasPrefix(entry, "LC_ALL=") {
			result = append(result, entry)
		}
	}
	return append(result, "LC_ALL=C")
}

func ParseMacOS(jsonText []byte) Security {
	var data struct {
		Interfaces []struct {
			AirPortInterfaces []struct {
				CurrentNetwork struct {
					SecurityMode string `json:"spairport_security_mode"`
				} `json:"spairport_current_network_information"`
			} `json:"spairport_airport_interfaces"`
		} `json:"SPAirPortDataType"`
	}
	if err := json.Unmarshal(jsonText, &data); err != nil {
		return Unknown
	}
	best := Unknown
	for _, section := range data.Interfaces {
		for _, iface := range section.AirPortInterfaces {
			candidate := macSecurity(iface.CurrentNetwork.SecurityMode)
			if securityRank(candidate) > securityRank(best) {
				best = candidate
			}
		}
	}
	return best
}

func macSecurity(mode string) Security {
	mode = strings.ToLower(strings.TrimSpace(mode))
	switch {
	case strings.Contains(mode, "wpa2"), strings.Contains(mode, "wpa3"):
		return Modern
	case strings.Contains(mode, "wep"):
		return WEP
	case strings.Contains(mode, "none"):
		return Open
	case strings.Contains(mode, "wpa_personal"), strings.Contains(mode, "wpa_enterprise"), strings.Contains(mode, "wpa1"):
		return WPA1orTKIP
	default:
		return Unknown
	}
}

func ParseNmcli(text string) Security {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if strings.HasPrefix(line, "yes:") {
			security := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(line, "yes:")))
			switch {
			case security == "":
				return Open
			case strings.Contains(security, "WEP"):
				return WEP
			case strings.Contains(security, "WPA2"), strings.Contains(security, "WPA3"):
				return Modern
			case strings.Contains(security, "WPA1"), strings.Contains(security, "TKIP"):
				return WPA1orTKIP
			default:
				return Unknown
			}
		}
	}
	return Unknown
}

func ParseNetsh(text string) Security {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	var blocks []string
	var block strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			if block.Len() > 0 {
				blocks = append(blocks, block.String())
				block.Reset()
			}
			continue
		}
		block.WriteString(line)
		block.WriteByte('\n')
	}
	if block.Len() > 0 {
		blocks = append(blocks, block.String())
	}

	best := Unknown
	for _, block := range blocks {
		var values []string
		connected := false
		for _, line := range strings.Split(block, "\n") {
			_, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			value = strings.TrimSpace(value)
			values = append(values, value)
			if isConnectedValue(value) {
				connected = true
			}
		}
		if !connected {
			continue
		}
		candidate := netshSecurity(values)
		if securityRank(candidate) > securityRank(best) {
			best = candidate
		}
	}
	return best
}

func isConnectedValue(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, prefix := range []string{
		"connected", "conectado", "conectada", "verbunden", "connecté", "connectée",
		"подключено", "подключен", "подключена", "연결됨", "terhubung", "已连接", "已連接", "接続",
	} {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func netshSecurity(values []string) Security {
	hasOpen, hasNone, hasTKIP, hasModern, hasModernCipher := false, false, false, false, false
	for _, value := range values {
		upper := strings.ToUpper(strings.TrimSpace(value))
		switch upper {
		case "WEP":
			return WEP
		case "TKIP":
			hasTKIP = true
		case "OPEN":
			hasOpen = true
		case "NONE":
			hasNone = true
		}
		if strings.Contains(upper, "WPA2") || strings.Contains(upper, "WPA3") {
			hasModern = true
		}
		if strings.Contains(upper, "CCMP") || strings.Contains(upper, "GCMP") {
			hasModernCipher = true
		}
	}
	switch {
	case hasTKIP:
		return WPA1orTKIP
	case hasOpen && hasNone:
		return Open
	case hasModern && hasModernCipher:
		return Modern
	default:
		return Unknown
	}
}

func securityRank(security Security) int {
	switch security {
	case Open, WEP:
		return 3
	case WPA1orTKIP:
		return 2
	case Modern:
		return 1
	default:
		return 0
	}
}
