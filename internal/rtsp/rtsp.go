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

var headerSeparator = []byte("\r\n\r\n")

func ParseDescribeResponse(b []byte) (status int, hasSDP bool, ok bool) {
	headerEnd := bytes.Index(b, headerSeparator)
	if headerEnd < 0 {
		return 0, false, false
	}

	header := b[:headerEnd]
	lineEnd := bytes.Index(header, []byte("\r\n"))
	statusLine := header
	var headerLines []byte
	if lineEnd >= 0 {
		statusLine = header[:lineEnd]
		headerLines = header[lineEnd+2:]
	}

	fields := strings.Fields(string(statusLine))
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

	if len(headerLines) > 0 {
		for _, line := range bytes.Split(headerLines, []byte("\r\n")) {
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
				headerEnd := bytes.Index(response, headerSeparator)
				if headerEnd >= 0 && len(response) > headerEnd+len(headerSeparator) {
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
