package ssdp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/textproto"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type Response struct {
	IP       netip.Addr
	ST       string
	USN      string
	Server   string
	Location string
}

func ParseResponse(b []byte, from netip.Addr) (Response, bool) {
	reader := bufio.NewReader(bytes.NewReader(b))
	status, err := reader.ReadString('\n')
	if err != nil || !strings.HasPrefix(strings.TrimSpace(status), "HTTP/1.1 200") {
		return Response{}, false
	}
	headers, err := textproto.NewReader(reader).ReadMIMEHeader()
	if err != nil {
		return Response{}, false
	}
	st := headers.Get("ST")
	if st == "" {
		st = headers.Get("NT")
	}
	return Response{
		IP:       from,
		ST:       st,
		USN:      headers.Get("USN"),
		Server:   headers.Get("SERVER"),
		Location: headers.Get("LOCATION"),
	}, true
}

func Discover(ctx context.Context, wait time.Duration, unicast <-chan netip.Addr) ([]Response, error) {
	multicast, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return nil, fmt.Errorf("resolve SSDP multicast address: %w", err)
	}

	var targetChannel <-chan *net.UDPAddr
	var cancel context.CancelFunc
	var converter sync.WaitGroup
	if unicast != nil {
		targets := make(chan *net.UDPAddr)
		conversionContext, conversionCancel := context.WithCancel(ctx)
		cancel = conversionCancel
		converter.Add(1)
		go func() {
			defer converter.Done()
			defer close(targets)
			for {
				select {
				case <-conversionContext.Done():
					return
				case address, ok := <-unicast:
					if !ok {
						return
					}
					if !address.IsValid() || !address.Is4() || !address.IsPrivate() || !address.IsGlobalUnicast() ||
						address.IsUnspecified() || address.IsMulticast() {
						continue
					}
					select {
					case targets <- &net.UDPAddr{IP: net.IP(address.AsSlice()), Port: 1900}:
					case <-conversionContext.Done():
						return
					}
				}
			}
		}()
		targetChannel = targets
	}

	responses, err := discoverTargets(
		ctx,
		wait,
		multicast,
		targetChannel,
		func(ip netip.Addr) bool { return ip.Is4() && ip.IsPrivate() },
		300*time.Millisecond,
	)
	if cancel != nil {
		cancel()
		converter.Wait()
	}
	return responses, err
}

func discover(
	ctx context.Context,
	wait time.Duration,
	multicast *net.UDPAddr,
	unicast []*net.UDPAddr,
	accept func(netip.Addr) bool,
	resend time.Duration,
) ([]Response, error) {
	targets := make(chan *net.UDPAddr, len(unicast))
	for _, target := range unicast {
		targets <- target
	}
	close(targets)
	return discoverTargets(ctx, wait, multicast, targets, accept, resend)
}

var afterRetrySend = func() {}

func discoverTargets(
	ctx context.Context,
	wait time.Duration,
	multicast *net.UDPAddr,
	unicast <-chan *net.UDPAddr,
	accept func(netip.Addr) bool,
	resend time.Duration,
) ([]Response, error) {
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("listen for SSDP: %w", err)
	}
	defer conn.Close()

	if accept == nil {
		accept = func(ip netip.Addr) bool { return ip.Is4() && ip.IsPrivate() }
	}

	start := time.Now()
	deadline := start.Add(wait)
	if wait <= 0 {
		deadline = start
	}
	if contextDeadline, ok := ctx.Deadline(); ok && contextDeadline.Before(deadline) {
		deadline = contextDeadline
	}
	retryAt := start.Add(resend)
	retryPending := resend > 0 && retryAt.Before(deadline)

	var sendMu sync.Mutex
	var sendErrors []error
	sentSuccessfully := false
	send := func(target *net.UDPAddr, message []byte) {
		if target == nil {
			return
		}
		_, writeErr := conn.WriteTo(message, target)
		sendMu.Lock()
		defer sendMu.Unlock()
		if writeErr != nil {
			sendErrors = append(sendErrors, fmt.Errorf("%s: %w", target, writeErr))
			return
		}
		sentSuccessfully = true
	}

	sendMulticast := func() {
		if multicast != nil {
			send(multicast, ssdpSearchMessage("239.255.255.250:1900", true))
		}
	}
	sendMulticast()

	readDeadline := deadline
	if retryPending {
		readDeadline = retryAt
	}
	setReadDeadline := func(deadline time.Time) {
		_ = conn.SetReadDeadline(deadline)
		if ctx.Err() != nil {
			_ = conn.SetReadDeadline(time.Now())
		}
	}
	setReadDeadline(readDeadline)
	stopCancel := context.AfterFunc(ctx, func() {
		_ = conn.SetReadDeadline(time.Now())
	})
	defer stopCancel()

	var workers sync.WaitGroup
	workers.Add(1)
	go func() {
		defer workers.Done()
		pending := make([]struct {
			target *net.UDPAddr
			at     time.Time
		}, 0)
		deadlineTimer := time.NewTimer(time.Until(deadline))
		defer deadlineTimer.Stop()
		for {
			var retryTimer *time.Timer
			var retryChannel <-chan time.Time
			if len(pending) > 0 {
				retryTimer = time.NewTimer(time.Until(pending[0].at))
				retryChannel = retryTimer.C
			}
			select {
			case <-ctx.Done():
				if retryTimer != nil {
					retryTimer.Stop()
				}
				return
			case <-deadlineTimer.C:
				if retryTimer != nil {
					retryTimer.Stop()
				}
				return
			case target, ok := <-unicast:
				if retryTimer != nil {
					retryTimer.Stop()
				}
				if !ok {
					unicast = nil
					if len(pending) == 0 {
						return
					}
					continue
				}
				if target == nil || !time.Now().Before(deadline) {
					continue
				}
				send(target, ssdpSearchMessage(target.IP.String()+":1900", false))
				if resend > 0 {
					at := time.Now().Add(resend)
					if at.Before(deadline) {
						pending = append(pending, struct {
							target *net.UDPAddr
							at     time.Time
						}{target: target, at: at})
					}
				}
			case <-retryChannel:
				if retryTimer != nil {
					retryTimer.Stop()
				}
				now := time.Now()
				due := pending[:0]
				for _, item := range pending {
					if !item.at.After(now) {
						send(item.target, ssdpSearchMessage(item.target.IP.String()+":1900", false))
					} else {
						due = append(due, item)
					}
				}
				pending = due
			}
			sort.Slice(pending, func(i, j int) bool { return pending[i].at.Before(pending[j].at) })
		}
	}()

	responses := make([]Response, 0)
	seen := make(map[string]struct{})
	buffer := make([]byte, 2048)
	retried := false
	for {
		n, sender, err := conn.ReadFrom(buffer)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				if retryPending && !retried && ctx.Err() == nil {
					retried = true
					sendMulticast()
					afterRetrySend()
					setReadDeadline(deadline)
					continue
				}
				break
			}
			if ctx.Err() != nil {
				break
			}
			return nil, fmt.Errorf("read SSDP response: %w", err)
		}
		ip, ok := netip.AddrFromSlice(sender.(*net.UDPAddr).IP)
		if !ok {
			continue
		}
		ip = ip.Unmap()
		if !accept(ip) {
			continue
		}
		response, ok := ParseResponse(buffer[:n], ip)
		if !ok {
			continue
		}
		key := ip.String() + "\x00" + response.USN + "\x00" + response.ST
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		responses = append(responses, response)
	}
	workers.Wait()
	sendMu.Lock()
	defer sendMu.Unlock()
	if !sentSuccessfully {
		return nil, sendSSDPError(sendErrors)
	}
	return responses, nil
}

func ssdpSearchMessage(host string, multicast bool) []byte {
	var message strings.Builder
	message.WriteString("M-SEARCH * HTTP/1.1\r\n")
	message.WriteString("HOST: ")
	message.WriteString(host)
	message.WriteString("\r\nMAN: \"ssdp:discover\"\r\n")
	if multicast {
		message.WriteString("MX: 2\r\n")
	}
	message.WriteString("ST: ssdp:all\r\n\r\n")
	return []byte(message.String())
}

func sendSSDPError(sendErrors []error) error {
	if len(sendErrors) == 0 {
		return errors.New("send SSDP search: no targets")
	}
	return fmt.Errorf("send SSDP search: %w", errors.Join(sendErrors...))
}

type Service struct {
	ServiceType string `xml:"serviceType"`
	ControlURL  string `xml:"controlURL"`
}

type Device struct {
	FriendlyName string    `xml:"friendlyName"`
	Manufacturer string    `xml:"manufacturer"`
	ModelName    string    `xml:"modelName"`
	DeviceType   string    `xml:"deviceType"`
	Services     []Service `xml:"serviceList>service"`
	Devices      []Device  `xml:"deviceList>device"`
}

func (d Device) IsIGD() bool {
	if strings.Contains(strings.ToLower(d.DeviceType), "internetgatewaydevice") {
		return true
	}
	for _, device := range d.Devices {
		if device.IsIGD() {
			return true
		}
	}
	return false
}

type Description struct {
	Device   Device `xml:"device"`
	URLBase  string `xml:"URLBase"`
	Location string `xml:"-"`
}

func (d Device) WANConnection() (Service, bool) {
	for _, service := range d.Services {
		if strings.Contains(service.ServiceType, ":service:WANIPConnection:") ||
			strings.Contains(service.ServiceType, ":service:WANPPPConnection:") {
			return service, true
		}
	}
	for _, device := range d.Devices {
		if service, ok := device.WANConnection(); ok {
			return service, true
		}
	}
	return Service{}, false
}

func ParseDescription(r io.Reader) (Description, error) {
	var description Description
	if err := xml.NewDecoder(r).Decode(&description); err != nil {
		return Description{}, err
	}
	return description, nil
}

func AllowedLocation(location string, from netip.Addr) bool {
	if !from.Is4() || !from.IsPrivate() {
		return false
	}
	parsed, err := url.Parse(location)
	if err != nil || !strings.EqualFold(parsed.Scheme, "http") || parsed.User != nil {
		return false
	}
	host, err := netip.ParseAddr(parsed.Hostname())
	return err == nil && host.Is4() && host.IsPrivate() && host == from
}

func ControlURL(location, urlBase, controlURL string, from netip.Addr) (string, bool) {
	if strings.TrimSpace(controlURL) == "" {
		return "", false
	}
	base := location
	if urlBase != "" {
		base = urlBase
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", false
	}
	reference, err := url.Parse(controlURL)
	if err != nil {
		return "", false
	}
	resolved := baseURL.ResolveReference(reference).String()
	if !AllowedLocation(resolved, from) {
		return "", false
	}
	return resolved, true
}

func FetchDescription(ctx context.Context, client *http.Client, location string) (Description, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, location, nil)
	if err != nil {
		return Description{}, err
	}
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return Description{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Description{}, fmt.Errorf("description request returned HTTP %d", response.StatusCode)
	}
	description, err := ParseDescription(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return Description{}, err
	}
	description.Location = location
	return description, nil
}

func IsGateway(r Response) bool {
	return strings.Contains(strings.ToLower(r.ST), "internetgatewaydevice") ||
		strings.Contains(strings.ToLower(r.USN), "internetgatewaydevice")
}
