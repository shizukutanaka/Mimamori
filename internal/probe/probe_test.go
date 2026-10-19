package probe

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"strconv"
	"testing"
	"time"
)

func TestScanFindsOpenPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	closedListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	closedPort := closedListener.Addr().(*net.TCPAddr).Port
	if err := closedListener.Close(); err != nil {
		t.Fatal(err)
	}

	got := Scan(context.Background(), []netip.Addr{netip.MustParseAddr("127.0.0.1")}, []int{closedPort, port}, 300*time.Millisecond, 2)
	if len(got) != 1 || got[0].Port != port {
		t.Fatalf("Scan() = %#v, want only port %d", got, port)
	}
}

func TestScanRetriesTimedOutPortsOnAliveHost(t *testing.T) {
	host := netip.MustParseAddr("192.0.2.1")
	attempts := make(map[string]int)
	var peers []net.Conn
	defer func() { closePipePeers(peers) }()
	dial := func(_ context.Context, address string) (net.Conn, error) {
		attempts[address]++
		switch address {
		case scanAddress(host, 23):
			if attempts[address] == 1 {
				return nil, fakeTimeoutError{}
			}
			return pipeConnection(&peers), nil
		case scanAddress(host, 21):
			return nil, errors.New("refused")
		default:
			return nil, errors.New("unexpected address")
		}
	}

	got := scan(context.Background(), []netip.Addr{host}, []int{23, 21}, 700*time.Millisecond, 1, dial)
	if len(got) != 1 || got[0] != (Open{IP: host, Port: 23}) {
		t.Fatalf("scan() = %#v, want only %s:23", got, host)
	}
	assertAttemptCount(t, attempts, scanAddress(host, 23), 2)
	assertAttemptCount(t, attempts, scanAddress(host, 21), 1)
}

func TestScanDoesNotRetryHostWithOnlyTimeouts(t *testing.T) {
	host := netip.MustParseAddr("192.0.2.2")
	attempts := make(map[string]int)
	dial := func(_ context.Context, address string) (net.Conn, error) {
		attempts[address]++
		return nil, fakeTimeoutError{}
	}

	got := scan(context.Background(), []netip.Addr{host}, []int{23, 21}, 700*time.Millisecond, 1, dial)
	if len(got) != 0 {
		t.Fatalf("scan() = %#v, want no open ports", got)
	}
	assertAttemptCount(t, attempts, scanAddress(host, 23), 1)
	assertAttemptCount(t, attempts, scanAddress(host, 21), 1)
}

func TestScanRetriesTimeoutWhenAnotherPortIsOpen(t *testing.T) {
	host := netip.MustParseAddr("192.0.2.3")
	attempts := make(map[string]int)
	var peers []net.Conn
	defer func() { closePipePeers(peers) }()
	dial := func(_ context.Context, address string) (net.Conn, error) {
		attempts[address]++
		if address == scanAddress(host, 23) {
			return pipeConnection(&peers), nil
		}
		return nil, fakeTimeoutError{}
	}

	got := scan(context.Background(), []netip.Addr{host}, []int{23, 21}, 700*time.Millisecond, 1, dial)
	if len(got) != 1 || got[0] != (Open{IP: host, Port: 23}) {
		t.Fatalf("scan() = %#v, want only %s:23", got, host)
	}
	assertAttemptCount(t, attempts, scanAddress(host, 23), 1)
	assertAttemptCount(t, attempts, scanAddress(host, 21), 2)
}

func TestScanDoesNotRetryAfterContextCancellation(t *testing.T) {
	host := netip.MustParseAddr("192.0.2.4")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	attempts := make(map[string]int)
	dial := func(_ context.Context, address string) (net.Conn, error) {
		attempts[address]++
		if address == scanAddress(host, 23) {
			return nil, fakeTimeoutError{}
		}
		cancel()
		return nil, errors.New("refused")
	}

	got := scan(ctx, []netip.Addr{host}, []int{23, 21}, 700*time.Millisecond, 1, dial)
	if len(got) != 0 {
		t.Fatalf("scan() = %#v, want no open ports", got)
	}
	assertAttemptCount(t, attempts, scanAddress(host, 23), 1)
	assertAttemptCount(t, attempts, scanAddress(host, 21), 1)
}

type fakeTimeoutError struct{}

func (fakeTimeoutError) Error() string   { return "timeout" }
func (fakeTimeoutError) Timeout() bool   { return true }
func (fakeTimeoutError) Temporary() bool { return true }

var _ net.Error = fakeTimeoutError{}

func scanAddress(ip netip.Addr, port int) string {
	return net.JoinHostPort(ip.String(), strconv.Itoa(port))
}

func pipeConnection(peers *[]net.Conn) net.Conn {
	connection, peer := net.Pipe()
	*peers = append(*peers, peer)
	return connection
}

func closePipePeers(peers []net.Conn) {
	for _, peer := range peers {
		_ = peer.Close()
	}
}

func assertAttemptCount(t *testing.T, attempts map[string]int, address string, want int) {
	t.Helper()
	if got := attempts[address]; got != want {
		t.Errorf("attempts for %s = %d, want %d", address, got, want)
	}
}
