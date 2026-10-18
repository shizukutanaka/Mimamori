package rtsp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strings"
	"testing"
	"time"
)

func TestParseDescribeResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		status   int
		hasSDP   bool
		ok       bool
	}{
		{
			name:     "200 with SDP",
			response: "RTSP/1.0 200 OK\r\nCSeq: 1\r\nContent-Type: application/sdp; charset=utf-8\r\n\r\nv=0\r\n",
			status:   200,
			hasSDP:   true,
			ok:       true,
		},
		{
			name:     "200 without SDP",
			response: "RTSP/1.0 200 OK\r\nContent-Type: text/plain\r\n\r\nok",
			status:   200,
			ok:       true,
		},
		{
			name:     "401",
			response: "RTSP/1.0 401 Unauthorized\r\nWWW-Authenticate: Basic realm=\"camera\"\r\n\r\n",
			status:   401,
			ok:       true,
		},
		{
			name:     "404",
			response: "RTSP/1.0 404 Not Found\r\n\r\n",
			status:   404,
			ok:       true,
		},
		{name: "garbage", response: "not an RTSP response\r\n\r\n"},
		{name: "truncated", response: "RTSP/1.0 200 OK\r\nContent-Type: application/sdp\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, hasSDP, ok := ParseDescribeResponse([]byte(test.response))
			if status != test.status || hasSDP != test.hasSDP || ok != test.ok {
				t.Fatalf("ParseDescribeResponse() = (%d, %t, %t), want (%d, %t, %t)",
					status, hasSDP, ok, test.status, test.hasSDP, test.ok)
			}
		})
	}
}

func TestUnauthenticated(t *testing.T) {
	const sdpBody = "v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\ns=Camera\r\nt=0 0\r\n"
	response := fmt.Sprintf(
		"RTSP/1.0 200 OK\r\nCSeq: 1\r\nContent-Type: application/sdp\r\nContent-Length: %d\r\n\r\n%s",
		len(sdpBody), sdpBody,
	)
	tests := []struct {
		name     string
		response string
		timeout  time.Duration
		want     bool
		silent   bool
	}{
		{name: "200 with SDP", response: response, timeout: time.Second, want: true},
		{
			name:     "401",
			response: "RTSP/1.0 401 Unauthorized\r\nCSeq: 1\r\nWWW-Authenticate: Basic\r\n\r\n",
			timeout:  time.Second,
		},
		{name: "silent server", timeout: 200 * time.Millisecond, silent: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listener, requests, done := startFakeRTSPServer(t, test.response, test.silent)
			defer listener.Close()
			port := listener.Addr().(*net.TCPAddr).Port
			ip := netip.MustParseAddr("127.0.0.1")

			started := time.Now()
			got := Unauthenticated(context.Background(), ip, port, test.timeout)
			elapsed := time.Since(started)
			if got != test.want {
				t.Fatalf("Unauthenticated() = %t, want %t", got, test.want)
			}
			if test.silent && elapsed > time.Second {
				t.Fatalf("silent response took %v, want return within the timeout", elapsed)
			}

			wantRequest := fmt.Sprintf(
				"DESCRIBE rtsp://127.0.0.1:%d/ RTSP/1.0\r\nCSeq: 1\r\nAccept: application/sdp\r\nUser-Agent: Mimamori\r\n\r\n",
				port,
			)
			select {
			case request := <-requests:
				if request != wantRequest {
					t.Fatalf("request = %q, want %q", request, wantRequest)
				}
			case <-time.After(time.Second):
				t.Fatal("fake server did not receive the request")
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("fake server did not finish")
			}
		})
	}
}

func startFakeRTSPServer(t *testing.T, response string, silent bool) (net.Listener, <-chan string, <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot listen on loopback: %v", err)
	}

	requests := make(chan string, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		defer connection.Close()

		var request strings.Builder
		reader := bufio.NewReader(connection)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			request.WriteString(line)
			if line == "\r\n" {
				break
			}
		}
		requests <- request.String()

		if silent {
			_, _ = io.Copy(io.Discard, connection)
			return
		}
		_, _ = io.WriteString(connection, response)
	}()
	return listener, requests, done
}
