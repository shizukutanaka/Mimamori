package wifi

import (
	"bytes"
	"encoding/json"
)

// macConnected reports whether system_profiler SPAirPortDataType -json shows a Wi-Fi interface associated to a network.
func macConnected(jsonText []byte) bool {
	var data struct {
		Sections []struct {
			Interfaces []struct {
				Status         string          `json:"spairport_status_information"`
				CurrentNetwork json.RawMessage `json:"spairport_current_network_information"`
			} `json:"spairport_airport_interfaces"`
		} `json:"SPAirPortDataType"`
	}
	if err := json.Unmarshal(jsonText, &data); err != nil {
		return false
	}
	for _, section := range data.Sections {
		for _, iface := range section.Interfaces {
			if iface.Status == "spairport_status_connected" {
				return true
			}
			switch string(bytes.TrimSpace(iface.CurrentNetwork)) {
			case "", "null", "{}", "[]":
			default:
				return true
			}
		}
	}
	return false
}

func macSecurityResult(jsonText []byte) (Security, error) {
	security := ParseMacOS(jsonText)
	if security == Unknown && macConnected(jsonText) {
		return Unknown, ErrWiFiUnreadable
	}
	return security, nil
}
