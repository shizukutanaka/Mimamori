package vnc

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"net/netip"
	"strconv"
	"time"
)

func ParseVersion(version []byte) (major, minor int, ok bool) {
	if len(version) != 12 ||
		version[0] != 'R' || version[1] != 'F' || version[2] != 'B' || version[3] != ' ' ||
		version[7] != '.' || version[11] != '\n' {
		return 0, 0, false
	}
	for _, index := range []int{4, 5, 6, 8, 9, 10} {
		if version[index] < '0' || version[index] > '9' {
			return 0, 0, false
		}
	}

	major = int(version[4]-'0')*100 + int(version[5]-'0')*10 + int(version[6]-'0')
	minor = int(version[8]-'0')*100 + int(version[9]-'0')*10 + int(version[10]-'0')
	return major, minor, true
}

// NoAuth reads the RFB version, echoes it and reports a server that offers security type 1 (None). Per RFC
// 6143 servers offer None only when no password is configured; 3.3 servers send a single type, 3.7+ a list.
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

	version := make([]byte, 12)
	if _, err := io.ReadFull(connection, version); err != nil {
		return false
	}
	major, minor, ok := ParseVersion(version)
	if !ok {
		return false
	}

	clientVersion, modern := clientVersionFor(major, minor)
	if len(clientVersion) == 0 || !writeAll(connection, clientVersion) {
		return false
	}
	if !modern {
		var securityType [4]byte
		if _, err := io.ReadFull(connection, securityType[:]); err != nil {
			return false
		}
		return binary.BigEndian.Uint32(securityType[:]) == 1
	}

	var count [1]byte
	if _, err := io.ReadFull(connection, count[:]); err != nil || count[0] == 0 {
		return false
	}
	securityTypes := make([]byte, int(count[0]))
	if _, err := io.ReadFull(connection, securityTypes); err != nil {
		return false
	}
	for _, securityType := range securityTypes {
		if securityType == 1 {
			return true
		}
	}
	return false
}

func clientVersionFor(major, minor int) ([]byte, bool) {
	switch {
	case major > 3 || major == 3 && minor >= 8:
		return []byte("RFB 003.008\n"), true
	case major == 3 && minor >= 7:
		return []byte("RFB 003.007\n"), true
	case major == 3 && minor >= 3:
		return []byte("RFB 003.003\n"), false
	default:
		return nil, false
	}
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
