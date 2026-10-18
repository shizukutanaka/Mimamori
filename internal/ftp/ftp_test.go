package ftp

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/netip"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestReadReply(t *testing.T) {
	tooManyLines := "220-start\r\n" + strings.Repeat("body\r\n", maxReplyLines) + "220 done\r\n"
	tests := []struct {
		name     string
		response string
		code     int
		ok       bool
	}{
		{name: "single-line CRLF", response: "220 ready\r\n", code: 220, ok: true},
		{name: "single-line LF", response: "230 logged in\n", code: 230, ok: true},
		{name: "multi-line", response: "220-start\r\n230 x\r\n220 ready\r\n", code: 220, ok: true},
		{name: "code without text", response: "331\n", code: 331, ok: true},
		{name: "non-digit code", response: "x20 invalid\r\n", ok: false},
		{name: "code below range", response: "099 invalid\r\n", ok: false},
		{name: "code above range", response: "600 invalid\r\n", ok: false},
		{name: "invalid separator", response: "220x invalid\r\n", ok: false},
		{name: "EOF during multi-line", response: "220-start\r\nbody", ok: false},
		{name: "overlong line", response: strings.Repeat("x", maxReplyLineLength+1) + "\r\n", ok: false},
		{name: "more than 100 lines", response: tooManyLines, ok: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, ok := ReadReply(bufio.NewReader(strings.NewReader(test.response)))
			if code != test.code || ok != test.ok {
				t.Fatalf("ReadReply() = (%d, %t), want (%d, %t)", code, ok, test.code, test.ok)
			}
		})
	}
}

func FuzzReadReply(f *testing.F) {
	for _, response := range []string{
		"220 ready\r\n",
		"230 logged in\n",
		"220-start\r\n230 x\r\n220 ready\r\n",
		"x20 invalid\r\n",
		"099 invalid\r\n",
		"600 invalid\r\n",
		"220-start\r\nbody",
		strings.Repeat("x", maxReplyLineLength+1) + "\r\n",
		"220-start\r\n" + strings.Repeat("body\r\n", maxReplyLines),
	} {
		f.Add(response)
	}

	f.Fuzz(func(t *testing.T, response string) {
		code, ok := ReadReply(bufio.NewReader(strings.NewReader(response)))
		if ok && (code < 100 || code > 599) {
			t.Fatalf("ReadReply() returned out-of-range code %d", code)
		}
	})
}

type ftpServerStep struct {
	command  string
	response string
}

type ftpServerResult struct {
	commands []string
	err      error
}

func startFTPServer(t *testing.T, greeting string, steps []ftpServerStep) (netip.Addr, int, <-chan ftpServerResult) {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	result := make(chan ftpServerResult, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			result <- ftpServerResult{err: err}
			return
		}
		defer connection.Close()
		_ = connection.SetDeadline(time.Now().Add(2 * time.Second))
		if !writeAll(connection, []byte(greeting)) {
			result <- ftpServerResult{err: io.ErrShortWrite}
			return
		}
		reader := bufio.NewReader(connection)
		commands := make([]string, 0, len(steps))
		for _, step := range steps {
			command, err := reader.ReadString('\n')
			if err != nil {
				result <- ftpServerResult{commands: commands, err: err}
				return
			}
			commands = append(commands, command)
			if command != step.command {
				result <- ftpServerResult{commands: commands, err: io.ErrUnexpectedEOF}
				return
			}
			if step.response != "" && !writeAll(connection, []byte(step.response)) {
				result <- ftpServerResult{commands: commands, err: io.ErrShortWrite}
				return
			}
		}
		result <- ftpServerResult{commands: commands}
	}()
	address := listener.Addr().(*net.TCPAddr)
	return netip.MustParseAddr("127.0.0.1"), address.Port, result
}

func TestAnonymousNegotiation(t *testing.T) {
	tests := []struct {
		name     string
		greeting string
		steps    []ftpServerStep
		want     bool
	}{
		{
			name:     "password accepted",
			greeting: "220 ready\r\n",
			steps: []ftpServerStep{
				{command: "USER anonymous\r\n", response: "331 password required\r\n"},
				{command: "PASS mimamori@\r\n", response: "230 logged in\r\n"},
				{command: "QUIT\r\n"},
			},
			want: true,
		},
		{
			name:     "accepted after USER",
			greeting: "220 ready\r\n",
			steps: []ftpServerStep{
				{command: "USER anonymous\r\n", response: "230 logged in\r\n"},
				{command: "QUIT\r\n"},
			},
			want: true,
		},
		{
			name:     "password rejected",
			greeting: "220 ready\r\n",
			steps: []ftpServerStep{
				{command: "USER anonymous\r\n", response: "331 password required\r\n"},
				{command: "PASS mimamori@\r\n", response: "530 not authorized\r\n"},
				{command: "QUIT\r\n"},
			},
			want: false,
		},
		{
			name:     "user rejected",
			greeting: "220 ready\r\n",
			steps: []ftpServerStep{
				{command: "USER anonymous\r\n", response: "530 not authorized\r\n"},
				{command: "QUIT\r\n"},
			},
			want: false,
		},
		{
			name:     "bad greeting",
			greeting: "421 service unavailable\r\n",
			steps:    []ftpServerStep{{command: "QUIT\r\n"}},
			want:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ip, port, resultCh := startFTPServer(t, test.greeting, test.steps)
			if got := Anonymous(context.Background(), ip, port, time.Second); got != test.want {
				t.Fatalf("Anonymous() = %t, want %t", got, test.want)
			}
			result := <-resultCh
			if result.err != nil {
				t.Fatal(result.err)
			}
			wantCommands := make([]string, len(test.steps))
			for i, step := range test.steps {
				wantCommands[i] = step.command
			}
			if !reflect.DeepEqual(result.commands, wantCommands) {
				t.Fatalf("server saw commands %q, want %q", result.commands, wantCommands)
			}
		})
	}
}

func startSilentFTPServer(t *testing.T) (netip.Addr, int, <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	accepted := make(chan struct{})
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		close(accepted)
		defer connection.Close()
		_ = connection.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, _ = io.Copy(io.Discard, connection)
	}()
	address := listener.Addr().(*net.TCPAddr)
	return netip.MustParseAddr("127.0.0.1"), address.Port, accepted
}

func TestAnonymousSilentServerTimesOut(t *testing.T) {
	ip, port, accepted := startSilentFTPServer(t)
	started := time.Now()
	if Anonymous(context.Background(), ip, port, 100*time.Millisecond) {
		t.Fatal("Anonymous() = true, want false")
	}
	if elapsed := time.Since(started); elapsed > 600*time.Millisecond {
		t.Fatalf("Anonymous() took %s, want no more than timeout + 500ms", elapsed)
	}
	select {
	case <-accepted:
	case <-time.After(time.Second):
		t.Fatal("FTP client did not connect to the silent server")
	}
}

func TestAnonymousCanceledContext(t *testing.T) {
	ip, port, accepted := startSilentFTPServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan bool, 1)
	go func() {
		result <- Anonymous(ctx, ip, port, 2*time.Second)
	}()
	select {
	case <-accepted:
	case <-time.After(time.Second):
		cancel()
		t.Fatal("FTP client did not connect to the silent server")
	}
	cancel()
	select {
	case got := <-result:
		if got {
			t.Fatal("Anonymous() = true, want false for canceled context")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Anonymous() did not return promptly after cancellation")
	}
}
