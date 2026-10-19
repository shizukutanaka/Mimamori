package ssdp

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestDiscoverResendsAfterDelay(t *testing.T) {
	target, requests, _ := startSSDPTestResponder(t, func(request int, _ []byte) bool {
		return request == 2
	})

	responses, err := discover(
		context.Background(),
		200*time.Millisecond,
		target,
		nil,
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err != nil {
		t.Fatal(err)
	}
	if requests.Load() < 2 {
		t.Fatalf("received %d searches, want retransmission", requests.Load())
	}
	if len(responses) != 1 {
		t.Fatalf("Discover() returned %d responses, want 1", len(responses))
	}
}

func TestDiscoverUnicastSucceedsWhenMulticastSendFails(t *testing.T) {
	target, requests, _ := startSSDPTestResponder(t, func(int, []byte) bool { return true })
	unreachableMulticast := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}

	responses, err := discover(
		context.Background(),
		200*time.Millisecond,
		unreachableMulticast,
		[]*net.UDPAddr{target},
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err != nil {
		t.Fatal(err)
	}
	if requests.Load() < 2 {
		t.Fatalf("received %d unicast searches, want both message forms", requests.Load())
	}
	if len(responses) != 1 {
		t.Fatalf("Discover() returned %d responses, want 1", len(responses))
	}
}

func TestDiscoverReturnsErrorWhenAllSendsFail(t *testing.T) {
	unreachable := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}

	_, err := discover(
		context.Background(),
		200*time.Millisecond,
		unreachable,
		[]*net.UDPAddr{unreachable},
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err == nil {
		t.Fatal("Discover() error = nil, want all sends failed error")
	}
}

func TestDiscoverDeduplicatesRepliesAcrossRounds(t *testing.T) {
	target, requests, _ := startSSDPTestResponder(t, func(int, []byte) bool { return true })

	responses, err := discover(
		context.Background(),
		100*time.Millisecond,
		target,
		nil,
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err != nil {
		t.Fatal(err)
	}
	if requests.Load() < 2 {
		t.Fatalf("received %d searches, want both rounds", requests.Load())
	}
	if len(responses) != 1 {
		t.Fatalf("Discover() returned %d duplicate responses, want 1", len(responses))
	}
}

func TestDiscoverUnicastSearchMessage(t *testing.T) {
	target, _, received := startSSDPTestResponder(t, func(int, []byte) bool { return true })

	_, err := discover(
		context.Background(),
		40*time.Millisecond,
		nil,
		[]*net.UDPAddr{target},
		netip.Addr.IsLoopback,
		100*time.Millisecond,
	)
	if err != nil {
		t.Fatal(err)
	}
	readRequest := func() string {
		t.Helper()
		select {
		case request := <-received:
			return string(request)
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for unicast SSDP search")
			return ""
		}
	}
	request := readRequest()
	if !strings.Contains(request, "HOST: 127.0.0.1:1900\r\n") {
		t.Fatalf("unicast search HOST header = %q, want target IP with port 1900", request)
	}
	if !strings.Contains(request, "MAN: \"ssdp:discover\"\r\n") || !strings.Contains(request, "ST: ssdp:all\r\n") {
		t.Fatalf("unicast search is missing MAN or ST: %q", request)
	}
	if strings.Contains(request, "\r\nMX:") {
		t.Fatalf("unicast search unexpectedly includes MX: %q", request)
	}

	request = readRequest()
	if !strings.Contains(request, "HOST: 239.255.255.250:1900\r\n") {
		t.Fatalf("multicast-form unicast search HOST header = %q, want multicast HOST", request)
	}
	if !strings.Contains(request, "MAN: \"ssdp:discover\"\r\n") ||
		!strings.Contains(request, "MX: 2\r\n") ||
		!strings.Contains(request, "ST: ssdp:all\r\n") {
		t.Fatalf("multicast-form unicast search is missing MAN, MX, or ST: %q", request)
	}
}

func TestDiscoverDelayedUnicast(t *testing.T) {
	target, _, _ := startSSDPTestResponder(t, func(int, []byte) bool { return true })
	unicast := make(chan *net.UDPAddr)
	go func() {
		time.Sleep(150 * time.Millisecond)
		unicast <- target
		close(unicast)
	}()

	responses, err := discoverTargets(
		context.Background(),
		400*time.Millisecond,
		nil,
		unicast,
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(responses) != 1 {
		t.Fatalf("discoverTargets() returned %d responses, want 1", len(responses))
	}
}

func TestDiscoverDelayedUnicastAfterDeadline(t *testing.T) {
	target, requests, _ := startSSDPTestResponder(t, func(int, []byte) bool { return true })
	unicast := make(chan *net.UDPAddr, 1)
	go func() {
		time.Sleep(80 * time.Millisecond)
		unicast <- target
		close(unicast)
	}()

	responses, err := discoverTargets(
		context.Background(),
		50*time.Millisecond,
		nil,
		unicast,
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(responses) != 1 {
		t.Fatalf("discoverTargets() returned %d responses, want 1", len(responses))
	}
	if requests.Load() == 0 {
		t.Fatal("no unicast search was sent")
	}
}

func TestDiscoverNilUnicast(t *testing.T) {
	target, _, _ := startSSDPTestResponder(t, func(int, []byte) bool { return true })

	responses, err := discoverTargets(
		context.Background(),
		60*time.Millisecond,
		target,
		nil,
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(responses) != 1 {
		t.Fatalf("discoverTargets() returned %d responses, want 1", len(responses))
	}
}

func TestDiscoverClosedUnicastReturnsAtDeadline(t *testing.T) {
	wait := 80 * time.Millisecond
	unicast := make(chan *net.UDPAddr)
	close(unicast)
	start := time.Now()
	_, err := discoverTargets(
		context.Background(),
		wait,
		nil,
		unicast,
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if err == nil {
		t.Fatal("discoverTargets() error = nil, want no-send error")
	}
	if elapsed := time.Since(start); elapsed < wait-10*time.Millisecond || elapsed > wait+500*time.Millisecond {
		t.Fatalf("discoverTargets() took %s, want return near %s", elapsed, wait)
	}
}

func TestDiscoverNeverClosedUnicastReturnsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	unicast := make(chan *net.UDPAddr)
	returned := make(chan error, 1)
	started := make(chan struct{})
	go func() {
		close(started)
		_, err := discoverTargets(
			ctx,
			40*time.Millisecond,
			nil,
			unicast,
			netip.Addr.IsLoopback,
			20*time.Millisecond,
		)
		returned <- err
	}()
	<-started

	time.Sleep(80 * time.Millisecond)
	select {
	case err := <-returned:
		t.Fatalf("discoverTargets() returned before cancellation: %v", err)
	default:
	}

	cancelledAt := time.Now()
	cancel()
	select {
	case err := <-returned:
		if time.Since(cancelledAt) > 500*time.Millisecond {
			t.Fatal("discoverTargets() returned too slowly after cancellation")
		}
		if err == nil {
			t.Fatal("discoverTargets() error = nil, want no-send error")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("discoverTargets() did not return after cancellation")
	}
}

func TestDiscoverCancellationAfterRetrySend(t *testing.T) {
	target, _, _ := startSSDPTestResponder(t, func(int, []byte) bool { return true })
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	previous := afterRetrySend
	t.Cleanup(func() { afterRetrySend = previous })
	afterRetrySend = cancel

	start := time.Now()
	_, _ = discoverTargets(
		ctx,
		time.Second,
		target,
		nil,
		netip.Addr.IsLoopback,
		20*time.Millisecond,
	)
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("discoverTargets() took %s after cancellation, want under 500ms", elapsed)
	}
}

func startSSDPTestResponder(
	t *testing.T,
	respond func(int, []byte) bool,
) (*net.UDPAddr, *atomic.Int32, <-chan []byte) {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	requests := &atomic.Int32{}
	received := make(chan []byte, 8)
	go func() {
		buffer := make([]byte, 2048)
		for {
			n, sender, err := conn.ReadFromUDP(buffer)
			if err != nil {
				return
			}
			request := bytes.Clone(buffer[:n])
			requestNumber := int(requests.Add(1))
			select {
			case received <- request:
			default:
			}
			if respond(requestNumber, request) {
				response := []byte("HTTP/1.1 200 OK\r\n" +
					"ST: upnp:rootdevice\r\n" +
					"USN: uuid:test::upnp:rootdevice\r\n\r\n")
				_, _ = conn.WriteToUDP(response, sender)
			}
		}
	}()
	return conn.LocalAddr().(*net.UDPAddr), requests, received
}

func TestParseResponse(t *testing.T) {
	from := netip.MustParseAddr("192.168.1.20")
	tests := []struct {
		name string
		data string
		ok   bool
		st   string
	}{
		{
			name: "normal",
			data: "HTTP/1.1 200 OK\r\nST: upnp:rootdevice\r\nUSN: uuid:abc::upnp:rootdevice\r\nSERVER: router\r\nLOCATION: http://192.168.1.20/device.xml\r\n\r\n",
			ok:   true,
			st:   "upnp:rootdevice",
		},
		{
			name: "lowercase headers and NT fallback",
			data: "HTTP/1.1 200 OK\r\nnt: ssdp:all\r\nusn: uuid:def\r\n\r\n",
			ok:   true,
			st:   "ssdp:all",
		},
		{name: "NOTIFY rejected", data: "NOTIFY * HTTP/1.1\r\nST: ssdp:all\r\n\r\n"},
		{name: "non-200 rejected", data: "HTTP/1.1 404 Not Found\r\nST: ssdp:all\r\n\r\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := ParseResponse([]byte(test.data), from)
			if ok != test.ok {
				t.Fatalf("ParseResponse() ok = %v, want %v", ok, test.ok)
			}
			if ok && (got.IP != from || got.ST != test.st) {
				t.Fatalf("ParseResponse() = %#v", got)
			}
		})
	}
}

func TestParseDescription(t *testing.T) {
	xmlDescription := `<root><specVersion><major>1</major><minor>0</minor></specVersion>
<URLBase>http://192.168.1.20:54321/</URLBase><device>
<deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
<friendlyName>Home Router</friendlyName><manufacturer>Example</manufacturer><modelName>Router X</modelName>
</device></root>`
	got, err := ParseDescription(strings.NewReader(xmlDescription))
	if err != nil {
		t.Fatal(err)
	}
	if got.Device.FriendlyName != "Home Router" || got.Device.Manufacturer != "Example" ||
		got.Device.ModelName != "Router X" || got.Device.DeviceType == "" {
		t.Fatalf("ParseDescription() = %#v", got)
	}
	if got.URLBase != "http://192.168.1.20:54321/" {
		t.Fatalf("URLBase = %q", got.URLBase)
	}
}

func TestWANConnectionsNestedDevice(t *testing.T) {
	xmlDescription := `<root><device>
<deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
<deviceList><device><deviceType>urn:schemas-upnp-org:device:WANDevice:1</deviceType>
<deviceList><device><deviceType>urn:schemas-upnp-org:device:WANConnectionDevice:1</deviceType>
<serviceList><service><serviceType>urn:schemas-upnp-org:service:WANIPConnection:1</serviceType><controlURL>/ctl/IPConn</controlURL></service></serviceList>
</device></deviceList></device></deviceList>
</device></root>`
	description, err := ParseDescription(strings.NewReader(xmlDescription))
	if err != nil {
		t.Fatal(err)
	}
	services := description.Device.WANConnections()
	if len(services) != 1 || services[0].ServiceType != "urn:schemas-upnp-org:service:WANIPConnection:1" ||
		services[0].ControlURL != "/ctl/IPConn" {
		t.Fatalf("WANConnections() = %#v", services)
	}
}

func TestWANConnectionsTrimsWhitespace(t *testing.T) {
	xmlDescription := `<root><device>
<deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
<deviceList><device><deviceType>urn:schemas-upnp-org:device:WANDevice:1</deviceType>
<deviceList><device><deviceType>urn:schemas-upnp-org:device:WANConnectionDevice:1</deviceType>
<serviceList><service>
<serviceType>
  urn:schemas-upnp-org:service:WANIPConnection:1
</serviceType>
<controlURL>
  /ctl/IPConn
</controlURL>
</service></serviceList>
</device></deviceList></device></deviceList>
</device></root>`
	description, err := ParseDescription(strings.NewReader(xmlDescription))
	if err != nil {
		t.Fatal(err)
	}
	want := Service{ServiceType: "urn:schemas-upnp-org:service:WANIPConnection:1", ControlURL: "/ctl/IPConn"}
	services := description.Device.WANConnections()
	if len(services) != 1 || services[0] != want {
		t.Fatalf("WANConnections() = %#v, want [%#v]", services, want)
	}
}

func TestWANConnectionsPPP(t *testing.T) {
	description, err := ParseDescription(strings.NewReader(`<root><device><deviceList><device><serviceList>
<service><serviceType>urn:schemas-upnp-org:service:WANPPPConnection:1</serviceType><controlURL>/ctl/PPPConn</controlURL></service>
</serviceList></device></deviceList></device></root>`))
	if err != nil {
		t.Fatal(err)
	}
	services := description.Device.WANConnections()
	if len(services) != 1 || services[0].ServiceType != "urn:schemas-upnp-org:service:WANPPPConnection:1" {
		t.Fatalf("WANConnections() = %#v", services)
	}
}

func TestWANConnectionsReturnsUniqueServicesInDocumentOrder(t *testing.T) {
	description, err := ParseDescription(strings.NewReader(`<root><device><deviceList>
<device><deviceType>urn:schemas-upnp-org:device:WANConnectionDevice:1</deviceType><serviceList>
<service><serviceType>urn:schemas-upnp-org:service:WANIPConnection:1</serviceType><controlURL>/ctl/IP</controlURL></service>
<service><serviceType>urn:schemas-upnp-org:service:WANPPPConnection:1</serviceType><controlURL>/ctl/PPP</controlURL></service>
</serviceList></device>
<device><deviceType>urn:schemas-upnp-org:device:WANConnectionDevice:1</deviceType><serviceList>
<service><serviceType>urn:schemas-upnp-org:service:WANPPPConnection:1</serviceType><controlURL>/ctl/PPP</controlURL></service>
</serviceList></device>
</deviceList></device></root>`))
	if err != nil {
		t.Fatal(err)
	}
	want := []Service{
		{ServiceType: "urn:schemas-upnp-org:service:WANIPConnection:1", ControlURL: "/ctl/IP"},
		{ServiceType: "urn:schemas-upnp-org:service:WANPPPConnection:1", ControlURL: "/ctl/PPP"},
	}
	services := description.Device.WANConnections()
	if len(services) != len(want) {
		t.Fatalf("WANConnections() = %#v, want %#v", services, want)
	}
	for index := range want {
		if services[index] != want[index] {
			t.Fatalf("WANConnections() = %#v, want %#v", services, want)
		}
	}
}

func TestDeviceIsIGD(t *testing.T) {
	tests := []struct {
		name   string
		device Device
		want   bool
	}{
		{
			name:   "root IGD",
			device: Device{DeviceType: "urn:schemas-upnp-org:device:InternetGatewayDevice:1"},
			want:   true,
		},
		{
			name: "nested IGD",
			device: Device{
				DeviceType: "urn:schemas-upnp-org:device:Basic:1",
				Devices: []Device{{
					DeviceType: "urn:schemas-upnp-org:device:INTERNETGATEWAYDEVICE:1",
				}},
			},
			want: true,
		},
		{name: "not an IGD", device: Device{DeviceType: "urn:schemas-upnp-org:device:Basic:1"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.device.IsIGD(); got != test.want {
				t.Fatalf("IsIGD() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestAllowedLocation(t *testing.T) {
	from := netip.MustParseAddr("192.168.1.20")
	tests := []struct {
		location string
		want     bool
	}{
		{"http://192.168.1.20/description.xml", true},
		{"https://192.168.1.20/description.xml", false},
		{"http://router.local/description.xml", false},
		{"http://192.168.1.21/description.xml", false},
		{"http://8.8.8.8/description.xml", false},
		{"http://user@192.168.1.20/description.xml", false},
	}
	for _, test := range tests {
		if got := AllowedLocation(test.location, from); got != test.want {
			t.Errorf("AllowedLocation(%q) = %v, want %v", test.location, got, test.want)
		}
	}
}

func TestControlURL(t *testing.T) {
	from := netip.MustParseAddr("192.168.1.20")
	tests := []struct {
		name        string
		location    string
		urlBase     string
		controlURL  string
		want        string
		wantAllowed bool
	}{
		{
			name:        "relative to location",
			location:    "http://192.168.1.20:54321/device/description.xml",
			controlURL:  "../control/WANIPConn",
			want:        "http://192.168.1.20:54321/control/WANIPConn",
			wantAllowed: true,
		},
		{
			name:        "relative to URLBase",
			location:    "http://192.168.1.20:54321/device.xml",
			urlBase:     "http://192.168.1.20:1234/base/",
			controlURL:  "WANIPConn",
			want:        "http://192.168.1.20:1234/base/WANIPConn",
			wantAllowed: true,
		},
		{
			name:        "absolute",
			location:    "http://192.168.1.20/device.xml",
			controlURL:  "http://192.168.1.20:1234/control",
			want:        "http://192.168.1.20:1234/control",
			wantAllowed: true,
		},
		{name: "different host", location: "http://192.168.1.20/device.xml", controlURL: "http://192.168.1.21/control"},
		{name: "public host", location: "http://192.168.1.20/device.xml", controlURL: "http://8.8.8.8/control"},
		{name: "https", location: "http://192.168.1.20/device.xml", controlURL: "https://192.168.1.20/control"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := ControlURL(test.location, test.urlBase, test.controlURL, from)
			if ok != test.wantAllowed || (ok && got != test.want) {
				t.Fatalf("ControlURL() = %q, %v; want %q, %v", got, ok, test.want, test.wantAllowed)
			}
		})
	}
}

func TestControlURLTrimsWhitespace(t *testing.T) {
	from := netip.MustParseAddr("192.168.1.1")
	location := "http://192.168.1.1:5000/rootDesc.xml"
	tests := []struct {
		name        string
		urlBase     string
		controlURL  string
		want        string
		wantAllowed bool
	}{
		{
			name:        "padded control URL",
			controlURL:  "\n  /ctl/IPConn\n",
			want:        "http://192.168.1.1:5000/ctl/IPConn",
			wantAllowed: true,
		},
		{
			name:        "padded URLBase",
			urlBase:     "\n http://192.168.1.1:5000/ \n",
			controlURL:  "ctl/IPConn",
			want:        "http://192.168.1.1:5000/ctl/IPConn",
			wantAllowed: true,
		},
		{
			name:       "padded off-router control URL",
			controlURL: "\n http://8.8.8.8/ctl\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := ControlURL(location, test.urlBase, test.controlURL, from)
			if ok != test.wantAllowed || (ok && got != test.want) {
				t.Fatalf("ControlURL() = %q, %v; want %q, %v", got, ok, test.want, test.wantAllowed)
			}
		})
	}
}

func TestFetchDescription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		_, _ = w.Write([]byte(`<root><device><friendlyName>Test device</friendlyName></device></root>`))
	}))
	defer server.Close()

	description, err := FetchDescription(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if description.Device.FriendlyName != "Test device" || description.Location != server.URL {
		t.Fatalf("FetchDescription() = %#v", description)
	}
}

func TestFetchDescriptionRequiresStatusOK(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	if _, err := FetchDescription(context.Background(), server.Client(), server.URL); err == nil {
		t.Fatal("expected non-200 response error")
	}
}
