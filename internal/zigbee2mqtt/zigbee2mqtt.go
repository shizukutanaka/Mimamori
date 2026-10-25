package zigbee2mqtt

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

const (
	maxFrameBytes = 1 << 20
	websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

func ParseMessage(payload []byte) bool {
	var message map[string]json.RawMessage
	if err := json.Unmarshal(payload, &message); err != nil || message == nil {
		return false
	}
	rawTopic, ok := message["topic"]
	if !ok {
		return false
	}
	var topic string
	if err := json.Unmarshal(rawTopic, &topic); err != nil || !strings.HasPrefix(topic, "bridge/") {
		return false
	}
	_, ok = message["payload"]
	return ok
}

// Unauthenticated opens the frontend WebSocket (GET /api with Upgrade) and reads one frame. A text frame
// with a bridge/… topic means the frontend accepted the connection without auth_token (the default is
// none); with a token set the handshake is rejected.
func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.Is4() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}

	deadline := time.Now().Add(timeout)
	address := net.JoinHostPort(ip.String(), strconv.Itoa(port))
	connection, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "tcp", address)
	if err != nil {
		return false
	}
	defer connection.Close()
	_ = connection.SetDeadline(deadline)
	stopCancellation := context.AfterFunc(ctx, func() {
		_ = connection.SetDeadline(time.Now())
	})
	defer stopCancellation()

	var randomKey [16]byte
	if _, err := rand.Read(randomKey[:]); err != nil {
		return false
	}
	key := base64.StdEncoding.EncodeToString(randomKey[:])
	request := "GET /api HTTP/1.1\r\n" +
		"Host: " + address + "\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + key + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"\r\n"
	if !writeAll(connection, []byte(request)) {
		return false
	}

	reader := bufio.NewReader(connection)
	response, err := http.ReadResponse(reader, nil)
	if err != nil || response.StatusCode != http.StatusSwitchingProtocols {
		return false
	}
	if response.Header.Get("Sec-WebSocket-Accept") != acceptKey(key) {
		return false
	}

	message, ok := readTextFrame(reader)
	return ok && ParseMessage(message)
}

func acceptKey(key string) string {
	hash := sha1.Sum([]byte(key + websocketGUID))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func readTextFrame(reader io.Reader) ([]byte, bool) {
	var header [2]byte
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return nil, false
	}
	if header[0]&0x80 == 0 || header[0]&0x70 != 0 || header[0]&0x0f != 1 || header[1]&0x80 != 0 {
		return nil, false
	}

	lengthCode := header[1] & 0x7f
	payloadLength := uint64(lengthCode)
	switch lengthCode {
	case 126:
		var extendedLength [2]byte
		if _, err := io.ReadFull(reader, extendedLength[:]); err != nil {
			return nil, false
		}
		payloadLength = uint64(binary.BigEndian.Uint16(extendedLength[:]))
		if payloadLength < 126 {
			return nil, false
		}
	case 127:
		var extendedLength [8]byte
		if _, err := io.ReadFull(reader, extendedLength[:]); err != nil {
			return nil, false
		}
		if extendedLength[0]&0x80 != 0 {
			return nil, false
		}
		payloadLength = binary.BigEndian.Uint64(extendedLength[:])
		if payloadLength < 1<<16 {
			return nil, false
		}
	}
	if payloadLength > maxFrameBytes {
		return nil, false
	}

	payload := make([]byte, int(payloadLength))
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, false
	}
	return payload, true
}

func writeAll(writer io.Writer, data []byte) bool {
	for len(data) > 0 {
		written, err := writer.Write(data)
		if written > 0 {
			data = data[written:]
		}
		if err != nil || written == 0 {
			return false
		}
	}
	return true
}
