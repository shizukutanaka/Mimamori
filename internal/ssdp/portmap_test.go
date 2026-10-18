package ssdp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const testWANIPService = "urn:schemas-upnp-org:service:WANIPConnection:1"

func TestPortMappingsReadsEntriesUntilInvalidIndex(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		assertPortMappingRequest(t, r, requests-1)
		if requests <= 2 {
			_, _ = io.WriteString(w, portMappingEntryXML(requests))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, portMappingErrorXML(713))
	}))
	defer server.Close()

	mappings, err := PortMappings(context.Background(), server.Client(), server.URL, testWANIPService)
	if err != nil {
		t.Fatal(err)
	}
	if len(mappings) != 2 {
		t.Fatalf("PortMappings() returned %d mappings, want 2: %#v", len(mappings), mappings)
	}
	if mappings[0] != (PortMapping{
		RemoteHost: "", ExternalPort: 80, Protocol: "TCP", InternalPort: 8080,
		InternalClient: "192.168.1.30", Enabled: true, Description: "Web camera",
	}) {
		t.Fatalf("first mapping = %#v", mappings[0])
	}
	if mappings[1].ExternalPort != 81 || mappings[1].Enabled {
		t.Fatalf("second mapping = %#v", mappings[1])
	}
	if requests != 3 {
		t.Fatalf("request count = %d, want 3", requests)
	}
}

func TestPortMappingsReturnsPartialEntriesAndError(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		assertPortMappingRequest(t, r, requests-1)
		if requests == 1 {
			_, _ = io.WriteString(w, portMappingEntryXML(1))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, portMappingErrorXML(501))
	}))
	defer server.Close()

	mappings, err := PortMappings(context.Background(), server.Client(), server.URL, testWANIPService)
	if err == nil {
		t.Fatal("PortMappings() returned no error for an unexpected UPnP error")
	}
	if len(mappings) != 1 {
		t.Fatalf("PortMappings() returned %d mappings, want 1", len(mappings))
	}
}

func TestPortMappingsEmptyList(t *testing.T) {
	for _, code := range []int{713, 714} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				assertPortMappingRequest(t, r, 0)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, portMappingErrorXML(code))
			}))
			defer server.Close()

			mappings, err := PortMappings(context.Background(), server.Client(), server.URL, testWANIPService)
			if err != nil || len(mappings) != 0 {
				t.Fatalf("PortMappings() = %#v, %v; want empty list without error", mappings, err)
			}
			if requests != 1 {
				t.Fatalf("request count = %d, want 1", requests)
			}
		})
	}
}

func TestPortMappingsErrorCode402(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		assertPortMappingRequest(t, r, 0)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, portMappingErrorXML(402))
	}))
	defer server.Close()

	mappings, err := PortMappings(context.Background(), server.Client(), server.URL, testWANIPService)
	if err == nil || len(mappings) != 0 {
		t.Fatalf("PortMappings() = %#v, %v; want error with no mappings", mappings, err)
	}
	if requests != 1 {
		t.Fatalf("request count = %d, want 1", requests)
	}
}

func assertPortMappingRequest(t *testing.T, r *http.Request, wantIndex int) {
	t.Helper()
	if r.Method != http.MethodPost {
		t.Errorf("method = %s, want POST", r.Method)
	}
	if got, want := r.Header.Get("SOAPAction"), `"`+testWANIPService+`#GetGenericPortMappingEntry"`; got != want {
		t.Errorf("SOAPAction = %q, want %q", got, want)
	}
	if got, want := r.Header.Get("Content-Type"), `text/xml; charset="utf-8"`; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("read request body: %v", err)
		return
	}
	if !strings.Contains(string(body), fmt.Sprintf("<NewPortMappingIndex>%d</NewPortMappingIndex>", wantIndex)) {
		t.Errorf("request body has no index %d: %s", wantIndex, body)
	}
}

func portMappingEntryXML(externalPort int) string {
	return fmt.Sprintf(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>
<u:GetGenericPortMappingEntryResponse xmlns:u="%s">
<NewRemoteHost></NewRemoteHost><NewExternalPort>%d</NewExternalPort><NewProtocol>TCP</NewProtocol>
<NewInternalPort>8080</NewInternalPort><NewInternalClient>192.168.1.30</NewInternalClient>
<NewEnabled>%s</NewEnabled><NewPortMappingDescription>Web camera</NewPortMappingDescription>
</u:GetGenericPortMappingEntryResponse></s:Body></s:Envelope>`,
		testWANIPService, externalPort+79, map[bool]string{true: "1", false: "false"}[externalPort == 1])
}

func portMappingErrorXML(code int) string {
	return fmt.Sprintf(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><s:Fault><detail>
<UPnPError xmlns="urn:schemas-upnp-org:control-1-0"><errorCode>%d</errorCode></UPnPError>
</detail></s:Fault></s:Body></s:Envelope>`, code)
}
