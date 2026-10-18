package netbios

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"net/netip"
	"testing"
	"time"
)

type testNameEntry struct {
	name   string
	suffix byte
	flags  uint16
}

func TestParseNodeStatus(t *testing.T) {
	entries := []testNameEntry{
		{name: "WORKGROUP", suffix: 0x00, flags: 0x8000},
		{name: "DESKTOP-AB12CD", suffix: 0x20},
		{name: "DESKTOP-AB12CD", suffix: 0x00},
	}
	tests := []struct {
		name string
		data []byte
		want string
		ok   bool
	}{
		{
			name: "selects unique workstation name",
			data: makeNodeStatusResponse(0x1234, 0x8000, 0x0021, 1, entries),
			want: "DESKTOP-AB12CD",
			ok:   true,
		},
		{
			name: "compressed answer owner",
			data: makeCompressedNodeStatusResponse(0x1234, 0x8000, 0x0021, 1, entries),
			want: "DESKTOP-AB12CD",
			ok:   true,
		},
		{
			name: "wrong type",
			data: makeNodeStatusResponse(0x1234, 0x8000, 0x0001, 1, entries),
		},
		{
			name: "wrong class",
			data: makeNodeStatusResponse(0x1234, 0x8000, 0x0021, 2, entries),
		},
		{
			name: "response bit missing",
			data: makeNodeStatusResponse(0x1234, 0, 0x0021, 1, entries),
		},
		{
			name: "nonzero RCODE",
			data: makeNodeStatusResponse(0x1234, 0x8003, 0x0021, 1, entries),
		},
		{
			name: "truncated RDLENGTH",
			data: truncatedRecordLength(makeNodeStatusResponse(0x1234, 0x8000, 0x0021, 1, entries)),
		},
		{
			name: "truncated entries",
			data: truncatedEntries(makeNodeStatusResponse(0x1234, 0x8000, 0x0021, 1, entries[:1])),
		},
		{
			name: "skips nonprintable name",
			data: makeNodeStatusResponse(0x1234, 0x8000, 0x0021, 1, []testNameEntry{
				{name: "BAD\x01NAME", suffix: 0},
				{name: "PRINTER01", suffix: 0},
			}),
			want: "PRINTER01",
			ok:   true,
		},
		{
			name: "no printable unique name",
			data: makeNodeStatusResponse(0x1234, 0x8000, 0x0021, 1, []testNameEntry{
				{name: "BAD\x01NAME", suffix: 0},
			}),
		},
		{
			name: "truncated answer name",
			data: func() []byte {
				packet := make([]byte, 14)
				binary.BigEndian.PutUint16(packet[2:4], 0x8000)
				binary.BigEndian.PutUint16(packet[6:8], 1)
				packet[12], packet[13] = 3, 'X'
				return packet
			}(),
		},
		{name: "truncated header", data: make([]byte, 11)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := ParseNodeStatus(test.data)
			if got != test.want || ok != test.ok {
				t.Fatalf("ParseNodeStatus() = (%q, %t), want (%q, %t)", got, ok, test.want, test.ok)
			}
		})
	}
}

func FuzzParseNodeStatus(f *testing.F) {
	f.Add(makeNodeStatusResponse(0x1234, 0x8000, 0x0021, 1, []testNameEntry{
		{name: "DESKTOP-AB12CD", suffix: 0},
	}))
	f.Add([]byte{})
	f.Add(make([]byte, 11))
	f.Fuzz(func(t *testing.T, packet []byte) {
		ParseNodeStatus(packet)
	})
}

func TestNodeStatusRequestBytes(t *testing.T) {
	packet, ok := nodeStatusRequest()
	if !ok {
		t.Fatal("nodeStatusRequest() failed")
	}
	packet[0], packet[1] = 0, 0

	want := []byte{
		0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0,
		0x20,
	}
	want = append(want, "CKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"...)
	want = append(want, 0, 0, 0x21, 0, 1)
	if !bytes.Equal(packet, want) {
		t.Fatalf("nodeStatusRequest() = %x, want %x", packet, want)
	}
}

func TestNames(t *testing.T) {
	host := netip.MustParseAddr("127.0.0.1")

	t.Run("correct reply", func(t *testing.T) {
		port := startResponder(t, nil, func(query []byte) []byte {
			return makeNodeStatusResponse(
				binary.BigEndian.Uint16(query[:2]),
				0x8000,
				0x0021,
				1,
				[]testNameEntry{{name: "DESKTOP-AB12CD", suffix: 0}},
			)
		})
		got := names(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, time.Second)
		if got[host] != "DESKTOP-AB12CD" {
			t.Fatalf("names() = %v, want %q", got, "DESKTOP-AB12CD")
		}
	})

	t.Run("reply window starts after all queries are sent", func(t *testing.T) {
		previousHook := afterNameQuerySend
		sends := 0
		afterNameQuerySend = func() {
			sends++
			if sends == 1 {
				time.Sleep(50 * time.Millisecond)
			}
		}
		t.Cleanup(func() { afterNameQuerySend = previousHook })

		port := startResponder(t, nil, func(query []byte) []byte {
			return makeNodeStatusResponse(
				binary.BigEndian.Uint16(query[:2]),
				0x8000,
				0x0021,
				1,
				[]testNameEntry{{name: "DESKTOP-AB12CD", suffix: 0}},
			)
		})
		hosts := []netip.Addr{host, netip.MustParseAddr("127.0.0.2")}
		got := names(context.Background(), hosts, port, netip.Addr.IsLoopback, 10*time.Millisecond)
		if got[host] != "DESKTOP-AB12CD" {
			t.Fatalf("names() = %v, want %q", got, "DESKTOP-AB12CD")
		}
		if sends != len(hosts) {
			t.Fatalf("sent %d queries, want %d", sends, len(hosts))
		}
	})

	t.Run("wrong transaction ID is ignored", func(t *testing.T) {
		port := startResponder(t, nil, func(query []byte) []byte {
			return makeNodeStatusResponse(
				binary.BigEndian.Uint16(query[:2])+1,
				0x8000,
				0x0021,
				1,
				[]testNameEntry{{name: "DESKTOP-AB12CD", suffix: 0}},
			)
		})
		got := names(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, 100*time.Millisecond)
		if len(got) != 0 {
			t.Fatalf("names() = %v, want no names", got)
		}
	})

	t.Run("reply from unqueried IP is ignored", func(t *testing.T) {
		port := startResponder(t, nonLoopbackIPv4(t), func(query []byte) []byte {
			return makeNodeStatusResponse(
				binary.BigEndian.Uint16(query[:2]),
				0x8000,
				0x0021,
				1,
				[]testNameEntry{{name: "DESKTOP-AB12CD", suffix: 0}},
			)
		})
		got := names(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, 100*time.Millisecond)
		if len(got) != 0 {
			t.Fatalf("names() = %v, want no names", got)
		}
	})

	t.Run("returns early after all hosts answer", func(t *testing.T) {
		port := startResponder(t, nil, func(query []byte) []byte {
			return makeNodeStatusResponse(
				binary.BigEndian.Uint16(query[:2]),
				0x8000,
				0x0021,
				1,
				[]testNameEntry{{name: "DESKTOP-AB12CD", suffix: 0}},
			)
		})
		start := time.Now()
		got := names(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, 3*time.Second)
		elapsed := time.Since(start)
		if got[host] != "DESKTOP-AB12CD" {
			t.Fatalf("names() = %v, want %q", got, "DESKTOP-AB12CD")
		}
		if elapsed >= time.Second {
			t.Fatalf("names() took %v, want early return well before the 3 second timeout", elapsed)
		}
	})
}

func makeNodeStatusResponse(id, flags, recordType, class uint16, entries []testNameEntry) []byte {
	packet := make([]byte, 12)
	binary.BigEndian.PutUint16(packet[:2], id)
	binary.BigEndian.PutUint16(packet[2:4], flags)
	binary.BigEndian.PutUint16(packet[6:8], 1)
	packet = append(packet, 1, '*', 0)
	return appendNodeStatusRecord(packet, recordType, class, entries)
}

func makeCompressedNodeStatusResponse(id, flags, recordType, class uint16, entries []testNameEntry) []byte {
	packet := make([]byte, 12)
	binary.BigEndian.PutUint16(packet[:2], id)
	binary.BigEndian.PutUint16(packet[2:4], flags)
	binary.BigEndian.PutUint16(packet[4:6], 1)
	binary.BigEndian.PutUint16(packet[6:8], 1)
	packet = append(packet, 1, '*', 0, 0, 0x21, 0, 1, 0xc0, 0x0c)
	return appendNodeStatusRecord(packet, recordType, class, entries)
}

func appendNodeStatusRecord(packet []byte, recordType, class uint16, entries []testNameEntry) []byte {
	packet = append(packet,
		byte(recordType>>8), byte(recordType),
		byte(class>>8), byte(class),
		0, 0, 0, 0,
		0, 0,
	)
	dataStart := len(packet)
	data := []byte{byte(len(entries))}
	for _, entry := range entries {
		name := make([]byte, 15)
		for i := range name {
			name[i] = ' '
		}
		copy(name, entry.name)
		data = append(data, name...)
		data = append(data, entry.suffix, byte(entry.flags>>8), byte(entry.flags))
	}
	packet = append(packet, data...)
	binary.BigEndian.PutUint16(packet[dataStart-2:dataStart], uint16(len(data)))
	return packet
}

func truncatedRecordLength(packet []byte) []byte {
	offset, ok := skipDNSName(packet, 12)
	if !ok {
		panic("test fixture has an invalid answer name")
	}
	binary.BigEndian.PutUint16(packet[offset+8:offset+10], 0xffff)
	return packet
}

func truncatedEntries(packet []byte) []byte {
	offset, ok := skipDNSName(packet, 12)
	if !ok {
		panic("test fixture has an invalid answer name")
	}
	packet[offset+10] = 2
	return packet
}

func startResponder(t *testing.T, sourceIP net.IP, reply func([]byte) []byte) int {
	t.Helper()
	server, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = server.Close() })

	var replyConn *net.UDPConn
	if sourceIP != nil {
		replyConn, err = net.ListenUDP("udp4", &net.UDPAddr{IP: sourceIP})
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = replyConn.Close() })
	}

	go func() {
		query := make([]byte, 1024)
		n, remote, err := server.ReadFromUDP(query)
		if err != nil || reply == nil {
			return
		}
		response := reply(append([]byte(nil), query[:n]...))
		writer := server
		if replyConn != nil {
			writer = replyConn
		}
		_, _ = writer.WriteToUDP(response, remote)
	}()
	return server.LocalAddr().(*net.UDPAddr).Port
}

func nonLoopbackIPv4(t *testing.T) net.IP {
	t.Helper()
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		t.Fatal(err)
	}
	for _, address := range addresses {
		var ip net.IP
		switch address := address.(type) {
		case *net.IPNet:
			ip = address.IP
		case *net.IPAddr:
			ip = address.IP
		}
		if ip4 := ip.To4(); ip4 != nil && !ip4.IsLoopback() {
			return append(net.IP(nil), ip4...)
		}
	}
	t.Skip("no non-loopback IPv4 address available for reply source")
	return nil
}
