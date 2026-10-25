package mqtt

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"
	"net/netip"
	"strconv"
	"time"
)

func ParseConnack(b []byte) (code byte, ok bool) {
	if len(b) < 4 || b[0] != 0x20 || b[1] != 0x02 {
		return 0, false
	}
	return b[3], true
}

// Anonymous sends an MQTT 3.1.1 CONNECT without credentials, reports CONNACK code 0 and sends
// DISCONNECT. Mosquitto answers code 5 (not authorized) unless allow_anonymous is true.
func Anonymous(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.IsValid() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}

	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return false
	}
	clientID := "mimamori-" + hex.EncodeToString(suffix[:])
	connect := connectPacket(clientID)

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

	if !writeAll(connection, connect) {
		return false
	}
	response := make([]byte, 4)
	if _, err := io.ReadFull(connection, response); err != nil {
		return false
	}
	_ = writeAll(connection, []byte{0xe0, 0x00})

	code, ok := ParseConnack(response)
	return ok && code == 0
}

func connectPacket(clientID string) []byte {
	body := make([]byte, 0, 29)
	body = append(body, 0x00, 0x04)
	body = append(body, "MQTT"...)
	body = append(body, 0x04, 0x02, 0x00, 0x0a)
	var length [2]byte
	binary.BigEndian.PutUint16(length[:], uint16(len(clientID)))
	body = append(body, length[:]...)
	body = append(body, clientID...)
	return append([]byte{0x10, byte(len(body))}, body...)
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
