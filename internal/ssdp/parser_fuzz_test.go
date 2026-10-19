package ssdp

import (
	"bytes"
	"encoding/xml"
	"net/netip"
	"strings"
	"testing"
)

func FuzzParseResponse(f *testing.F) {
	seeds := []string{
		"HTTP/1.1 200 OK\r\nST: upnp:rootdevice\r\nUSN: uuid:abc::upnp:rootdevice\r\nSERVER: router\r\nLOCATION: http://192.168.1.20/device.xml\r\n\r\n",
		"HTTP/1.1 200 OK\r\nnt: ssdp:all\r\nusn: uuid:def\r\n\r\n",
		"NOTIFY * HTTP/1.1\r\nST: ssdp:all\r\n\r\n",
		"HTTP/1.1 404 Not Found\r\nST: ssdp:all\r\n\r\n",
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	}
	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	from := netip.MustParseAddr("192.168.1.20")
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = ParseResponse(data, from)
	})
}

func FuzzParseDescription(f *testing.F) {
	seeds := []string{
		`<root><specVersion><major>1</major><minor>0</minor></specVersion>
<URLBase>http://192.168.1.20:54321/</URLBase><device>
<deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
<friendlyName>Home Router</friendlyName><manufacturer>Example</manufacturer><modelName>Router X</modelName>
</device></root>`,
		`<root><device><deviceList><device><deviceType>urn:schemas-upnp-org:device:WANDevice:1</deviceType>
<deviceList><device><deviceType>urn:schemas-upnp-org:device:WANConnectionDevice:1</deviceType>
<serviceList><service><serviceType>urn:schemas-upnp-org:service:WANIPConnection:1</serviceType><controlURL>/ctl/IPConn</controlURL></service></serviceList>
</device></deviceList></device></deviceList>
</device></root>`,
		`<root><device><deviceList><device><serviceList>
<service><serviceType>urn:schemas-upnp-org:service:WANPPPConnection:1</serviceType><controlURL>/ctl/PPPConn</controlURL></service>
</serviceList></device></deviceList></device></root>`,
		"",
		strings.Repeat("x", 1<<16),
		string([]byte{0xff, 0xfe, 0x80}),
	}
	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = ParseDescription(bytes.NewReader(data))
	})
}

func FuzzParsePortMappingSOAP(f *testing.F) {
	seeds := [][]byte{
		[]byte(portMappingEntryXML(1)),
		[]byte(portMappingEntryXML(2)),
		[]byte(portMappingErrorXML(713)),
		[]byte(portMappingErrorXML(501)),
		[]byte(portMappingErrorXML(402)),
		{},
		bytes.Repeat([]byte("x"), 1<<16),
		{0xff, 0xfe, 0x80},
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		var response portMappingResponse
		_ = xml.Unmarshal(data, &response)
		_, _ = upnpErrorCode(data)
	})
}

func FuzzControlURL(f *testing.F) {
	seeds := [][3]string{
		{"http://192.168.1.20:54321/device/description.xml", "", "../control/WANIPConn"},
		{"http://192.168.1.20:54321/device.xml", "http://192.168.1.20:1234/base/", "WANIPConn"},
		{"http://192.168.1.20/device.xml", "", "http://192.168.1.20:1234/control"},
		{"http://192.168.1.20/device.xml", "", "http://192.168.1.21/control"},
		{"http://192.168.1.20/device.xml", "", "http://8.8.8.8/control"},
		{"http://192.168.1.20/device.xml", "", "https://192.168.1.20/control"},
		{"", "", ""},
		{strings.Repeat("x", 1<<16), "", ""},
		{string([]byte{0xff, 0xfe}), "", ""},
	}
	for _, seed := range seeds {
		f.Add(seed[0], seed[1], seed[2])
	}

	from := netip.MustParseAddr("192.168.1.20")
	f.Fuzz(func(t *testing.T, location, urlBase, controlURL string) {
		_, _ = ControlURL(location, urlBase, controlURL, from)
	})
}
