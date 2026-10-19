package main

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"testing"
	"time"

	"github.com/shizukutanaka/Mimamori/internal/probe"
)

func TestMarkUnauthenticatedZigbee2MQTT(t *testing.T) {
	message := []byte(`{"topic":"bridge/state","payload":{"state":"online"}}`)
	frame := append([]byte{0x81, byte(len(message))}, message...)
	tests := []struct {
		name     string
		response zigbee2MQTTTestResponse
		want     bool
	}{
		{
			name:     "unauthenticated retained bridge message is kept",
			response: zigbee2MQTTTestResponse{status: http.StatusSwitchingProtocols, frame: frame},
			want:     true,
		},
		{
			name:     "authentication close frame is dropped",
			response: zigbee2MQTTTestResponse{status: http.StatusSwitchingProtocols, frame: []byte{0x88, 2, 0x11, 0x31}},
		},
		{
			name:     "HTTP page is dropped",
			response: zigbee2MQTTTestResponse{status: http.StatusOK, body: "<html>not Zigbee2MQTT</html>"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			z2mOpen, requests := startZigbee2MQTTTestServer(t, test.response)
			otherPort := probe.Open{IP: z2mOpen.IP, Port: 5900}
			opened := []probe.Open{z2mOpen, otherPort}
			want := []probe.Open{otherPort}
			if test.want {
				want = opened
			}
			got := markUnauthenticatedZigbee2MQTT(context.Background(), opened)
			if len(got) != len(want) {
				t.Fatalf("markUnauthenticatedZigbee2MQTT() returned %v, want %v", got, want)
			}
			for index := range want {
				if got[index] != want[index] {
					t.Fatalf("markUnauthenticatedZigbee2MQTT() returned %v, want %v", got, want)
				}
			}
			select {
			case request := <-requests:
				if request != "GET /api HTTP/1.1" {
					t.Fatalf("request = %q, want GET /api HTTP/1.1", request)
				}
			case <-time.After(time.Second):
				t.Fatal("Zigbee2MQTT fake server received no request")
			}
		})
	}
}

type zigbee2MQTTTestResponse struct {
	status int
	body   string
	frame  []byte
}

func startZigbee2MQTTTestServer(t *testing.T, response zigbee2MQTTTestResponse) (probe.Open, <-chan string) {
	t.Helper()
	const address = "127.0.0.1"
	listener, err := net.Listen("tcp4", net.JoinHostPort(address, "8080"))
	if err != nil {
		t.Skipf("cannot listen on %s:8080: %v", address, err)
	}
	requests := make(chan string, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		connection, err := listener.Accept()
		if err != nil {
			requests <- "error: " + err.Error()
			return
		}
		defer connection.Close()
		_ = connection.SetDeadline(time.Now().Add(3 * time.Second))
		reader := bufio.NewReader(connection)
		request, err := http.ReadRequest(reader)
		if err != nil {
			requests <- "error: " + err.Error()
			return
		}
		requests <- fmt.Sprintf("%s %s %s", request.Method, request.RequestURI, request.Proto)
		if response.status == http.StatusSwitchingProtocols {
			hash := sha1.Sum([]byte(request.Header.Get("Sec-WebSocket-Key") + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
			accept := base64.StdEncoding.EncodeToString(hash[:])
			_, _ = io.WriteString(connection, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "+accept+"\r\n\r\n")
			_, _ = connection.Write(response.frame)
		} else {
			_, _ = fmt.Fprintf(connection, "HTTP/1.1 %d %s\r\nContent-Length: %d\r\n\r\n%s", response.status, http.StatusText(response.status), len(response.body), response.body)
		}
		_, _ = io.Copy(io.Discard, reader)
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		select {
		case <-done:
		case <-time.After(4 * time.Second):
			t.Error("Zigbee2MQTT fake server did not stop")
		}
	})
	return probe.Open{IP: netip.MustParseAddr(address), Port: 8080}, requests
}
