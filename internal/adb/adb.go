package adb

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
	commandCNXN uint32 = 0x4e584e43
	commandAUTH uint32 = 0x48545541
	commandSTLS uint32 = 0x534c5453
	headerSize         = 24
)

const (
	// adbd negotiates the minimum of both sides for each value.
	cnxnVersion uint32 = 0x01000001
	cnxnMaxData uint32 = 256 * 1024
)

func ParseHeader(b []byte) (command, arg0, arg1, length uint32, ok bool) {
	if len(b) != headerSize {
		return 0, 0, 0, 0, false
	}

	command = binary.LittleEndian.Uint32(b[0:4])
	arg0 = binary.LittleEndian.Uint32(b[4:8])
	arg1 = binary.LittleEndian.Uint32(b[8:12])
	length = binary.LittleEndian.Uint32(b[12:16])
	magic := binary.LittleEndian.Uint32(b[20:24])
	if magic != command^uint32(0xffffffff) {
		return 0, 0, 0, 0, false
	}
	return command, arg0, arg1, length, true
}

// NoAuth sends one ADB CNXN (host::) and reports a daemon that answers CNXN instead of AUTH. adbd asks
// for AUTH unless authentication is disabled (userdebug/eng builds or an unlocked device), so a CNXN reply
// means anyone on the LAN gets a shell.
func NoAuth(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
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
	_ = connection.SetDeadline(deadline)
	stopCancellation := context.AfterFunc(ctx, func() {
		_ = connection.SetDeadline(time.Now())
	})
	defer stopCancellation()

	payload := []byte("host::\x00")
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint32(header[0:4], commandCNXN)
	binary.LittleEndian.PutUint32(header[4:8], cnxnVersion)
	binary.LittleEndian.PutUint32(header[8:12], cnxnMaxData)
	binary.LittleEndian.PutUint32(header[12:16], uint32(len(payload)))
	var dataCheck uint32
	for _, value := range payload {
		dataCheck += uint32(value)
	}
	binary.LittleEndian.PutUint32(header[16:20], dataCheck)
	binary.LittleEndian.PutUint32(header[20:24], commandCNXN^uint32(0xffffffff))
	if !writeAll(connection, header) || !writeAll(connection, payload) {
		return false
	}

	reply := make([]byte, headerSize)
	if _, err := io.ReadFull(connection, reply); err != nil {
		return false
	}
	command, _, _, _, ok := ParseHeader(reply)
	return ok && command == commandCNXN
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
