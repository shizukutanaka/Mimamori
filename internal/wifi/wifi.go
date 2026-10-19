package wifi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
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

var ErrLocationPermission = errors.New("Windows location permission is required to read Wi-Fi security")

func Current(ctx context.Context) (Security, error) {
	commandContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := commandContext.Err(); err != nil {
		return Unknown, err
	}

	if runtime.GOOS == "windows" {
		if security, handled, err := currentWLAN(); handled {
			if err := commandContext.Err(); err != nil {
				return Unknown, err
			}
			return security, err
		}
	}

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

	var errorOutput bytes.Buffer
	command.Stderr = &errorOutput
	output, err := command.Output()
	if commandContext.Err() != nil {
		return Unknown, commandContext.Err()
	}
	if runtime.GOOS == "windows" && (NetshNeedsLocation(string(output)) || NetshNeedsLocation(errorOutput.String())) {
		return Unknown, ErrLocationPermission
	}
	if err != nil {
		if runtime.GOOS == "linux" {
			return wpaCliSecurity(commandContext)
		}
		return Unknown, nil
	}

	switch runtime.GOOS {
	case "darwin":
		return macSecurityResult(output)
	case "linux":
		security := ParseNmcli(string(output))
		if security == Unknown {
			return wpaCliSecurity(commandContext)
		}
		return security, nil
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
			case strings.Contains(security, "WPA1"), strings.Contains(security, "TKIP"):
				return WPA1orTKIP
			case strings.Contains(security, "WPA2"), strings.Contains(security, "WPA3"):
				return Modern
			default:
				return Unknown
			}
		}
	}
	return Unknown
}

func wpaCliSecurity(ctx context.Context) (Security, error) {
	interfaces := supplicantInterfaces([]string{"/run/wpa_supplicant", "/var/run/wpa_supplicant"})
	if len(interfaces) == 0 {
		security, err := wpaCliStatus(ctx, "status")
		return unreadableIfConnected("/sys/class/net", security, err)
	}

	results := make([]Security, 0, len(interfaces))
	for _, iface := range interfaces {
		security, err := wpaCliStatus(ctx, "-p", iface.dir, "-i", iface.name, "status")
		if err != nil {
			return Unknown, err
		}
		results = append(results, security)
	}
	return unreadableIfConnected("/sys/class/net", worstSecurity(results), nil)
}

type supplicantInterface struct {
	dir  string
	name string
}

func supplicantInterfaces(dirs []string) []supplicantInterface {
	seenDirs := make(map[string]struct{}, len(dirs))
	var interfaces []supplicantInterface
	for _, dir := range dirs {
		resolvedDir, err := filepath.EvalSymlinks(dir)
		if err != nil {
			continue
		}
		if _, ok := seenDirs[resolvedDir]; ok {
			continue
		}
		seenDirs[resolvedDir] = struct{}{}

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || strings.HasPrefix(name, "p2p-dev-") {
				continue
			}
			interfaces = append(interfaces, supplicantInterface{dir: dir, name: name})
		}
	}
	sort.Slice(interfaces, func(i, j int) bool {
		if interfaces[i].dir != interfaces[j].dir {
			return interfaces[i].dir < interfaces[j].dir
		}
		return interfaces[i].name < interfaces[j].name
	})
	return interfaces
}

func worstSecurity(results []Security) Security {
	worst := Unknown
	for _, result := range results {
		if result == Unknown {
			continue
		}
		if worst == Unknown || result < worst {
			worst = result
		}
	}
	return worst
}

func wpaCliStatus(ctx context.Context, args ...string) (Security, error) {
	command := exec.CommandContext(ctx, "wpa_cli", args...)
	command.Env = withLCAllC(os.Environ())

	var errorOutput bytes.Buffer
	command.Stderr = &errorOutput
	output, err := command.Output()
	if ctx.Err() != nil {
		return Unknown, ctx.Err()
	}
	if err != nil {
		return Unknown, nil
	}
	return ParseWPACli(string(output)), nil
}

func ParseWPACli(text string) Security {
	var state, keyManagement, pairwiseCipher, groupCipher string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSuffix(line, "\r")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch key {
		case "wpa_state":
			state = value
		case "key_mgmt":
			keyManagement = value
		case "pairwise_cipher":
			pairwiseCipher = value
		case "group_cipher":
			groupCipher = value
		}
	}
	if state != "COMPLETED" {
		return Unknown
	}

	hasDynamicWEP, hasNone, hasWPA1 := false, false, false
	for _, method := range strings.Split(keyManagement, "+") {
		switch method {
		case "IEEE 802.1X (no WPA)":
			hasDynamicWEP = true
		case "NONE":
			hasNone = true
		case "WPA-PSK", "WPA-NONE":
			hasWPA1 = true
		default:
			if strings.HasPrefix(method, "WPA/") {
				hasWPA1 = true
			}
		}
	}
	if hasDynamicWEP {
		return WEP
	}
	if hasWPA1 {
		return WPA1orTKIP
	}
	if hasNone {
		if strings.HasPrefix(pairwiseCipher, "WEP") || strings.HasPrefix(groupCipher, "WEP") {
			return WEP
		}
		return Open
	}
	if pairwiseCipher == "TKIP" || groupCipher == "TKIP" {
		return WPA1orTKIP
	}
	if keyManagement != "" {
		return Modern
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
		var fields []string
		connected := false
		for _, line := range strings.Split(block, "\n") {
			label, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			label = strings.TrimSpace(label)
			value = strings.TrimSpace(value)
			isSSID := strings.Contains(strings.ToUpper(label), "SSID")
			if !isSSID && isConnectedValue(value) {
				connected = true
			}
			if isSSID {
				continue
			}
			fields = append(fields, value)
		}
		if !connected {
			continue
		}
		for i := 0; i+1 < len(fields); i++ {
			if !isNetshAuthentication(fields[i]) || !isNetshCipher(fields[i+1]) {
				continue
			}
			candidate := netshSecurity(fields[i], fields[i+1])
			if securityRank(candidate) > securityRank(best) {
				best = candidate
			}
		}
	}
	return best
}

func NetshNeedsLocation(text string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "ms-settings:privacy-location" {
			return true
		}
	}
	return false
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

func isNetshAuthentication(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "OPEN", "SHARED", "WPA-PERSONAL", "WPA-ENTERPRISE",
		"WPA2-PERSONAL", "WPA2-ENTERPRISE", "WPA3-PERSONAL",
		"WPA3-ENTERPRISE", "WPA3-SAE", "OWE":
		return true
	default:
		return false
	}
}

func isNetshCipher(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "NONE", "WEP", "TKIP", "CCMP", "GCMP", "CCMP-256", "GCMP-256":
		return true
	default:
		return false
	}
}

func netshSecurity(authentication, cipher string) Security {
	authentication = strings.ToUpper(strings.TrimSpace(authentication))
	cipher = strings.ToUpper(strings.TrimSpace(cipher))
	switch {
	case authentication == "SHARED" || cipher == "WEP":
		return WEP
	case authentication == "OPEN" && cipher == "NONE":
		return Open
	case cipher == "TKIP":
		return WPA1orTKIP
	case authentication == "WPA-PERSONAL" || authentication == "WPA-ENTERPRISE":
		return WPA1orTKIP
	case isModernNetshAuthentication(authentication) && isModernNetshCipher(cipher):
		return Modern
	default:
		return Unknown
	}
}

func isModernNetshAuthentication(authentication string) bool {
	switch authentication {
	case "WPA2-PERSONAL", "WPA2-ENTERPRISE", "WPA3-PERSONAL", "WPA3-ENTERPRISE", "WPA3-SAE", "OWE":
		return true
	default:
		return false
	}
}

func isModernNetshCipher(cipher string) bool {
	switch cipher {
	case "CCMP", "GCMP", "CCMP-256", "GCMP-256":
		return true
	default:
		return false
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
