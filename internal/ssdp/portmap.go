package ssdp

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PortMapping struct {
	RemoteHost     string
	ExternalPort   int
	Protocol       string
	InternalPort   int
	InternalClient string
	Enabled        bool
	Description    string
}

type portMappingResponse struct {
	RemoteHost     string
	ExternalPort   int
	Protocol       string
	InternalPort   int
	InternalClient string
	Enabled        string
	Description    string
}

func (mapping *portMappingResponse) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch value := token.(type) {
		case xml.StartElement:
			switch value.Name.Local {
			case "NewRemoteHost":
				err = decoder.DecodeElement(&mapping.RemoteHost, &value)
			case "NewExternalPort":
				err = decoder.DecodeElement(&mapping.ExternalPort, &value)
			case "NewProtocol":
				err = decoder.DecodeElement(&mapping.Protocol, &value)
			case "NewInternalPort":
				err = decoder.DecodeElement(&mapping.InternalPort, &value)
			case "NewInternalClient":
				err = decoder.DecodeElement(&mapping.InternalClient, &value)
			case "NewEnabled":
				err = decoder.DecodeElement(&mapping.Enabled, &value)
			case "NewPortMappingDescription":
				err = decoder.DecodeElement(&mapping.Description, &value)
			}
			if err != nil {
				return err
			}
		case xml.EndElement:
			if value.Name == start.Name {
				return nil
			}
		}
	}
}

func PortMappings(ctx context.Context, client *http.Client, controlURL, serviceType string) ([]PortMapping, error) {
	if strings.TrimSpace(serviceType) == "" || strings.ContainsAny(serviceType, "\r\n\"") {
		return nil, fmt.Errorf("invalid UPnP service type")
	}
	if client == nil {
		client = http.DefaultClient
	}

	mappings := make([]PortMapping, 0)
	for index := 0; index < 128; index++ {
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, controlURL, bytes.NewReader(portMappingRequest(serviceType, index)))
		if err != nil {
			return mappings, err
		}
		request.Header.Set("SOAPAction", `"`+serviceType+`#GetGenericPortMappingEntry"`)
		request.Header.Set("Content-Type", `text/xml; charset="utf-8"`)

		response, err := client.Do(request)
		if err != nil {
			return mappings, err
		}
		body, readErr := io.ReadAll(io.LimitReader(response.Body, 64*1024))
		closeErr := response.Body.Close()
		if readErr != nil {
			return mappings, fmt.Errorf("read port mapping response: %w", readErr)
		}
		if closeErr != nil {
			return mappings, fmt.Errorf("close port mapping response: %w", closeErr)
		}

		if response.StatusCode == http.StatusOK {
			var entry portMappingResponse
			if err := xml.Unmarshal(body, &entry); err != nil {
				return mappings, fmt.Errorf("parse port mapping entry %d: %w", index, err)
			}
			enabled := strings.TrimSpace(entry.Enabled)
			mappings = append(mappings, PortMapping{
				RemoteHost:     strings.TrimSpace(entry.RemoteHost),
				ExternalPort:   entry.ExternalPort,
				Protocol:       strings.TrimSpace(entry.Protocol),
				InternalPort:   entry.InternalPort,
				InternalClient: strings.TrimSpace(entry.InternalClient),
				Enabled:        enabled == "1" || strings.EqualFold(enabled, "true") || strings.EqualFold(enabled, "yes"),
				Description:    entry.Description,
			})
			continue
		}

		if response.StatusCode == http.StatusInternalServerError {
			switch code, ok := upnpErrorCode(body); {
			case ok && code == 713:
				return mappings, nil
			case ok && index == 0 && code == 714:
				return mappings, nil
			}
		}
		return mappings, fmt.Errorf("port mapping request at index %d returned HTTP %d", index, response.StatusCode)
	}
	return mappings, fmt.Errorf("router has at least 128 UPnP port mappings")
}

func portMappingRequest(serviceType string, index int) []byte {
	var body bytes.Buffer
	body.WriteString(`<?xml version="1.0" encoding="utf-8"?>` +
		`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" ` +
		`s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body>` +
		`<u:GetGenericPortMappingEntry xmlns:u="`)
	_ = xml.EscapeText(&body, []byte(serviceType))
	fmt.Fprintf(&body, `"><NewPortMappingIndex>%d</NewPortMappingIndex></u:GetGenericPortMappingEntry></s:Body></s:Envelope>`, index)
	return body.Bytes()
}

func upnpErrorCode(body []byte) (int, bool) {
	decoder := xml.NewDecoder(bytes.NewReader(body))
	for {
		token, err := decoder.Token()
		if err != nil {
			return 0, false
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "errorCode" {
			continue
		}
		var code int
		if err := decoder.DecodeElement(&code, &start); err != nil {
			return 0, false
		}
		return code, true
	}
}
