package smb

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"net/netip"
	"strconv"
	"time"
)

const (
	netBIOSHeaderSize        = 4
	smbHeaderSize            = 32
	maxSMBPayload            = 512
	smbCommandNegotiate byte = 0x72
)

// SMB1Enabled sends one SMB1 NEGOTIATE that offers a single dialect and reports a status-0 SMB1 reply that
// selects it. Servers with SMBv1 disabled answer with SMB2 or refuse the dialect, so this proves SMBv1 is
// still enabled (EternalBlue-class exposure).
func SMB1Enabled(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
	if ctx == nil || !ip.IsValid() || port < 1 || port > 65535 || timeout <= 0 {
		return false
	}

	deadline := time.Now().Add(timeout)
	address := net.JoinHostPort(ip.String(), strconv.Itoa(port))
	connection, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "tcp", address)
	if err != nil {
		return false
	}
	defer connection.Close()
	if err := connection.SetDeadline(deadline); err != nil {
		return false
	}
	stopCancellation := context.AfterFunc(ctx, func() {
		_ = connection.SetDeadline(time.Now())
	})
	defer stopCancellation()
	if ctx.Err() != nil {
		return false
	}

	if !writeAll(connection, negotiateRequest()) {
		return false
	}

	var sessionHeader [netBIOSHeaderSize]byte
	if _, err := io.ReadFull(connection, sessionHeader[:]); err != nil || sessionHeader[0] != 0x00 {
		return false
	}
	payloadLength := int(sessionHeader[1])<<16 | int(sessionHeader[2])<<8 | int(sessionHeader[3])
	if payloadLength < 37 {
		return false
	}
	if payloadLength > maxSMBPayload {
		payloadLength = maxSMBPayload
	}

	payload := make([]byte, payloadLength)
	if _, err := io.ReadFull(connection, payload); err != nil || ctx.Err() != nil {
		return false
	}
	return ParseNegotiateResponse(payload)
}

func ParseNegotiateResponse(payload []byte) bool {
	return len(payload) >= 37 &&
		payload[0] == 0xff && payload[1] == 'S' && payload[2] == 'M' && payload[3] == 'B' &&
		payload[4] == smbCommandNegotiate &&
		binary.LittleEndian.Uint32(payload[5:9]) == 0 &&
		payload[32] >= 1 &&
		binary.LittleEndian.Uint16(payload[33:35]) == 0
}

func negotiateRequest() []byte {
	dialect := []byte{0x02, 'N', 'T', ' ', 'L', 'M', ' ', '0', '.', '1', '2', 0}
	payloadLength := smbHeaderSize + 1 + 2 + len(dialect)
	request := make([]byte, netBIOSHeaderSize+payloadLength)
	request[1] = byte(payloadLength >> 16)
	request[2] = byte(payloadLength >> 8)
	request[3] = byte(payloadLength)

	header := request[netBIOSHeaderSize : netBIOSHeaderSize+smbHeaderSize]
	copy(header[:4], []byte{0xff, 'S', 'M', 'B'})
	header[4] = smbCommandNegotiate
	header[9] = 0x18
	binary.LittleEndian.PutUint16(header[10:12], 0x0001)
	binary.LittleEndian.PutUint16(header[26:28], 0x1234)
	binary.LittleEndian.PutUint16(header[30:32], 0x0001)

	body := request[netBIOSHeaderSize+smbHeaderSize:]
	binary.LittleEndian.PutUint16(body[1:3], uint16(len(dialect)))
	copy(body[3:], dialect)
	return request
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
