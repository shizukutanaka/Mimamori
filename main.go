package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/shizukutanaka/Mimamori/internal/baseline"
	"github.com/shizukutanaka/Mimamori/internal/console"
	"github.com/shizukutanaka/Mimamori/internal/i18n"
	"github.com/shizukutanaka/Mimamori/internal/lan"
	"github.com/shizukutanaka/Mimamori/internal/probe"
	"github.com/shizukutanaka/Mimamori/internal/report"
	"github.com/shizukutanaka/Mimamori/internal/ssdp"
	"github.com/shizukutanaka/Mimamori/internal/wifi"
)

var version = "dev"

type deviceOutput struct {
	IP   string `json:"ip"`
	MAC  string `json:"mac"`
	Name string `json:"name"`
	New  bool   `json:"new"`
}

type output struct {
	Version    string           `json:"version"`
	Score      int              `json:"score"`
	Status     string           `json:"status"`
	FirstRun   bool             `json:"first_run"`
	Incomplete bool             `json:"incomplete"`
	Gateway    string           `json:"gateway"`
	Devices    []deviceOutput   `json:"devices"`
	Findings   []report.Finding `json:"findings"`
}

type historyLoadState uint8

const (
	historyLoadReady historyLoadState = iota
	historyLoadQuarantined
	historyLoadUnavailable
)

type runResult struct {
	code  int
	lang  string
	pause bool
}

type discoveryResult struct {
	neighbors []lan.Neighbor
	err       error
}

type gatewayResult struct {
	gateway netip.Addr
	err     error
}

type ssdpResult struct {
	responses []ssdp.Response
	err       error
}

type wifiResult struct {
	security wifi.Security
	err      error
}

type descriptionResult struct {
	ip     netip.Addr
	device ssdp.Device
	err    error
}

func main() {
	result := run(os.Args[1:], os.Stdout, os.Stderr)
	if result.pause && console.ShouldPause() {
		fmt.Fprint(os.Stdout, i18n.T(result.lang, "press_enter"))
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	}
	os.Exit(result.code)
}

func run(args []string, stdout, stderr io.Writer) runResult {
	flags := flag.NewFlagSet("mimamori", flag.ContinueOnError)
	flags.SetOutput(stderr)
	jsonOutput := flags.Bool("json", false, "write JSON output")
	languageHelp := fmt.Sprintf("output language: auto or one of %s", strings.Join(i18n.Supported(), ", "))
	language := flags.String("lang", "auto", languageHelp)
	statePath := flags.String("state", "", "path to the local device history file")
	wait := flags.Duration("wait", 3*time.Second, "time to listen for local UPnP devices")
	showVersion := flags.Bool("version", false, "print version")
	lang := ""
	finish := func(code int) runResult {
		return runResult{code: code, lang: lang, pause: !*jsonOutput}
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return finish(0)
		}
		return finish(1)
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "unexpected argument: %s\n", flags.Arg(0))
		return finish(1)
	}
	if *showVersion {
		fmt.Fprintf(stdout, "Mimamori %s\n", version)
		return finish(0)
	}
	resolvedLang, err := resolveLanguage(*language, os.Getenv)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return finish(1)
	}
	lang = resolvedLang
	if *wait < 0 {
		fmt.Fprintln(stderr, "--wait cannot be negative")
		return finish(1)
	}
	if *statePath == "" {
		var err error
		*statePath, err = baseline.DefaultPath()
		if err != nil {
			fmt.Fprintf(stderr, "warning: cannot determine state path: %v\n", err)
			return finish(1)
		}
	}

	ctx := context.Background()
	arpCh := make(chan discoveryResult, 1)
	gatewayCh := make(chan gatewayResult, 1)
	ssdpCh := make(chan ssdpResult, 1)
	wifiCh := make(chan wifiResult, 1)
	go func() {
		neighbors, err := lan.ReadARP(ctx)
		arpCh <- discoveryResult{neighbors: neighbors, err: err}
	}()
	go func() {
		gateway, err := lan.DefaultGateway(ctx)
		gatewayCh <- gatewayResult{gateway: gateway, err: err}
	}()
	go func() {
		responses, err := ssdp.Discover(ctx, *wait)
		ssdpCh <- ssdpResult{responses: responses, err: err}
	}()
	go func() {
		security, err := wifi.Current(ctx)
		wifiCh <- wifiResult{security: security, err: err}
	}()

	arp := <-arpCh
	gateway := <-gatewayCh
	discovery := <-ssdpCh
	wifiCheck := <-wifiCh
	if arp.err != nil {
		fmt.Fprintf(stderr, "warning: cannot read the neighbor table: %v\n", arp.err)
	}
	var gatewayIP netip.Addr
	if gateway.err != nil {
		fmt.Fprintf(stderr, "warning: cannot find the default gateway: %v\n", gateway.err)
	} else if gateway.gateway.Is4() && gateway.gateway.IsPrivate() {
		gatewayIP = gateway.gateway
	}
	if discovery.err != nil {
		fmt.Fprintf(stderr, "warning: UPnP discovery failed: %v\n", discovery.err)
	}
	if wifiCheck.err != nil {
		fmt.Fprintf(stderr, "warning: cannot determine Wi-Fi security: %v\n", wifiCheck.err)
	}

	descriptions, descriptionErrors := fetchDescriptions(ctx, discovery.responses)
	for _, err := range descriptionErrors {
		fmt.Fprintf(stderr, "warning: cannot read a device description: %v\n", err)
	}
	incomplete := arp.err != nil || gateway.err != nil || discovery.err != nil || len(descriptionErrors) > 0
	wifiFs := wifiFindings(wifiCheck.security, lang)
	if incomplete && !*jsonOutput {
		fmt.Fprintln(stderr, "warning: some checks could not run; this result may be incomplete")
	}

	hosts := collectHosts(arp.neighbors, gatewayIP, discovery.responses)
	opened := probe.Scan(ctx, hosts, []int{23, 2323, 21, 5900}, 700*time.Millisecond, 64)
	canReport, forceIncomplete := reportable(scanObservedNetwork(arp.neighbors, discovery.responses, opened), wifiFs)
	if forceIncomplete {
		if !incomplete && !*jsonOutput {
			fmt.Fprintln(stderr, "warning: some checks could not run; this result may be incomplete")
		}
		incomplete = true
	}
	if !canReport {
		fmt.Fprintln(stderr, "error: the scan could not see any device on your network; the result would not be reliable")
		return finish(1)
	}

	store, existed, err := baseline.Load(*statePath)
	firstRun := !existed
	loadState := historyLoadReady
	if err != nil {
		firstRun = true
		store = baseline.Store{Version: 1, Devices: make(map[string]baseline.Seen)}
		if errors.Is(err, baseline.ErrInvalidJSON) {
			backupPath, quarantineErr := baseline.Quarantine(*statePath)
			if quarantineErr != nil {
				loadState = historyLoadUnavailable
				fmt.Fprintf(stderr, "warning: corrupt device history could not be moved to %s; history will not be saved: %v\n", backupPath, quarantineErr)
			} else {
				loadState = historyLoadQuarantined
				fmt.Fprintf(stderr, "warning: corrupt device history moved to %s; starting a new history\n", backupPath)
			}
		} else {
			loadState = historyLoadUnavailable
			fmt.Fprintf(stderr, "warning: cannot read device history; history will not be saved this run: %v\n", err)
		}
	}

	neighborByIP := make(map[netip.Addr]lan.Neighbor, len(arp.neighbors))
	for _, neighbor := range arp.neighbors {
		neighborByIP[neighbor.IP] = neighbor
	}
	deviceNames := make(map[netip.Addr]string, len(hosts))
	newDevices := make(map[netip.Addr]bool)
	now := time.Now().UTC()
	for _, host := range hosts {
		neighbor := neighborByIP[host]
		name := deviceName(lang, host, neighbor, gatewayIP, descriptions)
		deviceNames[host] = name
		if neighbor.MAC == "" {
			continue
		}
		isNew := store.Observe(neighbor.MAC, host.String(), name, now)
		newDevices[host] = !firstRun && isNew
	}

	devices := make([]deviceOutput, 0, len(hosts))
	for _, host := range hosts {
		neighbor := neighborByIP[host]
		devices = append(devices, deviceOutput{
			IP: host.String(), MAC: neighbor.MAC, Name: deviceNames[host], New: newDevices[host],
		})
	}

	findings := makeFindings(opened, discovery.responses, gatewayIP, descriptions, deviceNames, newDevices, lang)
	findings = append(findings, wifiFs...)
	report.Sort(findings)
	status := report.Status(report.Score(findings), findings)
	result := output{
		Version: version, Score: report.Score(findings), Status: status, FirstRun: firstRun, Incomplete: incomplete,
		Gateway: gatewayString(gatewayIP), Devices: devices, Findings: findings,
	}
	if err := writeOutput(stdout, result, *jsonOutput, lang); err != nil {
		if *jsonOutput {
			fmt.Fprintf(stderr, "warning: cannot write JSON output: %v\n", err)
		} else {
			fmt.Fprintf(stderr, "warning: cannot write output: %v\n", err)
		}
		return finish(1)
	}
	if shouldSaveHistory(true, firstRun, store, loadState) {
		if err := baseline.Save(*statePath, store); err != nil {
			fmt.Fprintf(stderr, "warning: cannot save device history: %v\n", err)
			return finish(1)
		}
	}

	for _, finding := range findings {
		if finding.Level == report.Danger {
			return finish(2)
		}
	}
	return finish(0)
}

func scanObservedNetwork(neighbors []lan.Neighbor, responses []ssdp.Response, opened []probe.Open) bool {
	return len(neighbors) > 0 || len(responses) > 0 || len(opened) > 0
}

func reportable(observed bool, wifiFindings []report.Finding) (bool, bool) {
	if observed {
		return true, false
	}
	if len(wifiFindings) > 0 {
		return true, true
	}
	return false, false
}

func shouldSaveHistory(outputSucceeded, firstRun bool, store baseline.Store, loadState historyLoadState) bool {
	if !outputSucceeded || loadState == historyLoadUnavailable {
		return false
	}
	return !firstRun || len(store.Devices) > 0
}

func writeOutput(writer io.Writer, result output, jsonMode bool, lang string) error {
	var rendered bytes.Buffer
	if jsonMode {
		encoder := json.NewEncoder(&rendered)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			return err
		}
	} else {
		renderText(&rendered, result, lang)
	}
	written, err := writer.Write(rendered.Bytes())
	if err != nil {
		return err
	}
	if written != rendered.Len() {
		return io.ErrShortWrite
	}
	return nil
}

func resolveLanguage(flagValue string, getenv func(string) string) (string, error) {
	if flagValue == "auto" {
		return i18n.Detect(getenv), nil
	}
	if lang, ok := i18n.Match(flagValue); ok {
		return lang, nil
	}
	return "", fmt.Errorf("invalid --lang %q (use auto, or one of %s)", flagValue, strings.Join(i18n.Supported(), ", "))
}

func gatewayString(ip netip.Addr) string {
	if !ip.IsValid() {
		return ""
	}
	return ip.String()
}

func collectHosts(neighbors []lan.Neighbor, gateway netip.Addr, responses []ssdp.Response) []netip.Addr {
	hosts := make(map[netip.Addr]struct{})
	add := func(ip netip.Addr) {
		if ip.Is4() && ip.IsPrivate() && len(hosts) < 1024 {
			hosts[ip] = struct{}{}
		}
	}
	for _, neighbor := range neighbors {
		add(neighbor.IP)
	}
	add(gateway)
	for _, response := range responses {
		add(response.IP)
	}
	result := make([]netip.Addr, 0, len(hosts))
	for ip := range hosts {
		result = append(result, ip)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Less(result[j]) })
	return result
}

func fetchDescriptions(ctx context.Context, responses []ssdp.Response) (map[netip.Addr]ssdp.Device, []error) {
	selected := make(map[netip.Addr]ssdp.Response)
	for _, response := range responses {
		if !ssdp.AllowedLocation(response.Location, response.IP) {
			continue
		}
		current, exists := selected[response.IP]
		if !exists || (response.ST == "upnp:rootdevice" && current.ST != "upnp:rootdevice") {
			selected[response.IP] = response
		}
	}

	clients := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	type request struct {
		ip       netip.Addr
		location string
	}
	requests := make(chan request)
	results := make(chan descriptionResult)
	var workers sync.WaitGroup
	workerCount := 8
	if len(selected) < workerCount {
		workerCount = len(selected)
	}
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			for item := range requests {
				device, err := ssdp.FetchDescription(ctx, clients, item.location)
				results <- descriptionResult{ip: item.ip, device: device, err: err}
			}
		}()
	}
	go func() {
		for ip, response := range selected {
			requests <- request{ip: ip, location: response.Location}
		}
		close(requests)
		workers.Wait()
		close(results)
	}()

	devices := make(map[netip.Addr]ssdp.Device, len(selected))
	var errs []error
	for result := range results {
		if result.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", result.ip, result.err))
			continue
		}
		devices[result.ip] = result.device
	}
	return devices, errs
}

func deviceName(lang string, ip netip.Addr, neighbor lan.Neighbor, gateway netip.Addr, descriptions map[netip.Addr]ssdp.Device) string {
	if device, ok := descriptions[ip]; ok {
		if name := sanitizeName(device.FriendlyName); name != "" {
			return name
		}
		name := sanitizeName(sanitizeName(device.Manufacturer) + " " + sanitizeName(device.ModelName))
		if name != "" {
			return name
		}
	}
	if gateway.IsValid() && ip == gateway {
		return i18n.T(lang, "router")
	}
	if neighbor.MAC != "" && lan.IsRandomizedMAC(neighbor.MAC) {
		return i18n.T(lang, "randomized_device")
	}
	return ip.String()
}

func sanitizeName(s string) string {
	filtered := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || r >= 0x200B && r <= 0x200F ||
			r >= 0x202A && r <= 0x202E || r >= 0x2066 && r <= 0x2069 || r == 0xFEFF {
			return -1
		}
		return r
	}, s)
	name := strings.Join(strings.Fields(filtered), " ")
	runes := []rune(name)
	if len(runes) > 64 {
		return string(runes[:64])
	}
	return name
}

func makeFindings(
	opened []probe.Open,
	responses []ssdp.Response,
	gateway netip.Addr,
	descriptions map[netip.Addr]ssdp.Device,
	deviceNames map[netip.Addr]string,
	newDevices map[netip.Addr]bool,
	lang string,
) []report.Finding {
	findings := make([]report.Finding, 0)
	telnetSeen := make(map[netip.Addr]struct{})
	for _, item := range opened {
		name := deviceNames[item.IP]
		switch item.Port {
		case 23, 2323:
			if _, exists := telnetSeen[item.IP]; exists {
				continue
			}
			telnetSeen[item.IP] = struct{}{}
			findings = append(findings, report.Finding{ID: "telnet_open", Level: report.Danger, Points: 25, IP: item.IP.String(), Device: name})
		case 21:
			findings = append(findings, report.Finding{ID: "ftp_open", Level: report.Caution, Points: 8, IP: item.IP.String(), Device: name})
		case 5900:
			findings = append(findings, report.Finding{ID: "vnc_open", Level: report.Caution, Points: 8, IP: item.IP.String(), Device: name})
		}
	}
	for ip, isNew := range newDevices {
		if isNew {
			findings = append(findings, report.Finding{ID: "new_device", Level: report.Caution, Points: 3, IP: ip.String(), Device: deviceNames[ip]})
		}
	}

	isIGD := func(response ssdp.Response) bool {
		if ssdp.IsGateway(response) {
			return true
		}
		if device, ok := descriptions[response.IP]; ok && strings.Contains(strings.ToLower(device.DeviceType), "internetgatewaydevice") {
			return true
		}
		return false
	}
	var upnpIP string
	var upnpDevice string
	if gateway.IsValid() {
		if device, ok := descriptions[gateway]; ok && strings.Contains(strings.ToLower(device.DeviceType), "internetgatewaydevice") {
			upnpIP = gateway.String()
			upnpDevice = deviceNames[gateway]
		}
	}
	for _, response := range responses {
		if gateway.IsValid() && response.IP != gateway {
			continue
		}
		if upnpIP == "" && isIGD(response) {
			upnpIP = response.IP.String()
			upnpDevice = deviceNames[response.IP]
			if !gateway.IsValid() {
				break
			}
		}
	}
	if upnpIP != "" {
		findings = append(findings, report.Finding{ID: "router_upnp", Level: report.Caution, Points: 5, IP: upnpIP, Device: upnpDevice})
	}

	for i := range findings {
		finding := &findings[i]
		finding.Message = i18n.T(lang, finding.ID)
		finding.Fix = i18n.T(lang, finding.ID+"_fix")
		switch finding.ID {
		case "telnet_open", "ftp_open", "vnc_open", "new_device":
			finding.Message = i18n.T(lang, finding.ID, finding.Device)
		}
	}
	return findings
}

func wifiFindings(security wifi.Security, lang string) []report.Finding {
	id := ""
	level := report.Level("")
	points := 0
	switch security {
	case wifi.Open:
		id, level, points = "wifi_open", report.Danger, 25
	case wifi.WEP:
		id, level, points = "wifi_wep", report.Danger, 25
	case wifi.WPA1orTKIP:
		id, level, points = "wifi_weak", report.Caution, 8
	default:
		return nil
	}
	return []report.Finding{{
		ID: id, Level: level, Points: points,
		Message: i18n.T(lang, id), Fix: i18n.T(lang, id+"_fix"),
	}}
}

func renderText(writer io.Writer, result output, lang string) {
	fmt.Fprintln(writer, i18n.T(lang, "title"))
	fmt.Fprintln(writer, i18n.T(lang, "score", i18n.T(lang, result.Status), result.Score))
	newCount := 0
	for _, device := range result.Devices {
		if device.New {
			newCount++
		}
	}
	fmt.Fprintln(writer, i18n.T(lang, "device_counts", len(result.Devices), newCount))
	if result.FirstRun {
		fmt.Fprintf(writer, "\n%s\n", i18n.T(lang, "first_run"))
	}
	if len(result.Findings) == 0 {
		fmt.Fprintf(writer, "\n%s\n", i18n.T(lang, "no_problems"))
	} else {
		fmt.Fprintf(writer, "\n%s\n", i18n.T(lang, "findings_heading"))
		for _, finding := range result.Findings {
			level := i18n.T(lang, string(finding.Level))
			fmt.Fprintf(writer, "  [%s] %s", level, finding.Message)
			if finding.IP != "" {
				fmt.Fprintf(writer, " (%s)", finding.IP)
			}
			fmt.Fprintf(writer, "\n         → %s  %s\n", finding.Fix, i18n.T(lang, "points_recovered", finding.Points))
		}
	}
	fmt.Fprintf(writer, "\n%s\n", i18n.T(lang, "devices_heading"))
	for _, device := range result.Devices {
		newLabel := ""
		if device.New {
			newLabel = " [" + i18n.T(lang, "new_badge") + "]"
		}
		fmt.Fprintf(writer, "  %-15s  %-17s  %s%s\n", device.IP, device.MAC, device.Name, newLabel)
	}
}
