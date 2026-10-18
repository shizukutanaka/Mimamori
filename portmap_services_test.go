package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func TestFetchPortMappingsAcrossWANServices(t *testing.T) {
	gateway := netip.MustParseAddr("192.168.1.1")
	ipService := ssdp.Service{
		ServiceType: "urn:schemas-upnp-org:service:WANIPConnection:1",
		ControlURL:  "/ctl/IP",
	}
	pppService := ssdp.Service{
		ServiceType: "urn:schemas-upnp-org:service:WANPPPConnection:1",
		ControlURL:  "/ctl/PPP",
	}

	tests := []struct {
		name             string
		ipControlURL     string
		pppControlURL    string
		ipIndexZero      portMapHTTPResponse
		wantMappings     int
		wantErr          string
		wantRan          bool
		wantHTTPRequests int
	}{
		{
			name:          "empty IP list and populated PPP list",
			ipControlURL:  ipService.ControlURL,
			pppControlURL: pppService.ControlURL,
			ipIndexZero: portMapHTTPResponse{
				status: http.StatusInternalServerError,
				body:   wanServiceErrorXML(713),
			},
			wantMappings:     1,
			wantRan:          true,
			wantHTTPRequests: 3,
		},
		{
			name:          "IP query error retains PPP mappings",
			ipControlURL:  ipService.ControlURL,
			pppControlURL: pppService.ControlURL,
			ipIndexZero: portMapHTTPResponse{
				status: http.StatusInternalServerError,
			},
			wantMappings:     1,
			wantErr:          "WANIPConnection",
			wantRan:          true,
			wantHTTPRequests: 3,
		},
		{
			name:          "rejected PPP URL makes query incomplete",
			ipControlURL:  ipService.ControlURL,
			pppControlURL: "http://192.168.1.2/ctl/PPP",
			ipIndexZero: portMapHTTPResponse{
				status: http.StatusInternalServerError,
				body:   wanServiceErrorXML(713),
			},
			wantErr:          "WANPPPConnection",
			wantRan:          true,
			wantHTTPRequests: 1,
		},
		{
			name:             "external control URLs are rejected",
			ipControlURL:     "http://203.0.113.9/ctl",
			pppControlURL:    "http://203.0.113.9/ctl",
			wantErr:          "router advertises UPnP but its port forwarding list cannot be queried",
			wantHTTPRequests: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			description := ssdp.Description{
				Device: ssdp.Device{
					DeviceType: "urn:schemas-upnp-org:device:InternetGatewayDevice:1",
					Services: []ssdp.Service{
						{ServiceType: ipService.ServiceType, ControlURL: test.ipControlURL},
						{ServiceType: pppService.ServiceType, ControlURL: test.pppControlURL},
					},
				},
				Location: "http://192.168.1.1:5000/rootDesc.xml",
			}
			requests := 0
			client := &http.Client{Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
				requests++
				requestBody, err := io.ReadAll(request.Body)
				if err != nil {
					return nil, err
				}
				index := 0
				if bytes.Contains(requestBody, []byte("<NewPortMappingIndex>1</NewPortMappingIndex>")) {
					index = 1
				}

				var response portMapHTTPResponse
				switch request.URL.Path {
				case "/ctl/IP":
					response = test.ipIndexZero
				case "/ctl/PPP":
					switch index {
					case 0:
						response = portMapHTTPResponse{
							status: http.StatusOK,
							body:   wanServicePortMappingEntryXML(),
						}
					case 1:
						response = portMapHTTPResponse{
							status: http.StatusInternalServerError,
							body:   wanServiceErrorXML(713),
						}
					}
				default:
					return nil, fmt.Errorf("unexpected control URL path %q", request.URL.Path)
				}
				if response.status == 0 {
					return nil, fmt.Errorf("no response configured for %s index %d", request.URL.Path, index)
				}
				return &http.Response{
					StatusCode: response.status,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(response.body)),
					Request:    request,
				}, nil
			})}

			mappings, ran, err := fetchPortMappings(
				context.Background(),
				client,
				gateway,
				nil,
				map[netip.Addr]ssdp.Description{gateway: description},
			)
			if len(mappings) != test.wantMappings || ran != test.wantRan {
				t.Fatalf("fetchPortMappings() = %#v, %t, %v; want %d mappings, ran=%t", mappings, ran, err, test.wantMappings, test.wantRan)
			}
			if test.wantErr == "" && err != nil {
				t.Fatalf("fetchPortMappings() error = %v, want nil", err)
			}
			if test.wantErr != "" && (err == nil || !strings.Contains(err.Error(), test.wantErr)) {
				t.Fatalf("fetchPortMappings() error = %v, want text containing %q", err, test.wantErr)
			}
			if requests != test.wantHTTPRequests {
				t.Fatalf("HTTP requests = %d, want %d", requests, test.wantHTTPRequests)
			}
			if test.wantMappings == 1 && mappings[0] != (ssdp.PortMapping{
				ExternalPort:   27017,
				Protocol:       "TCP",
				InternalPort:   27017,
				InternalClient: "192.168.1.30",
				Enabled:        true,
			}) {
				t.Fatalf("mapping = %#v, want exposed database mapping", mappings[0])
			}
		})
	}
}

type portMapHTTPResponse struct {
	status int
	body   string
}

func wanServicePortMappingEntryXML() string {
	return `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>
<u:GetGenericPortMappingEntryResponse xmlns:u="urn:schemas-upnp-org:service:WANPPPConnection:1">
<NewRemoteHost></NewRemoteHost><NewExternalPort>27017</NewExternalPort><NewProtocol>TCP</NewProtocol>
<NewInternalPort>27017</NewInternalPort><NewInternalClient>192.168.1.30</NewInternalClient>
<NewEnabled>1</NewEnabled><NewPortMappingDescription></NewPortMappingDescription>
</u:GetGenericPortMappingEntryResponse></s:Body></s:Envelope>`
}

func wanServiceErrorXML(code int) string {
	return fmt.Sprintf(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><s:Fault><detail>
<UPnPError xmlns="urn:schemas-upnp-org:control-1-0"><errorCode>%d</errorCode></UPnPError>
</detail></s:Fault></s:Body></s:Envelope>`, code)
}
