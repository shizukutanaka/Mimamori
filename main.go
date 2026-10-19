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
	"github.com/shizukutanaka/Mimamori/internal/mdns"
	"github.com/shizukutanaka/Mimamori/internal/netbios"
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
	Version              string           `json:"version"`
	Score                int              `json:"score"`
	Status               string           `json:"status"`
	FirstRun             bool             `json:"first_run"`
	Incomplete           bool             `json:"incomplete"`
	HistoryUnavailable   bool             `json:"history_unavailable,omitempty"`
	WiFiLocationRequired bool             `json:"wifi_location_required,omitempty"`
	Checks               checkedSummary   `json:"checked"`
	Gateway              string           `json:"gateway"`
	Devices              []deviceOutput   `json:"devices"`
	Findings             []report.Finding `json:"findings"`
}

type checkedSummary struct {
	Ports        bool `json:"ports"`
	PortForwards bool `json:"port_forwards"`
	WiFi         bool `json:"wifi"`
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

func earlyRunResult(code int, lang string, jsonMode bool) runResult {
	return runResult{code: code, lang: lang, pause: !jsonMode}
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
	device ssdp.Description
	err    error
}

func main() {
	restoreQuickEdit := console.DisableQuickEdit()
	result := run(os.Args[1:], os.Stdout, os.Stderr)
	restoreQuickEdit()
	if result.pause && console.ShouldPause() {
		fmt.Fprint(os.Stdout, i18n.T(result.lang, "press_enter"))
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	}
	os.Exit(result.code)
}

func run(args []string, stdout, stderr io.Writer) runResult {
	return runWithEnv(args, stdout, stderr, os.Getenv)
}

func runWithEnv(args []string, stdout, stderr io.Writer, getenv func(string) string) runResult {
	flags := flag.NewFlagSet("mimamori", flag.ContinueOnError)
	flags.SetOutput(stderr)
	registeredFlags := registerRunFlags(flags)
	jsonOutput := registeredFlags.jsonOutput
	language := registeredFlags.language
	statePath := registeredFlags.statePath
	wait := registeredFlags.wait
	showVersion := registeredFlags.showVersion
	pre := preLanguage(args, getenv, flags)
	flags.Lookup("json").Usage = i18n.T(pre, "flag_json")
	languageHelp := i18n.T(pre, "flag_lang", strings.Join(i18n.Supported(), ", "))
	flags.Lookup("lang").Usage = languageHelp
	flags.Lookup("state").Usage = i18n.T(pre, "flag_state")
	flags.Lookup("wait").Usage = i18n.T(pre, "flag_wait")
	flags.Lookup("version").Usage = i18n.T(pre, "flag_version")
	flags.Usage = func() {
		fmt.Fprintln(stderr, i18n.T(pre, "usage_header"))
		flags.PrintDefaults()
	}
	lang := ""
	incomplete := false
	localizedNoticeShown := false
	finish := func(code int) runResult {
		return finishRun(stderr, code, lang, *jsonOutput, incomplete, localizedNoticeShown)
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return earlyRunResult(0, pre, *jsonOutput)
		}
		return earlyRunResult(1, pre, *jsonOutput)
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, i18n.T(pre, "unexpected_argument", flags.Arg(0)))
		return earlyRunResult(1, pre, *jsonOutput)
	}
	if *showVersion {
		fmt.Fprintf(stdout, "Mimamori %s\n", version)
		return earlyRunResult(0, pre, *jsonOutput)
	}
	resolvedLang, err := resolveLanguage(*language, getenv)
	if err != nil {
		fmt.Fprintln(stderr, i18n.T(pre, "invalid_lang", *language, strings.Join(i18n.Supported(), ", ")))
		return earlyRunResult(1, pre, *jsonOutput)
	}
	lang = resolvedLang
	if *wait < 0 {
		fmt.Fprintln(stderr, i18n.T(pre, "wait_negative"))
		return earlyRunResult(1, pre, *jsonOutput)
	}
	if *statePath == "" {
		var err error
		*statePath, err = baseline.DefaultPath()
		if err != nil {
			fmt.Fprintf(stderr, "warning: cannot determine state path: %v\n", err)
			return finish(1)
		}
	}

	printScanningNotice(stderr, lang, *jsonOutput)

	ctx := context.Background()
	arpCh := make(chan discoveryResult, 1)
	gatewayCh := make(chan gatewayResult, 1)
	unicastCh := make(chan netip.Addr, 1)
	ssdpCh := make(chan ssdpResult, 1)
	wifiCh := make(chan wifiResult, 1)
	go func() {
		neighbors, err := lan.ReadARP(ctx)
		arpCh <- discoveryResult{neighbors: neighbors, err: err}
	}()
	go func() {
		gateway, err := lan.DefaultGateway(ctx)
		if err == nil {
			unicastCh <- gateway
		}
		close(unicastCh)
		gatewayCh <- gatewayResult{gateway: gateway, err: err}
	}()
	go func() {
		responses, err := ssdp.Discover(ctx, *wait, unicastCh)
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

	httpClient := newUPnPHTTPClient()
	descriptions, descriptionErrors := fetchDescriptions(ctx, httpClient, discovery.responses, gatewayIP)
	for _, err := range descriptionErrors {
		fmt.Fprintf(stderr, "warning: cannot read a device description: %s\n", stripUnsafeRunes(err.Error()))
	}
	incomplete = arp.err != nil || gateway.err != nil || discovery.err != nil || wifiCheck.err != nil || len(descriptionErrors) > 0
	wifiFs := wifiFindings(wifiCheck.security, lang)
	portMappings, portMappingQueryRan, portMappingErr := fetchPortMappings(ctx, httpClient, gatewayIP, discovery.responses, descriptions)
	if portMappingErr != nil {
		fmt.Fprintf(stderr, "warning: cannot read the router's port forwarding list: %s\n", portMappingWarning(portMappingErr))
		incomplete = true
	}

	var prefixes []netip.Prefix
	prefixesFailed := false
	if gatewayIP.IsValid() {
		var err error
		prefixes, err = lan.LocalPrefixes()
		if err != nil {
			fmt.Fprintf(stderr, "warning: cannot list network interfaces: %v\n", err)
			incomplete = true
			prefixesFailed = true
		}
	}
	scope := scanScope(arp.neighbors, prefixes, gatewayIP)
	if prefixesFailed {
		scope = nil
	}
	initialHosts := collectHosts(arp.neighbors, gatewayIP, discovery.responses, scope)
	scanTargets := initialHosts
	if gatewayIP.IsValid() {
		scanTargets = mergeHosts(initialHosts, scanSweepHosts(prefixes, scope, gatewayIP))
	}
	opened := markSMB1(ctx, probe.Scan(ctx, scanTargets, []int{23, 2323, 21, 5900, 5555, 445, 554, 1883}, 700*time.Millisecond, 128))
	opened = markUnauthenticatedRTSP(ctx, opened)
	opened = markUnauthenticatedMQTT(ctx, opened)
	neighbors := arp.neighbors
	opened = markUnauthenticatedVNC(ctx, opened)
	opened = markUnauthenticatedADB(ctx, opened)
	if len(scanTargets) > len(initialHosts) {
		refreshed, err := lan.ReadARP(ctx)
		if err != nil {
			fmt.Fprintf(stderr, "warning: cannot refresh the neighbor table: %v\n", err)
		} else {
			neighbors = mergeNeighbors(neighbors, refreshed)
		}
	}
	canReport, forceIncomplete := reportable(scanObservedNetwork(neighbors, discovery.responses, opened), wifiFs)
	if forceIncomplete {
		incomplete = true
	}
	if !canReport {
		printNoDeviceError(stderr, lang, wifiCheck.err)
		return finish(1)
	}

	if prefixesFailed {
		prefixes = nil
	}
	neighbors = homeNeighbors(neighbors, prefixes, gatewayIP)

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

	var lastCheckDays *int
	if existed && !firstRun && loadState == historyLoadReady {
		if lastRun, ok := store.LastRun(); ok {
			checkNow := time.Now().UTC()
			days := daysSinceLastCheck(checkNow, lastRun, time.Local)
			lastCheckDays = &days
		}
	}

	reportedHosts := deviceHosts(neighbors, gatewayIP, discovery.responses, opened)
	resolvedNames := lookupDeviceNames(ctx, reportedHosts)
	neighborByIP := make(map[netip.Addr]lan.Neighbor, len(neighbors))
	for _, neighbor := range neighbors {
		neighborByIP[neighbor.IP] = neighbor
	}
	deviceNames := make(map[netip.Addr]string, len(reportedHosts))
	newDeviceFindings := make(map[netip.Addr]string)
	now := time.Now().UTC()
	for _, host := range reportedHosts {
		neighbor := neighborByIP[host]
		name := deviceName(lang, host, neighbor, gatewayIP, descriptions, resolvedNames)
		deviceNames[host] = name
		if neighbor.MAC == "" {
			continue
		}
		isNew := store.Observe(neighbor.MAC, host.String(), name, now)
		if findingID := newDeviceFindingID(firstRun, isNew, neighbor.MAC); findingID != "" {
			newDeviceFindings[host] = findingID
		}
	}
	store.LastCheck = now

	devices := make([]deviceOutput, 0, len(reportedHosts))
	for _, host := range reportedHosts {
		neighbor := neighborByIP[host]
		devices = append(devices, deviceOutput{
			IP: host.String(), MAC: neighbor.MAC, Name: deviceNames[host], New: newDeviceFindings[host] != "",
		})
	}

	findings := makeFindings(opened, discovery.responses, gatewayIP, descriptions, deviceNames, newDeviceFindings, portMappings, lang)
	findings = markAnonymousFTP(findings, anonymousFTPHosts(ctx, opened), lang)
	findings = append(findings, wifiFs...)
	findings = append(findings, gatewaySpoofFindings(gatewayNeighbors(ctx, neighbors, gatewayIP, lan.ReadARP), gatewayIP, lang)...)
	report.Sort(findings)
	findings = append(findings, vulnerableUPnPStackFindings(discovery.responses, gatewayIP, deviceNames, lang)...)
	report.Sort(findings)
	status := report.Status(report.Score(findings), findings)
	result := output{
		Version: version, Score: report.Score(findings), Status: status, FirstRun: firstRun, Incomplete: incomplete,
		HistoryUnavailable:   loadState == historyLoadUnavailable,
		WiFiLocationRequired: errors.Is(wifiCheck.err, wifi.ErrLocationPermission),
		Checks: checkedSummary{
			Ports:        len(scanTargets) > 0,
			PortForwards: gatewayIP.IsValid() && portMappingQueryRan && portMappingErr == nil,
			WiFi:         wifiCheck.err == nil && wifiCheck.security != wifi.Unknown,
		},
		Gateway: gatewayString(gatewayIP), Devices: devices, Findings: findings,
	}
	if err := writeOutput(stdout, result, *jsonOutput, lang, lastCheckDays); err != nil {
		if *jsonOutput {
			fmt.Fprintf(stderr, "warning: cannot write JSON output: %v\n", err)
		} else {
			fmt.Fprintf(stderr, "warning: cannot write output: %v\n", err)
		}
		return finish(1)
	}
	localizedNoticeShown = !*jsonOutput
	if shouldSaveHistory(true, firstRun, store, loadState) {
		store.Prune(now, 30*24*time.Hour, lan.IsRandomizedMAC)
		if err := baseline.Save(*statePath, store); err != nil {
			fmt.Fprintf(stderr, "warning: cannot save device history: %v\n", err)
			return finish(exitCode(findings, true))
		}
	}

	return finish(exitCode(findings, false))
}

func exitCode(findings []report.Finding, historySaveFailed bool) int {
	for _, finding := range findings {
		if finding.Level == report.Danger {
			return 2
		}
	}
	if historySaveFailed {
		return 1
	}
	return 0
}

func printNoDeviceError(writer io.Writer, lang string, wifiErr error) {
	fmt.Fprintln(writer, "error: the scan could not see any device on your network; the result would not be reliable")
	if errors.Is(wifiErr, wifi.ErrLocationPermission) {
		fmt.Fprintln(writer, i18n.T(lang, "wifi_location_notice"))
	}
}

func printScanningNotice(writer io.Writer, lang string, jsonOutput bool) {
	if !jsonOutput {
		fmt.Fprintln(writer, i18n.T(lang, "scanning_notice"))
	}
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

func writeOutput(writer io.Writer, result output, jsonMode bool, lang string, lastCheckDays ...*int) error {
	var rendered bytes.Buffer
	if jsonMode {
		encoder := json.NewEncoder(&rendered)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			return err
		}
	} else {
		renderText(&rendered, result, lang, lastCheckDays...)
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

func finishRun(stderr io.Writer, code int, lang string, jsonMode, incomplete, localizedNoticeShown bool) runResult {
	if incomplete && (jsonMode || !localizedNoticeShown) {
		fmt.Fprintln(stderr, "warning: some checks could not run; this result may be incomplete")
	}
	return runResult{code: code, lang: lang, pause: !jsonMode}
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

func collectHosts(neighbors []lan.Neighbor, gateway netip.Addr, responses []ssdp.Response, scope []netip.Prefix) []netip.Addr {
	hosts := make(map[netip.Addr]struct{})
	add := func(ip netip.Addr) {
		if ip.Is4() && ip.IsPrivate() && len(hosts) < maxScanHosts {
			hosts[ip] = struct{}{}
		}
	}
	inScope := func(ip netip.Addr) bool {
		if len(scope) == 0 {
			return true
		}
		for _, prefix := range scope {
			if prefix.Contains(ip) {
				return true
			}
		}
		return false
	}
	for _, neighbor := range neighbors {
		if inScope(neighbor.IP) {
			add(neighbor.IP)
		}
	}
	add(gateway)
	for _, response := range responses {
		if inScope(response.IP) {
			add(response.IP)
		}
	}
	result := make([]netip.Addr, 0, len(hosts))
	for ip := range hosts {
		result = append(result, ip)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Less(result[j]) })
	return result
}

const maxScanHosts = 1024

func mergeHosts(hosts, swept []netip.Addr) []netip.Addr {
	seen := make(map[netip.Addr]struct{}, len(hosts)+len(swept))
	merged := make([]netip.Addr, 0, maxScanHosts)
	for _, list := range [][]netip.Addr{hosts, swept} {
		for _, host := range list {
			if !host.Is4() || !host.IsPrivate() {
				continue
			}
			if _, exists := seen[host]; exists {
				continue
			}
			if len(merged) == maxScanHosts {
				break
			}
			seen[host] = struct{}{}
			merged = append(merged, host)
		}
		if len(merged) == maxScanHosts {
			break
		}
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Less(merged[j])
	})
	return merged
}

func mergeNeighbors(original, refreshed []lan.Neighbor) []lan.Neighbor {
	seen := make(map[netip.Addr]struct{}, len(original)+len(refreshed))
	merged := make([]lan.Neighbor, 0, len(original)+len(refreshed))
	for _, list := range [][]lan.Neighbor{refreshed, original} {
		for _, neighbor := range list {
			if _, exists := seen[neighbor.IP]; exists {
				continue
			}
			seen[neighbor.IP] = struct{}{}
			merged = append(merged, neighbor)
		}
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].IP.Less(merged[j].IP)
	})
	return merged
}

func deviceHosts(neighbors []lan.Neighbor, gateway netip.Addr, responses []ssdp.Response, opened []probe.Open) []netip.Addr {
	hosts := collectHosts(neighbors, gateway, responses, nil)
	seen := make(map[netip.Addr]struct{}, len(hosts)+len(opened))
	for _, host := range hosts {
		seen[host] = struct{}{}
	}
	for _, item := range opened {
		if !item.IP.Is4() || !item.IP.IsPrivate() {
			continue
		}
		if _, exists := seen[item.IP]; exists {
			continue
		}
		seen[item.IP] = struct{}{}
		hosts = append(hosts, item.IP)
	}
	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].Less(hosts[j])
	})
	return hosts
}

func newUPnPHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func fetchPortMappings(
	ctx context.Context,
	client *http.Client,
	gateway netip.Addr,
	responses []ssdp.Response,
	descriptions map[netip.Addr]ssdp.Description,
) ([]ssdp.PortMapping, bool, error) {
	routerIP, selected := routerUPnPIP(gateway, responses, descriptions)
	if gateway.IsValid() {
		routerIP = gateway
		description, ok := descriptions[routerIP]
		if !ok {
			return nil, false, nil
		}
		isIGD := description.Device.IsIGD() ||
			(selected && routerIP == gateway)
		if !isIGD {
			return nil, false, nil
		}
	} else if !selected {
		return nil, false, nil
	}
	description, ok := descriptions[routerIP]
	if !ok {
		return nil, false, fmt.Errorf("router advertises UPnP but its port forwarding list cannot be queried")
	}
	type serviceEndpoint struct {
		service    ssdp.Service
		controlURL string
	}
	endpoints := make([]serviceEndpoint, 0)
	var errs []error
	for _, service := range description.Device.WANConnections() {
		controlURL, ok := ssdp.ControlURL(description.Location, description.URLBase, service.ControlURL, routerIP)
		if !ok {
			errs = append(errs, fmt.Errorf("%s: control URL is not on the router", service.ServiceType))
			continue
		}
		endpoints = append(endpoints, serviceEndpoint{service: service, controlURL: controlURL})
	}
	if len(endpoints) == 0 {
		return nil, false, fmt.Errorf("router advertises UPnP but its port forwarding list cannot be queried")
	}
	queryContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var mappings []ssdp.PortMapping
	for _, endpoint := range endpoints {
		serviceMappings, err := ssdp.PortMappings(queryContext, client, endpoint.controlURL, endpoint.service.ServiceType)
		mappings = append(mappings, serviceMappings...)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", endpoint.service.ServiceType, err))
		}
	}
	return mappings, true, errors.Join(errs...)
}

var descriptionFetchBudget = 5 * time.Second

func fetchDescriptions(ctx context.Context, client *http.Client, responses []ssdp.Response, gateway netip.Addr) (map[netip.Addr]ssdp.Description, []error) {
	ctx, cancel := context.WithTimeout(ctx, descriptionFetchBudget)
	defer cancel()

	selected := make(map[netip.Addr]ssdp.Response)
	for _, response := range responses {
		if !ssdp.AllowedLocation(response.Location, response.IP) {
			continue
		}
		current, exists := selected[response.IP]
		if !exists || descriptionPriority(response) > descriptionPriority(current) {
			selected[response.IP] = response
		}
	}

	type request struct {
		ip       netip.Addr
		location string
	}
	ordered := make([]request, 0, len(selected))
	for ip, response := range selected {
		ordered = append(ordered, request{ip: ip, location: response.Location})
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].ip.Less(ordered[j].ip)
	})
	for i := range ordered {
		if ordered[i].ip == gateway && i > 0 {
			gatewayRequest := ordered[i]
			copy(ordered[1:i+1], ordered[:i])
			ordered[0] = gatewayRequest
			break
		}
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
				device, err := ssdp.FetchDescription(ctx, client, item.location)
				results <- descriptionResult{ip: item.ip, device: device, err: err}
			}
		}()
	}
	skipped := 0
	go func() {
	feed:
		for index, item := range ordered {
			if ctx.Err() != nil {
				skipped = len(ordered) - index
				break
			}
			select {
			case requests <- item:
			case <-ctx.Done():
				skipped = len(ordered) - index
				break feed
			}
		}
		close(requests)
		workers.Wait()
		close(results)
	}()

	devices := make(map[netip.Addr]ssdp.Description, len(selected))
	var errs []error
	for result := range results {
		if result.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", result.ip, result.err))
			continue
		}
		devices[result.ip] = result.device
	}
	if skipped > 0 {
		errs = append(errs, fmt.Errorf("description fetch time limit reached; %d devices not read", skipped))
	}
	return devices, errs
}

func descriptionPriority(response ssdp.Response) int {
	switch {
	case isUPnPIGDResponse(response.ST):
		return 3
	case ssdp.IsGateway(response):
		return 2
	case response.ST == "upnp:rootdevice":
		return 1
	default:
		return 0
	}
}

func isUPnPIGDResponse(st string) bool {
	st = strings.ToLower(st)
	return strings.HasPrefix(st, "urn:schemas-upnp-org:") &&
		(strings.Contains(st, ":device:internetgatewaydevice:") ||
			strings.Contains(st, ":device:wandevice:") ||
			strings.Contains(st, ":device:wanconnectiondevice:") ||
			strings.Contains(st, ":service:wanipconnection:") ||
			strings.Contains(st, ":service:wanpppconnection:") ||
			strings.Contains(st, ":service:wancommoninterfaceconfig:") ||
			strings.Contains(st, ":service:layer3forwarding:"))
}

func lookupDeviceNames(ctx context.Context, hosts []netip.Addr) map[netip.Addr]string {
	lookupCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var mdnsNames, netbiosNames map[netip.Addr]string
	var lookups sync.WaitGroup
	lookups.Add(2)
	go func() {
		defer lookups.Done()
		mdnsNames = mdns.ReverseNames(lookupCtx, hosts, time.Second)
	}()
	go func() {
		defer lookups.Done()
		netbiosNames = netbios.Names(ctx, hosts, time.Second)
	}()
	lookups.Wait()
	return mergeDeviceNames(mdnsNames, netbiosNames)
}

func mergeDeviceNames(mdnsNames, netbiosNames map[netip.Addr]string) map[netip.Addr]string {
	resolved := make(map[netip.Addr]string, len(mdnsNames)+len(netbiosNames))
	for ip, name := range netbiosNames {
		if sanitizeName(name) != "" {
			resolved[ip] = name
		}
	}
	for ip, name := range mdnsNames {
		if sanitizeName(name) != "" {
			resolved[ip] = name
		}
	}
	return resolved
}

func deviceName(
	lang string,
	ip netip.Addr,
	neighbor lan.Neighbor,
	gateway netip.Addr,
	descriptions map[netip.Addr]ssdp.Description,
	resolvedNames map[netip.Addr]string,
) string {
	if description, ok := descriptions[ip]; ok {
		if name := sanitizeName(description.Device.FriendlyName); name != "" {
			return name
		}
		name := sanitizeName(sanitizeName(description.Device.Manufacturer) + " " + sanitizeName(description.Device.ModelName))
		if name != "" {
			return name
		}
	}
	if gateway.IsValid() && ip == gateway {
		return i18n.T(lang, "router")
	}
	if name := sanitizeName(resolvedNames[ip]); name != "" {
		return name
	}
	if neighbor.MAC != "" && lan.IsRandomizedMAC(neighbor.MAC) {
		return i18n.T(lang, "randomized_device")
	}
	return ip.String()
}

func portMappingWarning(err error) string {
	return stripUnsafeRunes(strings.ReplaceAll(err.Error(), "\n", "; "))
}

func stripUnsafeRunes(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.Is(unicode.Bidi_Control, r) ||
			r >= 0x200B && r <= 0x200D || r == 0xFEFF {
			return -1
		}
		return r
	}, s)
}

func sanitizeName(s string) string {
	filtered := stripUnsafeRunes(s)
	name := strings.Join(strings.Fields(filtered), " ")
	runes := []rune(name)
	if len(runes) > 64 {
		return string(runes[:64])
	}
	return name
}

func newDeviceFindingID(firstRun, isNew bool, mac string) string {
	if firstRun || !isNew {
		return ""
	}
	if lan.IsRandomizedMAC(mac) {
		return "new_device_private"
	}
	return "new_device"
}

func makeFindings(
	opened []probe.Open,
	responses []ssdp.Response,
	gateway netip.Addr,
	descriptions map[netip.Addr]ssdp.Description,
	deviceNames map[netip.Addr]string,
	newDeviceFindings map[netip.Addr]string,
	portMappings []ssdp.PortMapping,
	lang string,
) []report.Finding {
	findings := make([]report.Finding, 0)
	telnetSeen := make(map[netip.Addr]struct{})
	for _, item := range opened {
		name := deviceNames[item.IP]
		switch item.Port {
		case smb1Port:
			findings = append(findings, report.Finding{ID: "smb1_enabled", Level: report.Caution, Points: 8, IP: item.IP.String(), Device: name})
		case 1883:
			findings = append(findings, report.Finding{ID: "mqtt_anonymous", Level: report.Caution, Points: 8, IP: item.IP.String(), Device: name})
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
		case vncNoAuthPort:
			findings = append(findings, report.Finding{ID: "vnc_noauth", Level: report.Danger, Points: 25, IP: item.IP.String(), Device: name})
		case 5555:
			findings = append(findings, report.Finding{ID: "adb_open", Level: report.Caution, Points: 8, IP: item.IP.String(), Device: name})
		case 554:
			findings = append(findings, report.Finding{ID: "rtsp_noauth", Level: report.Caution, Points: 8, IP: item.IP.String(), Device: name})
		case adbNoAuthPort:
			findings = append(findings, report.Finding{ID: "adb_noauth", Level: report.Danger, Points: 25, IP: item.IP.String(), Device: name})
		}
	}
	for ip, findingID := range newDeviceFindings {
		if findingID != "" {
			findings = append(findings, report.Finding{ID: findingID, Level: report.Caution, Points: 3, IP: ip.String(), Device: deviceNames[ip]})
		}
	}

	type portForwardKey struct {
		client netip.Addr
		port   int
	}
	type portForwardGroup struct {
		client netip.Addr
		port   int
		device string
		id     string
		level  report.Level
		points int
	}
	groupIndexes := make(map[portForwardKey]int)
	groups := make([]portForwardGroup, 0, len(portMappings))
	for _, mapping := range portMappings {
		if !mapping.Enabled {
			continue
		}
		internalIP, err := netip.ParseAddr(mapping.InternalClient)
		if err != nil || !internalIP.Is4() || !internalIP.IsPrivate() {
			continue
		}
		key := portForwardKey{client: internalIP, port: mapping.ExternalPort}
		device := deviceNames[internalIP]
		if device == "" {
			device = sanitizeName(mapping.Description)
		}
		if device == "" {
			device = internalIP.String()
		}
		group := portForwardGroup{
			client: internalIP, port: mapping.ExternalPort, device: device,
			id: "port_forward", level: report.Caution, points: 3,
		}
		if riskyPortForward(mapping) {
			group.id, group.level, group.points = "port_forward_risky", report.Danger, 20
			switch mapping.InternalPort {
			case 554, 8554, 37777, 34567:
				group.id = "forward_camera"
			case 1433, 3306, 5432, 6379, 9200, 27017:
				group.id = "forward_database"
			}
		}
		if index, exists := groupIndexes[key]; exists {
			if group.points > groups[index].points {
				groups[index].id = group.id
				groups[index].level = group.level
				groups[index].points = group.points
				groups[index].device = group.device
			}
			continue
		}
		groupIndexes[key] = len(groups)
		groups = append(groups, group)
	}
	for _, group := range groups {
		findings = append(findings, report.Finding{
			ID: group.id, Level: group.level, Points: group.points, IP: group.client.String(),
			Device: group.device, Port: group.port,
		})
	}

	if upnpIP, ok := routerUPnPIP(gateway, responses, descriptions); ok {
		findings = append(findings, report.Finding{
			ID: "router_upnp", Level: report.Caution, Points: 5,
			IP: upnpIP.String(), Device: deviceNames[upnpIP],
		})
	}

	for i := range findings {
		finding := &findings[i]
		finding.Message = i18n.T(lang, finding.ID)
		finding.Fix = i18n.T(lang, finding.ID+"_fix")
		if finding.ID == "mqtt_anonymous" {
			finding.Message = i18n.T(lang, finding.ID, finding.Device)
		}
		switch finding.ID {
		case "telnet_open", "ftp_open", "vnc_open", "adb_open", "new_device", "new_device_private":
			finding.Message = i18n.T(lang, finding.ID, finding.Device)
		case "port_forward", "port_forward_risky", "forward_camera", "forward_database":
			finding.Message = i18n.T(lang, finding.ID, finding.Device, finding.Port)
		case "smb1_enabled":
			finding.Message = i18n.T(lang, finding.ID, finding.Device)
		}
		if finding.ID == "vnc_noauth" {
			finding.Message = i18n.T(lang, finding.ID, finding.Device)
		}
		if finding.ID == "adb_noauth" {
			finding.Message = i18n.T(lang, finding.ID, finding.Device)
		}
	}
	for i := range findings {
		finding := &findings[i]
		if finding.ID == "rtsp_noauth" {
			finding.Message = i18n.T(lang, finding.ID, finding.Device)
		}
	}
	return findings
}

func riskyPortForward(mapping ssdp.PortMapping) bool {
	if !strings.EqualFold(mapping.Protocol, "TCP") {
		return false
	}
	remoteHost := strings.TrimSpace(mapping.RemoteHost)
	if remoteHost != "" && remoteHost != "*" && remoteHost != "0.0.0.0" {
		return false
	}
	switch mapping.InternalPort {
	case 21, 22, 23, 135, 139, 445, 554, 8554, 1883, 2323, 2375, 3389, 5555, 5800, 5900, 1433, 3306, 5432, 6379, 9200, 27017, 34567, 37777:
		return true
	default:
		return false
	}
}

func routerUPnPIP(gateway netip.Addr, responses []ssdp.Response, descriptions map[netip.Addr]ssdp.Description) (netip.Addr, bool) {
	isIGDDevice := func(device ssdp.Device) bool {
		return device.IsIGD()
	}
	isIGDResponse := func(response ssdp.Response) bool {
		if ssdp.IsGateway(response) {
			return true
		}
		if description, ok := descriptions[response.IP]; ok && isIGDDevice(description.Device) {
			return true
		}
		return false
	}
	if gateway.IsValid() {
		if description, ok := descriptions[gateway]; ok && isIGDDevice(description.Device) {
			return gateway, true
		}
	}
	for _, response := range responses {
		if gateway.IsValid() && response.IP != gateway {
			continue
		}
		if isIGDResponse(response) {
			return response.IP, true
		}
	}
	return netip.Addr{}, false
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

func renderText(writer io.Writer, result output, lang string, lastCheckDays ...*int) {
	fmt.Fprintln(writer, i18n.T(lang, "title"))
	fmt.Fprintln(writer, i18n.T(lang, "score", i18n.T(lang, result.Status), result.Score))
	checkedItems := []string{
		i18n.T(lang, "checked_devices", i18n.CheckedDeviceCount(lang, len(result.Devices))),
	}
	if result.Checks.Ports {
		checkedItems = append(checkedItems, i18n.T(lang, "checked_ports"))
	}
	if result.Checks.PortForwards {
		checkedItems = append(checkedItems, i18n.T(lang, "checked_forwards"))
	}
	if result.Checks.WiFi {
		checkedItems = append(checkedItems, i18n.T(lang, "checked_wifi"))
	}
	fmt.Fprintln(writer, strings.Join(checkedItems, i18n.T(lang, "checked_sep")))
	if result.Incomplete {
		fmt.Fprintln(writer, i18n.T(lang, "incomplete_notice"))
	}
	if result.WiFiLocationRequired {
		fmt.Fprintln(writer, i18n.T(lang, "wifi_location_notice"))
	}
	newCount := 0
	for _, device := range result.Devices {
		if device.New {
			newCount++
		}
	}
	fmt.Fprintln(writer, i18n.T(lang, "device_counts", len(result.Devices), newCount))
	var lastCheckDayCount *int
	if len(lastCheckDays) > 0 {
		lastCheckDayCount = lastCheckDays[0]
	}
	if line, ok := lastCheckLine(result, lastCheckDayCount, lang); ok {
		fmt.Fprintln(writer, line)
	}
	if result.FirstRun {
		key := "first_run"
		if result.HistoryUnavailable {
			key = "history_unavailable"
		}
		fmt.Fprintf(writer, "\n%s\n", i18n.T(lang, key))
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
