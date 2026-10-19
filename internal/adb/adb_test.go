package adb

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/netip"
	"testing"
	"time"
)

func TestParseHeader(t *testing.T) {
	valid := encodeHeader(commandCNXN, cnxnVersion, cnxnMaxData, 7, 562)
	tests := []struct {
		name string
		data []byte
		ok   bool
	}{
		{name: "valid", data: valid, ok: true},
		{name: "short", data: valid[:headerSize-1]},
		{name: "long", data: append(valid, 0)},
		{name: "bad magic", data: append([]byte(nil), valid[:headerSize]...)},
	}
	tests[3].data[20] ^= 1

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command, arg0, arg1, length, ok := ParseHeader(test.data)
			if ok != test.ok {
				t.Fatalf("ParseHeader() ok = %t, want %t", ok, test.ok)
			}
			if !ok {
				return
			}
			if command != commandCNXN || arg0 != cnxnVersion || arg1 != cnxnMaxData || length != 7 {
				t.Fatalf("ParseHeader() = (%#x, %#x, %#x, %d), want CNXN header", command, arg0, arg1, length)
			}
		})
	}
}

func TestNoAuthReplies(t *testing.T) {
	tests := []struct {
		name     string
		response []byte
		hold     bool
		want     bool
	}{
		{name: "CNXN", response: encodeHeader(commandCNXN, 0, 0, 0, 0), want: true},
		{name: "AUTH", response: encodeHeader(commandAUTH, 1, 0, 0, 0)},
		{name: "STLS", response: encodeHeader(commandSTLS, 0, 0, 0, 0)},
		{name: "bad magic", response: invalidMagicHeader(commandCNXN)},
		{name: "close", response: nil},
		{name: "short header", response: encodeHeader(commandCNXN, 0, 0, 0, 0)[:8]},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port, ready, serverErr := startFakeServer(t, test.response, test.hold)
			got := NoAuth(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second)
			if got != test.want {
				t.Fatalf("NoAuth() = %t, want %t", got, test.want)
			}
			waitForRequest(t, ready)
			waitForServer(t, serverErr)
		})
	}
}

func TestNoAuthTimeout(t *testing.T) {
	port, ready, serverErr := startFakeServer(t, nil, true)
	timeout := 150 * time.Millisecond
	start := time.Now()
	got := NoAuth(context.Background(), netip.MustParseAddr("127.0.0.1"), port, timeout)
	elapsed := time.Since(start)
	if got {
		t.Fatal("NoAuth() = true, want false for silent server")
	}
	if elapsed > timeout+500*time.Millisecond {
		t.Fatalf("NoAuth() took %v, want no more than timeout plus 500ms", elapsed)
	}
	waitForRequest(t, ready)
	waitForServer(t, serverErr)
}

func TestNoAuthCanceledContext(t *testing.T) {
	port, ready, serverErr := startFakeServer(t, nil, true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan bool, 1)
	go func() {
		result <- NoAuth(ctx, netip.MustParseAddr("127.0.0.1"), port, time.Second)
	}()
	waitForRequest(t, ready)
	cancel()
	select {
	case got := <-result:
		if got {
			t.Fatal("NoAuth() = true, want false after cancellation")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("NoAuth() did not return promptly after cancellation")
	}
	waitForServer(t, serverErr)
}

func FuzzParseHeader(f *testing.F) {
	f.Add(encodeHeader(commandCNXN, cnxnVersion, cnxnMaxData, 7, 562))
	f.Add([]byte{})
	f.Add(make([]byte, headerSize))
	f.Fuzz(func(t *testing.T, data []byte) {
		command, arg0, arg1, length, ok := ParseHeader(data)
		if !ok {
			return
		}
		if len(data) != headerSize {
			t.Fatalf("ParseHeader() accepted %d bytes", len(data))
		}
		if binary.LittleEndian.Uint32(data[20:24]) != command^uint32(0xffffffff) {
			t.Fatal("ParseHeader() accepted invalid magic")
		}
		if binary.LittleEndian.Uint32(data[4:8]) != arg0 ||
			binary.LittleEndian.Uint32(data[8:12]) != arg1 ||
			binary.LittleEndian.Uint32(data[12:16]) != length {
			t.Fatal("ParseHeader() returned values inconsistent with the header")
		}
	})
}

func startFakeServer(t *testing.T, response []byte, hold bool) (int, <-chan struct{}, <-chan error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	ready := make(chan struct{})
	serverErr := make(chan error, 1)
	go func() {
		defer listener.Close()
		connection, err := listener.Accept()
		if err != nil {
			serverErr <- err
			return
		}
		defer connection.Close()

		header := make([]byte, headerSize)
		if _, err := io.ReadFull(connection, header); err != nil {
			serverErr <- err
			return
		}
		command, arg0, arg1, length, ok := ParseHeader(header)
		payload := []byte("host::\x00")
		var expectedCheck uint32
		for _, value := range payload {
			expectedCheck += uint32(value)
		}
		dataCheck := binary.LittleEndian.Uint32(header[16:20])
		if !ok || command != commandCNXN || arg0 != cnxnVersion || arg1 != cnxnMaxData ||
			length != uint32(len(payload)) || dataCheck != expectedCheck {
			serverErr <- fmt.Errorf("invalid CNXN header: %#v", header)
			return
		}
		requestPayload := make([]byte, length)
		if _, err := io.ReadFull(connection, requestPayload); err != nil {
			serverErr <- err
			return
		}
		if !bytes.Equal(requestPayload, payload) {
			serverErr <- fmt.Errorf("request payload = %q, want %q", requestPayload, payload)
			return
		}
		close(ready)
		if len(response) > 0 {
			_, err := connection.Write(response)
			serverErr <- err
			return
		}
		if hold {
			_, err := io.Copy(io.Discard, connection)
			if err != nil {
				serverErr <- err
				return
			}
		}
		serverErr <- nil
	}()
	return port, ready, serverErr
}

func waitForRequest(t *testing.T, ready <-chan struct{}) {
	t.Helper()
	select {
	case <-ready:
	case <-time.After(time.Second):
		t.Fatal("server did not receive a valid CNXN request")
	}
}

func waitForServer(t *testing.T, serverErr <-chan error) {
	t.Helper()
	select {
	case err := <-serverErr:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("fake ADB server did not finish")
	}
}

func encodeHeader(command, arg0, arg1, length, dataCheck uint32) []byte {
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint32(header[0:4], command)
	binary.LittleEndian.PutUint32(header[4:8], arg0)
	binary.LittleEndian.PutUint32(header[8:12], arg1)
	binary.LittleEndian.PutUint32(header[12:16], length)
	binary.LittleEndian.PutUint32(header[16:20], dataCheck)
	binary.LittleEndian.PutUint32(header[20:24], command^uint32(0xffffffff))
	return header
}

func invalidMagicHeader(command uint32) []byte {
	header := encodeHeader(command, 0, 0, 0, 0)
	header[20] ^= 1
	return header
}
