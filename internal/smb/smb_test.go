package smb

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

func TestSMB1EnabledResponses(t *testing.T) {
	valid := validNegotiateResponse()
	dialectUnsupported := append([]byte(nil), valid...)
	binary.LittleEndian.PutUint16(dialectUnsupported[33:35], 0xffff)
	smb2 := append([]byte(nil), valid...)
	copy(smb2[:4], []byte{0xfe, 'S', 'M', 'B'})
	statusError := append([]byte(nil), valid...)
	statusError[5] = 1

	tests := []struct {
		name     string
		response []byte
		want     bool
	}{
		{name: "SMB1 accepted", response: valid, want: true},
		{name: "unsupported dialect", response: dialectUnsupported},
		{name: "SMB2 reply", response: smb2},
		{name: "nonzero status", response: statusError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedRequest := expectedNegotiateRequest()
			port, done := startFakeServer(t, func(connection net.Conn) error {
				request, err := readRequest(connection)
				if err != nil {
					return err
				}
				if !bytes.Equal(request, expectedRequest) {
					return fmt.Errorf("request bytes = % x, want % x", request, expectedRequest)
				}
				return writeFrame(connection, test.response)
			})

			got := SMB1Enabled(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second)
			waitForServer(t, done)
			if got != test.want {
				t.Fatalf("SMB1Enabled() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestSMB1EnabledResetAndShortResponse(t *testing.T) {
	tests := []struct {
		name    string
		handler func(net.Conn) error
	}{
		{
			name: "reset",
			handler: func(connection net.Conn) error {
				if _, err := readRequest(connection); err != nil {
					return err
				}
				if tcp, ok := connection.(*net.TCPConn); ok {
					return tcp.SetLinger(0)
				}
				return nil
			},
		},
		{
			name: "short response",
			handler: func(connection net.Conn) error {
				if _, err := readRequest(connection); err != nil {
					return err
				}
				return writeFrameWithLength(connection, 37, []byte{0xff, 'S', 'M', 'B'})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port, done := startFakeServer(t, test.handler)
			if got := SMB1Enabled(context.Background(), netip.MustParseAddr("127.0.0.1"), port, time.Second); got {
				t.Fatal("SMB1Enabled() = true, want false")
			}
			waitForServer(t, done)
		})
	}
}

func TestSMB1EnabledTimeout(t *testing.T) {
	port, done := startFakeServer(t, func(connection net.Conn) error {
		if _, err := readRequest(connection); err != nil {
			return err
		}
		time.Sleep(150 * time.Millisecond)
		return nil
	})

	if got := SMB1Enabled(context.Background(), netip.MustParseAddr("127.0.0.1"), port, 40*time.Millisecond); got {
		t.Fatal("SMB1Enabled() = true, want false")
	}
	waitForServer(t, done)
}

func TestSMB1EnabledContextCancel(t *testing.T) {
	requestReceived := make(chan struct{})
	port, done := startFakeServer(t, func(connection net.Conn) error {
		if _, err := readRequest(connection); err != nil {
			return err
		}
		close(requestReceived)
		var response [1]byte
		_, _ = connection.Read(response[:])
		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan bool, 1)
	go func() {
		result <- SMB1Enabled(ctx, netip.MustParseAddr("127.0.0.1"), port, time.Second)
	}()

	select {
	case <-requestReceived:
	case <-time.After(time.Second):
		t.Fatal("server did not receive SMB negotiate request")
	}
	cancel()
	select {
	case got := <-result:
		if got {
			t.Fatal("SMB1Enabled() = true after context cancellation, want false")
		}
	case <-time.After(time.Second):
		t.Fatal("SMB1Enabled() did not stop after context cancellation")
	}
	waitForServer(t, done)
}

func TestParseNegotiateResponse(t *testing.T) {
	valid := validNegotiateResponse()
	invalidDialect := append([]byte(nil), valid...)
	binary.LittleEndian.PutUint16(invalidDialect[33:35], 1)
	invalidStatus := append([]byte(nil), valid...)
	invalidStatus[5] = 1
	invalidWordCount := append([]byte(nil), valid...)
	invalidWordCount[32] = 0
	tests := []struct {
		name    string
		payload []byte
		want    bool
	}{
		{name: "valid", payload: valid, want: true},
		{name: "short", payload: valid[:36]},
		{name: "SMB2", payload: []byte{0xfe, 'S', 'M', 'B'}},
		{name: "wrong command", payload: responseWithCommand(0x73)},
		{name: "nonzero status", payload: invalidStatus},
		{name: "no words", payload: invalidWordCount},
		{name: "nonzero dialect", payload: invalidDialect},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseNegotiateResponse(test.payload); got != test.want {
				t.Fatalf("ParseNegotiateResponse() = %t, want %t", got, test.want)
			}
		})
	}
}

func FuzzParseNegotiateResponse(f *testing.F) {
	f.Add(validNegotiateResponse())
	f.Add([]byte{0xff, 'S', 'M', 'B'})
	f.Add([]byte(nil))
	f.Fuzz(func(t *testing.T, payload []byte) {
		ParseNegotiateResponse(payload)
	})
}

func startFakeServer(t *testing.T, handler func(net.Conn) error) (int, <-chan error) {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer connection.Close()
		_ = connection.SetDeadline(time.Now().Add(3 * time.Second))
		done <- handler(connection)
	}()
	t.Cleanup(func() {
		_ = listener.Close()
	})
	return listener.Addr().(*net.TCPAddr).Port, done
}

func waitForServer(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("fake server did not finish")
	}
}

func readRequest(connection net.Conn) ([]byte, error) {
	var header [4]byte
	if _, err := io.ReadFull(connection, header[:]); err != nil {
		return nil, err
	}
	length := int(header[1])<<16 | int(header[2])<<8 | int(header[3])
	if length > 512 {
		return nil, fmt.Errorf("request payload length %d exceeds test limit", length)
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(connection, payload); err != nil {
		return nil, err
	}
	return append(header[:], payload...), nil
}

func writeFrame(connection net.Conn, payload []byte) error {
	return writeFrameWithLength(connection, len(payload), payload)
}

func writeFrameWithLength(connection net.Conn, advertisedLength int, payload []byte) error {
	frame := make([]byte, 4+len(payload))
	frame[1] = byte(advertisedLength >> 16)
	frame[2] = byte(advertisedLength >> 8)
	frame[3] = byte(advertisedLength)
	copy(frame[4:], payload)
	_, err := io.Copy(connection, bytes.NewReader(frame))
	return err
}

func expectedNegotiateRequest() []byte {
	return []byte{
		0x00, 0x00, 0x00, 0x2f,
		0xff, 'S', 'M', 'B', 0x72,
		0x00, 0x00, 0x00, 0x00,
		0x18, 0x01, 0x00,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x34, 0x12,
		0x00, 0x00,
		0x01, 0x00,
		0x00,
		0x0c, 0x00,
		0x02, 'N', 'T', ' ', 'L', 'M', ' ', '0', '.', '1', '2', 0x00,
	}
}

func validNegotiateResponse() []byte {
	return responseWithCommand(smbCommandNegotiate)
}

func responseWithCommand(command byte) []byte {
	payload := make([]byte, 37)
	copy(payload[:4], []byte{0xff, 'S', 'M', 'B'})
	payload[4] = command
	payload[32] = 1
	return payload
}
