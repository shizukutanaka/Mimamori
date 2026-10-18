package mdns

import (
	"context"
	"encoding/binary"
	"net"
	"net/netip"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParsePTRResponse(t *testing.T) {
	tests := []struct {
		name     string
		packet   []byte
		question string
		target   string
		ok       bool
	}{
		{
			name:     "compressed answer",
			packet:   makePTRResponse("1.0.0.127.in-addr.arpa", "1.0.0.127.in-addr.arpa", "phone.local", 12, 1, 0x8000),
			question: "1.0.0.127.in-addr.arpa",
			target:   "phone.local",
			ok:       true,
		},
		{
			name:   "different owner",
			packet: makePTRResponse("1.0.0.127.in-addr.arpa", "other.local", "phone.local", 12, 1, 0x8000),
		},
		{
			name:   "non PTR answer",
			packet: makePTRResponse("1.0.0.127.in-addr.arpa", "1.0.0.127.in-addr.arpa", "phone.local", 1, 1, 0x8000),
		},
		{
			name:   "QR unset",
			packet: makePTRResponse("1.0.0.127.in-addr.arpa", "1.0.0.127.in-addr.arpa", "phone.local", 12, 1, 0),
		},
		{
			name:   "nonzero RCODE",
			packet: makePTRResponse("1.0.0.127.in-addr.arpa", "1.0.0.127.in-addr.arpa", "phone.local", 12, 1, 0x8003),
		},
		{name: "truncated header", packet: make([]byte, 11)},
		{name: "truncated name", packet: truncatedQuestion()},
		{name: "truncated rdata", packet: truncatedRData()},
		{name: "pointer loop", packet: pointerLoop()},
		{
			name:     "uppercase names are lowercased",
			packet:   makePTRResponse("1.0.0.127.IN-ADDR.ARPA", "1.0.0.127.IN-ADDR.ARPA", "Phone.LOCAL", 12, 0x8001, 0x8000),
			question: "1.0.0.127.in-addr.arpa",
			target:   "phone.local",
			ok:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			question, target, ok := ParsePTRResponse(test.packet)
			if ok != test.ok {
				t.Fatalf("ParsePTRResponse() ok = %t, want %t", ok, test.ok)
			}
			if question != test.question || target != test.target {
				t.Fatalf("ParsePTRResponse() = (%q, %q, %t), want (%q, %q, %t)",
					question, target, ok, test.question, test.target, test.ok)
			}
		})
	}
}

func FuzzParsePTRResponse(f *testing.F) {
	f.Add(makePTRResponse("1.0.0.127.in-addr.arpa", "1.0.0.127.in-addr.arpa", "phone.local", 12, 1, 0x8000))
	f.Add([]byte{})
	f.Add(pointerLoop())
	f.Fuzz(func(t *testing.T, packet []byte) {
		ParsePTRResponse(packet)
	})
}

func TestReverseNames(t *testing.T) {
	host := netip.MustParseAddr("127.0.0.1")
	question, _ := reverseQuestion(host)

	t.Run("strips local suffix", func(t *testing.T) {
		port := startResponder(t, nil, func(query []byte) []byte {
			response := makePTRResponse(question, question, "TaroS-iPhone.local", 12, 1, 0x8000)
			copy(response[:2], query[:2])
			return response
		})
		names := reverseNames(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, time.Second)
		if got := names[host]; got != "taros-iphone" {
			t.Fatalf("ReverseNames() name = %q, want %q", got, "taros-iphone")
		}
	})

	t.Run("ignores a reply from an unqueried IP", func(t *testing.T) {
		port := startResponder(t, nonLoopbackIPv4(t), func(query []byte) []byte {
			response := makePTRResponse(question, question, "phone.local", 12, 1, 0x8000)
			copy(response[:2], query[:2])
			return response
		})
		names := reverseNames(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, time.Second)
		if len(names) != 0 {
			t.Fatalf("ReverseNames() = %v, want no names", names)
		}
	})

	t.Run("ignores a mismatched question", func(t *testing.T) {
		port := startResponder(t, nil, func(query []byte) []byte {
			response := makePTRResponse("wrong.example", "wrong.example", "phone.local", 12, 1, 0x8000)
			copy(response[:2], query[:2])
			return response
		})
		names := reverseNames(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, time.Second)
		if len(names) != 0 {
			t.Fatalf("ReverseNames() = %v, want no names", names)
		}
	})

	t.Run("no reply", func(t *testing.T) {
		port := startResponder(t, nil, nil)
		start := time.Now()
		names := reverseNames(context.Background(), []netip.Addr{host}, port, netip.Addr.IsLoopback, 100*time.Millisecond)
		if len(names) != 0 {
			t.Fatalf("ReverseNames() = %v, want no names", names)
		}
		if elapsed := time.Since(start); elapsed > time.Second {
			t.Fatalf("ReverseNames() took %v without a reply", elapsed)
		}
	})
}

func makePTRResponse(question, owner, target string, recordType, class, flags uint16) []byte {
	packet := make([]byte, 12)
	binary.BigEndian.PutUint16(packet[2:4], flags)
	binary.BigEndian.PutUint16(packet[4:6], 1)
	binary.BigEndian.PutUint16(packet[6:8], 1)
	packet = appendDNSName(packet, question)
	packet = append(packet, 0, 12, 0, 1)
	if owner == question {
		packet = append(packet, 0xc0, 0x0c)
	} else {
		packet = appendDNSName(packet, owner)
	}
	packet = append(packet,
		byte(recordType>>8), byte(recordType),
		byte(class>>8), byte(class),
		0, 0, 0, 0,
		0, 0,
	)
	dataStart := len(packet)
	packet = appendDNSName(packet, target)
	binary.BigEndian.PutUint16(packet[dataStart-2:dataStart], uint16(len(packet)-dataStart))
	return packet
}

func appendDNSName(packet []byte, name string) []byte {
	if name == "" {
		return append(packet, 0)
	}
	for _, label := range strings.Split(name, ".") {
		packet = append(packet, byte(len(label)))
		packet = append(packet, label...)
	}
	return append(packet, 0)
}

func truncatedQuestion() []byte {
	packet := make([]byte, 12)
	binary.BigEndian.PutUint16(packet[4:6], 1)
	return append(packet, 5, 'a', 'b')
}

func truncatedRData() []byte {
	packet := make([]byte, 12)
	binary.BigEndian.PutUint16(packet[2:4], 0x8000)
	binary.BigEndian.PutUint16(packet[4:6], 1)
	binary.BigEndian.PutUint16(packet[6:8], 1)
	packet = appendDNSName(packet, "1.0.0.127.in-addr.arpa")
	packet = append(packet, 0, 12, 0, 1, 0xc0, 0x0c)
	packet = append(packet, 0, 12, 0, 1, 0, 0, 0, 0, 0, 4, 0xc0, 0x0c)
	return packet
}

func pointerLoop() []byte {
	packet := make([]byte, 12)
	binary.BigEndian.PutUint16(packet[4:6], 1)
	return append(packet, 0xc0, 0x0c)
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
		query := make([]byte, 2048)
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

func TestReverseQuestion(t *testing.T) {
	host := netip.MustParseAddr("192.168.1.23")
	got, ok := reverseQuestion(host)
	if !ok || got != "23.1.168.192.in-addr.arpa" {
		t.Fatalf("reverseQuestion() = (%q, %t)", got, ok)
	}
	if _, ok := reverseQuestion(netip.MustParseAddr("2001:db8::1")); ok {
		t.Fatal("reverseQuestion() accepted IPv6")
	}
}

func TestPTRQueryFormat(t *testing.T) {
	host := netip.MustParseAddr("192.168.1.23")
	packet, question, ok := ptrQuery(host)
	if !ok {
		t.Fatal("ptrQuery() failed")
	}
	if question != "23.1.168.192.in-addr.arpa" {
		t.Fatalf("ptrQuery() question = %q", question)
	}
	if binary.BigEndian.Uint16(packet[2:4]) != 0 || binary.BigEndian.Uint16(packet[4:6]) != 1 {
		t.Fatalf("ptrQuery() flags/count = %x/%d, want 0/1", packet[2:4], binary.BigEndian.Uint16(packet[4:6]))
	}
	decoded, next, valid := parseDNSName(packet, 12)
	if !valid || decoded != question || next+4 != len(packet) {
		t.Fatalf("ptrQuery() question encoding = %q (valid %t), next = %d", decoded, valid, next)
	}
	if binary.BigEndian.Uint16(packet[next:next+2]) != 12 || binary.BigEndian.Uint16(packet[next+2:next+4]) != 0x8001 {
		t.Fatalf("ptrQuery() type/class = %x/%x, want 12/0x8001", packet[next:next+2], packet[next+2:next+4])
	}
}

func TestPTRQuestionCompressionDoesNotMutateInput(t *testing.T) {
	packet := makePTRResponse("1.0.0.127.in-addr.arpa", "1.0.0.127.in-addr.arpa", "phone.local", 12, 1, 0x8000)
	before := append([]byte(nil), packet...)
	ParsePTRResponse(packet)
	if !reflect.DeepEqual(packet, before) {
		t.Fatal("ParsePTRResponse() mutated the packet")
	}
}
