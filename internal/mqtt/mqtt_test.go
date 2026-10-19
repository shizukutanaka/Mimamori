package mqtt

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"
	"net/netip"
	"strings"
	"testing"
	"time"
)

func TestParseConnack(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		code byte
		ok   bool
	}{
		{name: "accepted", data: []byte{0x20, 0x02, 0x00, 0x00}, code: 0x00, ok: true},
		{name: "not authorized", data: []byte{0x20, 0x02, 0x00, 0x05}, code: 0x05, ok: true},
		{name: "bad protocol", data: []byte{0x20, 0x02, 0x00, 0x04}, code: 0x04, ok: true},
		{name: "wrong packet type", data: []byte{0x30, 0x02, 0x00, 0x00}, ok: false},
		{name: "truncated", data: []byte{0x20, 0x02, 0x00}, ok: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, ok := ParseConnack(test.data)
			if code != test.code || ok != test.ok {
				t.Fatalf("ParseConnack(%#v) = (%#x, %t), want (%#x, %t)", test.data, code, ok, test.code, test.ok)
			}
		})
	}
}

type brokerResult struct {
	flags      byte
	clientID   string
	disconnect []byte
	err        error
}

func startBroker(t *testing.T, code byte, silent bool) (netip.Addr, int, <-chan brokerResult) {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	result := make(chan brokerResult, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			result <- brokerResult{err: err}
			return
		}
		defer connection.Close()
		if silent {
			_ = connection.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			_, _ = io.Copy(io.Discard, connection)
			result <- brokerResult{}
			return
		}

		header := make([]byte, 2)
		if _, err := io.ReadFull(connection, header); err != nil {
			result <- brokerResult{err: err}
			return
		}
		body := make([]byte, int(header[1]))
		if _, err := io.ReadFull(connection, body); err != nil {
			result <- brokerResult{err: err}
			return
		}
		if header[0] != 0x10 || len(body) < 12 {
			result <- brokerResult{err: io.ErrUnexpectedEOF}
			return
		}
		clientIDLength := int(binary.BigEndian.Uint16(body[10:12]))
		if clientIDLength < 0 || len(body) != 12+clientIDLength {
			result <- brokerResult{err: io.ErrUnexpectedEOF}
			return
		}
		clientID := string(body[12:])
		if !strings.HasPrefix(clientID, "mimamori-") || len(clientID) != len("mimamori-")+8 {
			result <- brokerResult{err: io.ErrUnexpectedEOF}
			return
		}
		if _, err := hex.DecodeString(clientID[len("mimamori-"):]); err != nil {
			result <- brokerResult{err: err}
			return
		}
		if _, err := connection.Write([]byte{0x20, 0x02, 0x00, code}); err != nil {
			result <- brokerResult{err: err}
			return
		}
		disconnect := make([]byte, 2)
		_ = connection.SetReadDeadline(time.Now().Add(time.Second))
		_, _ = io.ReadFull(connection, disconnect)
		result <- brokerResult{flags: body[7], clientID: clientID, disconnect: disconnect}
	}()
	address := listener.Addr().(*net.TCPAddr)
	return netip.MustParseAddr("127.0.0.1"), address.Port, result
}

func TestAnonymousAccepted(t *testing.T) {
	ip, port, resultCh := startBroker(t, 0x00, false)
	if !Anonymous(context.Background(), ip, port, time.Second) {
		t.Fatal("Anonymous() = false, want true")
	}
	result := <-resultCh
	if result.err != nil {
		t.Fatal(result.err)
	}
	if result.flags != 0x02 {
		t.Fatalf("CONNECT flags = %#x, want %#x", result.flags, byte(0x02))
	}
	if len(result.disconnect) != 2 {
		t.Fatalf("DISCONNECT = %#v, want two bytes", result.disconnect)
	}
	if result.disconnect[0] != 0xe0 || result.disconnect[1] != 0x00 {
		t.Fatalf("DISCONNECT = %#v, want [224 0]", result.disconnect)
	}
}

func TestAnonymousNotAuthorized(t *testing.T) {
	ip, port, resultCh := startBroker(t, 0x05, false)
	if Anonymous(context.Background(), ip, port, time.Second) {
		t.Fatal("Anonymous() = true, want false")
	}
	result := <-resultCh
	if result.err != nil {
		t.Fatal(result.err)
	}
	if result.flags != 0x02 {
		t.Fatalf("CONNECT flags = %#x, want %#x", result.flags, byte(0x02))
	}
}

func TestAnonymousSilentTimesOut(t *testing.T) {
	ip, port, resultCh := startBroker(t, 0x00, true)
	start := time.Now()
	if Anonymous(context.Background(), ip, port, 100*time.Millisecond) {
		t.Fatal("Anonymous() = true, want false")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("Anonymous() took %s, want less than 500ms", elapsed)
	}
	<-resultCh
}
