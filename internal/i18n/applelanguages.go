package i18n

import "strings"

func parseAppleLanguages(output string) string {
	output = strings.TrimSpace(output)
	if len(output) < 2 || output[0] != '(' || output[len(output)-1] != ')' {
		return ""
	}

	for _, entry := range strings.Split(output[1:len(output)-1], ",") {
		entry = strings.Trim(entry, " \t\r\n,\"")
		if entry == "" {
			continue
		}
		if strings.ContainsAny(entry, "\"(), \t\r\n") {
			return ""
		}
		return entry
	}
	return ""
}
