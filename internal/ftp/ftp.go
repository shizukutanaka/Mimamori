package ftp

import (
	"bufio"
	"context"
	"net"
	"net/netip"
	"strconv"
	"time"
)

const (
	maxReplyLineLength = 1024
	maxReplyLines      = 100
)

func ReadReply(reader *bufio.Reader) (code int, ok bool) {
	if reader == nil {
		return 0, false
	}
	first, ok := readReplyLine(reader)
	if !ok || len(first) < 3 || first[0] < '0' || first[0] > '9' || first[1] < '0' || first[1] > '9' || first[2] < '0' || first[2] > '9' {
		return 0, false
	}
	code = int(first[0]-'0')*100 + int(first[1]-'0')*10 + int(first[2]-'0')
	if code < 100 || code > 599 {
		return 0, false
	}
	if len(first) == 3 {
		return code, true
	}
	switch first[3] {
	case ' ':
		return code, true
	case '-':
		for lines := 1; lines < maxReplyLines; lines++ {
			line, ok := readReplyLine(reader)
			if !ok {
				return 0, false
			}
			if len(line) >= 4 && line[0] == first[0] && line[1] == first[1] && line[2] == first[2] && line[3] == ' ' {
				return code, true
			}
		}
	}
	return 0, false
}

func readReplyLine(reader *bufio.Reader) ([]byte, bool) {
	line := make([]byte, 0, 128)
	for {
		character, err := reader.ReadByte()
		if err != nil {
			return nil, false
		}
		if character == '\n' {
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			return line, len(line) <= maxReplyLineLength
		}
		if len(line) >= maxReplyLineLength+1 || len(line) == maxReplyLineLength && character != '\r' {
			return nil, false
		}
		line = append(line, character)
	}
}

// Anonymous logs in with USER anonymous / PASS mimamori@, reports a 230 reply and sends QUIT. Servers
// answer USER anonymous with 331 whether or not anonymous access is allowed, so only the PASS reply tells.
func Anonymous(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
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

	reader := bufio.NewReader(connection)
	greeting, ok := ReadReply(reader)
	if !ok {
		return false
	}
	defer func() {
		_ = writeAll(connection, []byte("QUIT\r\n"))
	}()
	if greeting != 220 {
		return false
	}

	if !writeAll(connection, []byte("USER anonymous\r\n")) {
		return false
	}
	code, ok := ReadReply(reader)
	if !ok {
		return false
	}
	if code == 230 {
		return true
	}
	if code != 331 || !writeAll(connection, []byte("PASS mimamori@\r\n")) {
		return false
	}
	code, ok = ReadReply(reader)
	return ok && code == 230
}

func writeAll(writer net.Conn, data []byte) bool {
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
