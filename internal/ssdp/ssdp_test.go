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
	description := `<root><specVersion><major>1</major><minor>0</minor></specVersion><device>
<deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
<friendlyName>Home Router</friendlyName><manufacturer>Example</manufacturer><modelName>Router X</modelName>
</device></root>`
	got, err := ParseDescription(strings.NewReader(description))
	if err != nil {
		t.Fatal(err)
	}
	if got.FriendlyName != "Home Router" || got.Manufacturer != "Example" || got.ModelName != "Router X" || got.DeviceType == "" {
		t.Fatalf("ParseDescription() = %#v", got)
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

func TestFetchDescription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		_, _ = w.Write([]byte(`<root><device><friendlyName>Test device</friendlyName></device></root>`))
	}))
	defer server.Close()

	device, err := FetchDescription(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if device.FriendlyName != "Test device" {
		t.Fatalf("FetchDescription() = %#v", device)
	}
}

func TestFetchDescriptionRequiresStatusOK(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	if _, err := FetchDescription(context.Background(), server.Client(), server.URL); err == nil {
		t.Fatal("expected non-200 response error")
	}
}
