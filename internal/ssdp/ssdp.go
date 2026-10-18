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
	"strings"
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

func Discover(ctx context.Context, wait time.Duration) ([]Response, error) {
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("listen for SSDP: %w", err)
	}
	defer conn.Close()

	target, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return nil, fmt.Errorf("resolve SSDP multicast address: %w", err)
	}
	message := []byte("M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"MAN: \"ssdp:discover\"\r\n" +
		"MX: 2\r\n" +
		"ST: ssdp:all\r\n\r\n")
	for range 2 {
		if _, err := conn.WriteTo(message, target); err != nil {
			return nil, fmt.Errorf("send SSDP search: %w", err)
		}
	}

	deadline := time.Now().Add(wait)
	if wait <= 0 {
		deadline = time.Now()
	}
	_ = conn.SetReadDeadline(deadline)
	stopCancel := context.AfterFunc(ctx, func() {
		_ = conn.SetReadDeadline(time.Now())
	})
	defer stopCancel()

	responses := make([]Response, 0)
	seen := make(map[string]struct{})
	buffer := make([]byte, 2048)
	for {
		n, sender, err := conn.ReadFrom(buffer)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
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
		if !ip.Is4() || !ip.IsPrivate() {
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
	return responses, nil
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
