package rtsp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

const maxResponseBytes = 8 * 1024

// RFC 2326 §4 recommends accepting LF as a line terminator.
func headerEnd(b []byte) (headerLength, bodyStart int, ok bool) {
	lineStart := 0
	for index, value := range b {
		if value != '\n' {
			continue
		}
		line := b[lineStart:index]
		if len(line) == 0 || (len(line) == 1 && line[0] == '\r') {
			return lineStart, index + 1, true
		}
		lineStart = index + 1
	}
	return 0, 0, false
}

func ParseDescribeResponse(b []byte) (status int, hasSDP bool, ok bool) {
	headerLength, _, found := headerEnd(b)
	if !found {
		return 0, false, false
	}

	lines := bytes.Split(b[:headerLength], []byte("\n"))
	if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	for index, line := range lines {
		if len(line) > 0 && line[len(line)-1] == '\r' {
			lines[index] = line[:len(line)-1]
		}
	}
	if len(lines) == 0 {
		return 0, false, false
	}

	fields := strings.Fields(string(lines[0]))
	if len(fields) < 3 || fields[0] != "RTSP/1.0" || len(fields[1]) != 3 {
		return 0, false, false
	}
	for _, digit := range fields[1] {
		if digit < '0' || digit > '9' {
			return 0, false, false
		}
	}
	status, err := strconv.Atoi(fields[1])
	if err != nil || status < 100 || status > 999 {
		return 0, false, false
	}

	for _, line := range lines[1:] {
		colon := bytes.IndexByte(line, ':')
		if colon <= 0 {
			return 0, false, false
		}
		name := strings.TrimSpace(string(line[:colon]))
		if name == "" {
			return 0, false, false
		}
		if strings.EqualFold(name, "Content-Type") &&
			strings.Contains(strings.ToLower(string(line[colon+1:])), "application/sdp") {
			hasSDP = true
		}
	}
	return status, hasSDP, true
}

func Unauthenticated(ctx context.Context, ip netip.Addr, port int, timeout time.Duration) bool {
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

	request := fmt.Sprintf(
		"DESCRIBE rtsp://%s/ RTSP/1.0\r\nCSeq: 1\r\nAccept: application/sdp\r\nUser-Agent: Mimamori\r\n\r\n",
		address,
	)
	written, err := io.WriteString(connection, request)
	if err != nil || written != len(request) {
		return false
	}

	response := make([]byte, 0, maxResponseBytes)
	buffer := make([]byte, 1024)
	for len(response) < maxResponseBytes {
		readSize := len(buffer)
		if remaining := maxResponseBytes - len(response); remaining < readSize {
			readSize = remaining
		}
		n, readErr := connection.Read(buffer[:readSize])
		if n > 0 {
			response = append(response, buffer[:n]...)
			status, hasSDP, ok := ParseDescribeResponse(response)
			if ok {
				if status != 200 || !hasSDP {
					return false
				}
				_, bodyStart, found := headerEnd(response)
				if found && len(response) > bodyStart {
					return true
				}
			}
		}
		if readErr != nil {
			return false
		}
	}
	return false
}
