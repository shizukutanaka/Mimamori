package zigbee2mqtt

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"testing"
	"time"
)

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    bool
	}{
		{name: "2.x topic first", payload: `{"topic":"bridge/state","payload":{"state":"online"}}`, want: true},
		{name: "1.x payload first", payload: `{"payload":"online","topic":"bridge/state"}`, want: true},
		{name: "topic without bridge prefix", payload: `{"topic":"device/state","payload":{}}`},
		{name: "topic is not a string", payload: `{"topic":1,"payload":{}}`},
		{name: "missing payload", payload: `{"topic":"bridge/state"}`},
		{name: "non-JSON", payload: `not JSON`},
		{name: "JSON array", payload: `[{"topic":"bridge/state","payload":{}}]`},
		{name: "null", payload: `null`},
		{name: "payload may be null", payload: `{"topic":"bridge/state","payload":null}`, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseMessage([]byte(test.payload)); got != test.want {
				t.Fatalf("ParseMessage() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	validMessage := []byte(`{"topic":"bridge/state","payload":{"state":"online"}}`)
	reorderedMessage := []byte(`{"payload":"online","topic":"bridge/state"}`)
	longMessage := append([]byte(nil), validMessage...)
	longMessage = append(longMessage, bytes.Repeat([]byte(" "), 126-len(longMessage))...)
	veryLongMessage := append([]byte(nil), validMessage...)
	veryLongMessage = append(veryLongMessage, bytes.Repeat([]byte(" "), 1<<16-len(veryLongMessage))...)
	tooLongFrame := []byte{0x81, 0x7f}
	var extendedLength [8]byte
	binary.BigEndian.PutUint64(extendedLength[:], maxFrameBytes+1)
	tooLongFrame = append(tooLongFrame, extendedLength[:]...)

	tests := []struct {
		name     string
		response fakeResponse
		want     bool
	}{
		{
			name:     "unauthenticated 2.x retained bridge message",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 1, validMessage)},
			want:     true,
		},
		{
			name:     "unauthenticated 1.x retained bridge message",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 1, reorderedMessage)},
			want:     true,
		},
		{
			name:     "126-length-form text frame",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 1, longMessage)},
			want:     true,
		},
		{
			name:     "127-length-form text frame",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 1, veryLongMessage)},
			want:     true,
		},
		{
			name:     "authentication close frame",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 8, []byte{0x11, 0x31})},
		},
		{
			name:     "HTTP page",
			response: fakeResponse{status: http.StatusOK, body: "<html>not Zigbee2MQTT</html>"},
		},
		{
			name:     "wrong WebSocket accept",
			response: fakeResponse{status: http.StatusSwitchingProtocols, accept: "wrong", frame: unmaskedFrame(true, 1, validMessage)},
		},
		{
			name:     "masked server frame",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: maskedFrame(validMessage)},
		},
		{
			name:     "frame length exceeds limit",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: tooLongFrame},
		},
		{
			name:     "non-bridge topic",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 1, []byte(`{"topic":"device/state","payload":{}}`))},
		},
		{
			name:     "non-JSON message",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 1, []byte("not JSON"))},
		},
		{
			name:     "message without payload",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(true, 1, []byte(`{"topic":"bridge/state"}`))},
		},
		{
			name:     "non-final text frame",
			response: fakeResponse{status: http.StatusSwitchingProtocols, frame: unmaskedFrame(false, 1, validMessage)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port, result := startFakeWebSocketServer(t, test.response)
			got := Unauthenticated(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second)
			if got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			select {
			case serverResult := <-result:
				if serverResult.err != nil {
					t.Fatal(serverResult.err)
				}
				request := serverResult.request
				if request.method != http.MethodGet || request.requestURI != "/api" || request.proto != "HTTP/1.1" {
					t.Errorf("request line = %q %q %q, want GET /api HTTP/1.1", request.method, request.requestURI, request.proto)
				}
				if request.upgrade != "websocket" || request.connection != "Upgrade" || request.version != "13" {
					t.Errorf("WebSocket headers = Upgrade:%q Connection:%q Version:%q", request.upgrade, request.connection, request.version)
				}
				key, err := base64.StdEncoding.DecodeString(request.key)
				if err != nil || len(key) != 16 {
					t.Errorf("Sec-WebSocket-Key = %q, decoded length %d, err %v; want 16 random bytes", request.key, len(key), err)
				}
				if request.extensions != "" || request.contentLength != 0 || len(request.body) != 0 {
					t.Errorf("unexpected extensions or request body: extensions=%q contentLength=%d body=%q", request.extensions, request.contentLength, request.body)
				}
				if len(serverResult.extra) != 0 {
					t.Errorf("client sent data after the HTTP request: %x", serverResult.extra)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("fake WebSocket server did not finish")
			}
		})
	}
}

type fakeResponse struct {
	status int
	body   string
	accept string
	frame  []byte
}

type capturedRequest struct {
	method        string
	requestURI    string
	proto         string
	key           string
	upgrade       string
	connection    string
	version       string
	extensions    string
	contentLength int64
	body          []byte
}

type fakeServerResult struct {
	request capturedRequest
	extra   []byte
	err     error
}

func startFakeWebSocketServer(t *testing.T, response fakeResponse) (int, <-chan fakeServerResult) {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	result := make(chan fakeServerResult, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		serverResult := fakeServerResult{}
		defer func() {
			result <- serverResult
		}()
		connection, err := listener.Accept()
		if err != nil {
			serverResult.err = err
			return
		}
		defer connection.Close()
		_ = connection.SetDeadline(time.Now().Add(3 * time.Second))
		reader := bufio.NewReader(connection)
		request, err := http.ReadRequest(reader)
		if err != nil {
			serverResult.err = err
			return
		}
		serverResult.request = capturedRequest{
			method:        request.Method,
			requestURI:    request.RequestURI,
			proto:         request.Proto,
			key:           request.Header.Get("Sec-WebSocket-Key"),
			upgrade:       request.Header.Get("Upgrade"),
			connection:    request.Header.Get("Connection"),
			version:       request.Header.Get("Sec-WebSocket-Version"),
			extensions:    request.Header.Get("Sec-WebSocket-Extensions"),
			contentLength: request.ContentLength,
		}
		if request.Body != nil {
			serverResult.request.body, err = io.ReadAll(request.Body)
			_ = request.Body.Close()
			if err != nil {
				serverResult.err = err
				return
			}
		}
		if response.status == http.StatusSwitchingProtocols {
			accept := response.accept
			if accept == "" {
				accept = websocketAccept(serverResult.request.key)
			}
			headers := "HTTP/1.1 101 Switching Protocols\r\n" +
				"Upgrade: websocket\r\n" +
				"Connection: Upgrade\r\n" +
				"Sec-WebSocket-Accept: " + accept + "\r\n" +
				"\r\n"
			if err := writeTestAll(connection, []byte(headers)); err != nil {
				serverResult.err = err
				return
			}
			if err := writeTestAll(connection, response.frame); err != nil {
				serverResult.err = err
				return
			}
		} else {
			statusText := http.StatusText(response.status)
			headers := fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Length: %d\r\n\r\n%s", response.status, statusText, len(response.body), response.body)
			if err := writeTestAll(connection, []byte(headers)); err != nil {
				serverResult.err = err
				return
			}
		}
		serverResult.extra, err = io.ReadAll(reader)
		if err != nil && !errors.Is(err, io.EOF) {
			serverResult.err = err
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		select {
		case <-done:
		case <-time.After(4 * time.Second):
			t.Error("fake WebSocket server did not stop")
		}
	})
	return listener.Addr().(*net.TCPAddr).Port, result
}

func unmaskedFrame(fin bool, opcode byte, payload []byte) []byte {
	first := opcode
	if fin {
		first |= 0x80
	}
	frame := []byte{first}
	switch {
	case len(payload) <= 125:
		frame = append(frame, byte(len(payload)))
	case len(payload) <= 1<<16-1:
		frame = append(frame, 126, 0, 0)
		binary.BigEndian.PutUint16(frame[len(frame)-2:], uint16(len(payload)))
	default:
		frame = append(frame, 127, 0, 0, 0, 0, 0, 0, 0, 0)
		binary.BigEndian.PutUint64(frame[len(frame)-8:], uint64(len(payload)))
	}
	return append(frame, payload...)
}

func maskedFrame(payload []byte) []byte {
	frame := []byte{0x81, 0x80 | byte(len(payload)), 0, 0, 0, 0}
	return append(frame, payload...)
}

func websocketAccept(key string) string {
	hash := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func writeTestAll(writer io.Writer, data []byte) error {
	for len(data) > 0 {
		written, err := writer.Write(data)
		if written > 0 {
			data = data[written:]
		}
		if err != nil {
			return err
		}
		if written == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}
