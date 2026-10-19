package main

import (
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

var (
	portableSDKVersionPattern = regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9_])Portable SDK for UPnP devices/([0-9]+)\.([0-9]+)\.([0-9]+)`)
	intelSDKPattern           = regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9_])Intel SDK for UPnP devices/`)
	miniupnpdVersionPattern   = regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9_])miniupnpd/([0-9]+)\.([0-9]+)`)
)

func vulnerableUPnPStackFindings(
	responses []ssdp.Response,
	gateway netip.Addr,
	deviceNames map[netip.Addr]string,
	lang string,
) []report.Finding {
	findings := make([]report.Finding, 0)
	seen := make(map[netip.Addr]struct{})
	for _, response := range responses {
		if !response.IP.IsValid() {
			continue
		}
		if _, exists := seen[response.IP]; exists {
			continue
		}
		component, vulnerable := vulnerableUPnPComponent(response.Server)
		if !vulnerable {
			continue
		}

		device := deviceNames[response.IP]
		if device == "" {
			device = response.IP.String()
		}
		level, points := report.Caution, 8
		if gateway.IsValid() && response.IP == gateway {
			level, points = report.Danger, 15
		}
		findings = append(findings, report.Finding{
			ID: "upnp_stack_vulnerable", Level: level, Points: points,
			IP: response.IP.String(), Device: device,
			Message: i18n.T(lang, "upnp_stack_vulnerable", device, component),
			Fix:     i18n.T(lang, "upnp_stack_vulnerable_fix"),
		})
		seen[response.IP] = struct{}{}
	}
	return findings
}

func vulnerableUPnPComponent(server string) (string, bool) {
	component := ""
	firstMatch := len(server) + 1
	for _, match := range portableSDKVersionPattern.FindAllStringSubmatchIndex(server, -1) {
		if !versionTokenBoundary(server, match[1]) {
			continue
		}
		major, majorErr := strconv.Atoi(server[match[2]:match[3]])
		minor, minorErr := strconv.Atoi(server[match[4]:match[5]])
		patch, patchErr := strconv.Atoi(server[match[6]:match[7]])
		if majorErr != nil || minorErr != nil || patchErr != nil {
			continue
		}
		if major > 1 || major == 1 && (minor > 6 || minor == 6 && patch >= 21) {
			continue
		}
		if match[0] < firstMatch {
			firstMatch = match[0]
			component = strings.Join([]string{
				"libupnp",
				server[match[2]:match[3]] + "." + server[match[4]:match[5]] + "." + server[match[6]:match[7]],
			}, " ")
		}
	}

	if match := intelSDKPattern.FindStringIndex(server); match != nil && match[0] < firstMatch {
		firstMatch = match[0]
		component = "Intel SDK for UPnP devices (all versions)"
	}

	for _, match := range miniupnpdVersionPattern.FindAllStringSubmatchIndex(server, -1) {
		if !versionTokenBoundary(server, match[1]) {
			continue
		}
		major, majorErr := strconv.Atoi(server[match[2]:match[3]])
		minor, minorErr := strconv.Atoi(server[match[4]:match[5]])
		if majorErr != nil || minorErr != nil {
			continue
		}
		if major > 1 || major == 1 && minor >= 4 {
			continue
		}
		if match[0] < firstMatch {
			firstMatch = match[0]
			component = strings.Join([]string{
				"MiniUPnPd",
				server[match[2]:match[3]] + "." + server[match[4]:match[5]],
			}, " ")
		}
	}
	return component, component != ""
}

func versionTokenBoundary(server string, end int) bool {
	if end >= len(server) {
		return true
	}
	next := server[end]
	return !(next >= '0' && next <= '9' ||
		next >= 'a' && next <= 'z' ||
		next >= 'A' && next <= 'Z' ||
		next == '.' || next == '_')
}
