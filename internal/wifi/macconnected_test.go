package wifi

import (
	"errors"
	"testing"
)

func TestMacSecurityResult(t *testing.T) {
	tests := []struct {
		name         string
		jsonText     string
		wantSecurity Security
		wantErr      error
	}{
		{
			name:         "connected WPA2",
			jsonText:     `{"SPAirPortDataType":[{"spairport_airport_interfaces":[{"spairport_status_information":"spairport_status_connected","spairport_current_network_information":{"_name":"Home","spairport_security_mode":"spairport_security_mode_wpa2_personal"}}]}]}`,
			wantSecurity: Modern,
		},
		{
			name:         "connected but security is missing",
			jsonText:     `{"SPAirPortDataType":[{"spairport_airport_interfaces":[{"spairport_status_information":"spairport_status_connected","spairport_current_network_information":{"_name":"Home"}}]}]}`,
			wantSecurity: Unknown,
			wantErr:      ErrWiFiUnreadable,
		},
		{
			name:         "network without status",
			jsonText:     `{"SPAirPortDataType":[{"spairport_airport_interfaces":[{"spairport_current_network_information":{"_name":"Home"}}]}]}`,
			wantSecurity: Unknown,
			wantErr:      ErrWiFiUnreadable,
		},
		{
			name:         "connected open network",
			jsonText:     `{"SPAirPortDataType":[{"spairport_airport_interfaces":[{"spairport_status_information":"spairport_status_connected","spairport_current_network_information":{"_name":"Guest","spairport_security_mode":"spairport_security_mode_none"}}]}]}`,
			wantSecurity: Open,
		},
		{
			name:         "inactive without network",
			jsonText:     `{"SPAirPortDataType":[{"spairport_airport_interfaces":[{"spairport_status_information":"spairport_status_inactive"}]}]}`,
			wantSecurity: Unknown,
		},
		{
			name:         "off without network",
			jsonText:     `{"SPAirPortDataType":[{"spairport_airport_interfaces":[{"spairport_status_information":"spairport_status_off"}]}]}`,
			wantSecurity: Unknown,
		},
		{
			name:         "inactive with empty network",
			jsonText:     `{"SPAirPortDataType":[{"spairport_airport_interfaces":[{"spairport_status_information":"spairport_status_inactive","spairport_current_network_information":{}}]}]}`,
			wantSecurity: Unknown,
		},
		{
			name:         "software information only",
			jsonText:     `{"SPAirPortDataType":[{"spairport_software_information":{"spairport_corewlan_version":"16.0 (1657)"}}]}`,
			wantSecurity: Unknown,
		},
		{
			name:         "invalid JSON",
			jsonText:     `{`,
			wantSecurity: Unknown,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotSecurity, gotErr := macSecurityResult([]byte(test.jsonText))
			if gotSecurity != test.wantSecurity {
				t.Errorf("macSecurityResult() security = %v, want %v", gotSecurity, test.wantSecurity)
			}
			if !errors.Is(gotErr, test.wantErr) {
				t.Errorf("macSecurityResult() error = %v, want %v", gotErr, test.wantErr)
			}
		})
	}
}
