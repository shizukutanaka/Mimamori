package probe

import (
	"context"
	"net"
	"net/netip"
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
	closedPort := port + 1

	got := Scan(context.Background(), []netip.Addr{netip.MustParseAddr("127.0.0.1")}, []int{closedPort, port}, 300*time.Millisecond, 2)
	if len(got) != 1 || got[0].Port != port {
		t.Fatalf("Scan() = %#v, want only port %d", got, port)
	}
}
