package ssdp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
)

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

func TestWANConnectionNestedDevice(t *testing.T) {
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
	service, ok := description.Device.WANConnection()
	if !ok || service.ServiceType != "urn:schemas-upnp-org:service:WANIPConnection:1" || service.ControlURL != "/ctl/IPConn" {
		t.Fatalf("WANConnection() = %#v, %v", service, ok)
	}
}

func TestWANConnectionPPP(t *testing.T) {
	description, err := ParseDescription(strings.NewReader(`<root><device><deviceList><device><serviceList>
<service><serviceType>urn:schemas-upnp-org:service:WANPPPConnection:1</serviceType><controlURL>/ctl/PPPConn</controlURL></service>
</serviceList></device></deviceList></device></root>`))
	if err != nil {
		t.Fatal(err)
	}
	service, ok := description.Device.WANConnection()
	if !ok || service.ServiceType != "urn:schemas-upnp-org:service:WANPPPConnection:1" {
		t.Fatalf("WANConnection() = %#v, %v", service, ok)
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
