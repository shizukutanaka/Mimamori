package vnc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/netip"
	"testing"
	"time"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantMajor int
		wantMinor int
		wantOK    bool
	}{
		{name: "standard", input: []byte("RFB 003.008\n"), wantMajor: 3, wantMinor: 8, wantOK: true},
		{name: "Apple high minor", input: []byte("RFB 003.889\n"), wantMajor: 3, wantMinor: 889, wantOK: true},
		{name: "higher major", input: []byte("RFB 005.000\n"), wantMajor: 5, wantMinor: 0, wantOK: true},
		{name: "zero version", input: []byte("RFB 000.000\n"), wantOK: true},
		{name: "wrong prefix", input: []byte("SSH 003.008\n")},
		{name: "wrong separator", input: []byte("RFB-003.008\n")},
		{name: "non-digit major", input: []byte("RFB 0a3.008\n")},
		{name: "non-digit minor", input: []byte("RFB 003.0a8\n")},
		{name: "missing newline", input: []byte("RFB 003.008")},
		{name: "extra byte", input: []byte("RFB 003.008\nx")},
		{name: "CRLF", input: []byte("RFB 003.008\r\n")},
		{name: "short", input: []byte("RFB")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			major, minor, ok := ParseVersion(test.input)
			if major != test.wantMajor || minor != test.wantMinor || ok != test.wantOK {
				t.Fatalf("ParseVersion(%q) = (%d, %d, %v), want (%d, %d, %v)",
					test.input, major, minor, ok, test.wantMajor, test.wantMinor, test.wantOK)
			}
		})
	}
}

func FuzzParseVersion(f *testing.F) {
	for _, seed := range [][]byte{
		[]byte("RFB 003.008\n"),
		[]byte("RFB 003.889\n"),
		[]byte(""),
		[]byte("SSH-2.0-x\r\n"),
		[]byte("RFB 999.999\n"),
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input []byte) {
		major, minor, ok := ParseVersion(input)
		if ok && (major < 0 || major > 999 || minor < 0 || minor > 999) {
			t.Fatalf("ParseVersion(%q) = (%d, %d, true), versions must be within 0..999", input, major, minor)
		}
	})
}

func TestNoAuthSecurityNegotiation(t *testing.T) {
	tests := []struct {
		name         string
		greeting     string
		expected     string
		response     []byte
		want         bool
		waitForClose bool
	}{
		{
			name:     "3.8 security list without None",
			greeting: "RFB 003.008\n", expected: "RFB 003.008\n",
			response: []byte{1, 2}, waitForClose: true,
		},
		{
			name:     "3.8 None",
			greeting: "RFB 003.008\n", expected: "RFB 003.008\n",
			response: []byte{1, 1}, want: true, waitForClose: true,
		},
		{
			name:     "3.8 None among multiple types",
			greeting: "RFB 003.008\n", expected: "RFB 003.008\n",
			response: []byte{2, 2, 1}, want: true, waitForClose: true,
		},
		{
			name:     "3.8 zero security types",
			greeting: "RFB 003.008\n", expected: "RFB 003.008\n",
			response: []byte{0}, waitForClose: true,
		},
		{
			name:     "3.7 None",
			greeting: "RFB 003.007\n", expected: "RFB 003.007\n",
			response: []byte{1, 1}, want: true, waitForClose: true,
		},
		{
			name:     "3.3 None",
			greeting: "RFB 003.003\n", expected: "RFB 003.003\n",
			response: []byte{0, 0, 0, 1}, want: true, waitForClose: true,
		},
		{
			name:     "3.3 password security",
			greeting: "RFB 003.003\n", expected: "RFB 003.003\n",
			response: []byte{0, 0, 0, 2}, waitForClose: true,
		},
		{
			name:     "RealVNC 5.0 selects 3.8",
			greeting: "RFB 005.000\n", expected: "RFB 003.008\n",
			response: []byte{1, 1}, want: true, waitForClose: true,
		},
		{
			name:     "3.5 selects 3.3",
			greeting: "RFB 003.005\n", expected: "RFB 003.003\n",
			response: []byte{0, 0, 0, 1}, want: true, waitForClose: true,
		},
		{
			name:     "unsupported 2.0",
			greeting: "RFB 002.000\n", waitForClose: true,
		},
		{
			name:     "non-RFB greeting",
			greeting: "SSH-2.0-x\r\n",
		},
		{
			name:     "truncated security list",
			greeting: "RFB 003.008\n", expected: "RFB 003.008\n",
			response: []byte{2, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ip, port, _, serverDone := startRFBServer(t, rfbScript{
				greeting:     test.greeting,
				expected:     test.expected,
				response:     test.response,
				waitForClose: test.waitForClose,
			})
			got := NoAuth(context.Background(), ip, port, time.Second)
			if got != test.want {
				t.Fatalf("NoAuth() = %v, want %v", got, test.want)
			}
			waitForRFBServer(t, serverDone)
		})
	}
}

func TestNoAuthSilentServerTimesOut(t *testing.T) {
	ip, port, _, serverDone := startRFBServer(t, rfbScript{waitForClose: true})
	timeout := 150 * time.Millisecond
	started := time.Now()
	if NoAuth(context.Background(), ip, port, timeout) {
		t.Fatal("NoAuth() accepted a silent server")
	}
	if elapsed := time.Since(started); elapsed > timeout+300*time.Millisecond {
		t.Fatalf("NoAuth() took %v with a %v timeout", elapsed, timeout)
	}
	waitForRFBServer(t, serverDone)
}

func TestNoAuthCanceledContext(t *testing.T) {
	ip, port, accepted, serverDone := startRFBServer(t, rfbScript{waitForClose: true})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan bool, 1)
	go func() {
		result <- NoAuth(ctx, ip, port, 5*time.Second)
	}()

	select {
	case <-accepted:
	case <-time.After(time.Second):
		t.Fatal("server did not accept the connection")
	}
	cancel()

	select {
	case got := <-result:
		if got {
			t.Fatal("NoAuth() accepted a canceled connection")
		}
	case <-time.After(time.Second):
		t.Fatal("NoAuth() did not return after context cancellation")
	}
	waitForRFBServer(t, serverDone)
}

func TestNoAuthRejectsInvalidArguments(t *testing.T) {
	ip := netip.MustParseAddr("127.0.0.1")
	for _, test := range []struct {
		name    string
		ctx     context.Context
		ip      netip.Addr
		port    int
		timeout time.Duration
	}{
		{name: "nil context", ctx: nil, ip: ip, port: 5900, timeout: time.Second},
		{name: "invalid IP", ctx: context.Background(), port: 5900, timeout: time.Second},
		{name: "zero port", ctx: context.Background(), ip: ip, timeout: time.Second},
		{name: "too large port", ctx: context.Background(), ip: ip, port: 65536, timeout: time.Second},
		{name: "zero timeout", ctx: context.Background(), ip: ip, port: 5900},
	} {
		t.Run(test.name, func(t *testing.T) {
			if NoAuth(test.ctx, test.ip, test.port, test.timeout) {
				t.Fatal("NoAuth() accepted invalid arguments")
			}
		})
	}
}

type rfbScript struct {
	greeting     string
	expected     string
	response     []byte
	waitForClose bool
}

func startRFBServer(t *testing.T, script rfbScript) (netip.Addr, int, <-chan struct{}, <-chan error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = listener.Close()
	})

	accepted := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer connection.Close()
		close(accepted)

		if script.greeting != "" {
			if err := writeTestData(connection, []byte(script.greeting)); err != nil {
				done <- err
				return
			}
		}
		if script.expected != "" {
			reply := make([]byte, len(script.expected))
			if _, err := io.ReadFull(connection, reply); err != nil {
				done <- err
				return
			}
			if !bytes.Equal(reply, []byte(script.expected)) {
				done <- fmt.Errorf("client reply = %q, want %q", reply, script.expected)
				return
			}
		}
		if len(script.response) > 0 {
			if err := writeTestData(connection, script.response); err != nil {
				done <- err
				return
			}
		}
		if script.waitForClose {
			_ = connection.SetReadDeadline(time.Now().Add(time.Second))
			var extra [1]byte
			n, err := connection.Read(extra[:])
			if n > 0 {
				done <- fmt.Errorf("client sent extra data %q", extra[:n])
				return
			}
			if !errors.Is(err, io.EOF) {
				done <- fmt.Errorf("connection did not close cleanly: %w", err)
				return
			}
		}
		done <- nil
	}()

	address := listener.Addr().(*net.TCPAddr)
	return netip.MustParseAddr(address.IP.String()), address.Port, accepted, done
}

func writeTestData(connection net.Conn, data []byte) error {
	for len(data) > 0 {
		written, err := connection.Write(data)
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

func waitForRFBServer(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("scripted RFB server did not finish")
	}
}
