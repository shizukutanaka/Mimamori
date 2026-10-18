package mdns

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

func ParsePTRResponse(b []byte) (question string, target string, ok bool) {
	if len(b) < 12 {
		return "", "", false
	}

	flags := binary.BigEndian.Uint16(b[2:4])
	if flags&0x8000 == 0 || flags&0x000f != 0 {
		return "", "", false
	}

	questionCount := int(binary.BigEndian.Uint16(b[4:6]))
	answerCount := int(binary.BigEndian.Uint16(b[6:8]))
	if questionCount == 0 {
		return "", "", false
	}

	offset := 12
	for i := 0; i < questionCount; i++ {
		name, next, valid := parseDNSName(b, offset)
		if !valid || next > len(b)-4 {
			return "", "", false
		}
		if i == 0 {
			question = name
		}
		offset = next + 4
	}

	for i := 0; i < answerCount; i++ {
		owner, next, valid := parseDNSName(b, offset)
		if !valid || next > len(b)-10 {
			return "", "", false
		}

		recordType := binary.BigEndian.Uint16(b[next : next+2])
		class := binary.BigEndian.Uint16(b[next+2 : next+4])
		dataLength := int(binary.BigEndian.Uint16(b[next+8 : next+10]))
		dataStart := next + 10
		if dataLength > len(b)-dataStart {
			return "", "", false
		}
		dataEnd := dataStart + dataLength

		if recordType == 12 && class&0x7fff == 1 && owner == question {
			name, end, valid := parseDNSName(b, dataStart)
			if !valid || end != dataEnd {
				return "", "", false
			}
			return question, name, true
		}
		offset = dataEnd
	}

	return "", "", false
}

func ReverseNames(ctx context.Context, hosts []netip.Addr, timeout time.Duration) map[netip.Addr]string {
	return reverseNames(ctx, hosts, 5353, func(host netip.Addr) bool {
		return host.Is4() && host.IsPrivate()
	}, timeout)
}

func reverseNames(
	ctx context.Context,
	hosts []netip.Addr,
	port int,
	accept func(netip.Addr) bool,
	timeout time.Duration,
) map[netip.Addr]string {
	names := make(map[netip.Addr]string)
	if ctx == nil || ctx.Err() != nil || port < 1 || port > 65535 || timeout <= 0 || accept == nil {
		return names
	}

	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return names
	}
	defer conn.Close()

	queries := make(map[netip.Addr]string)
	for _, host := range hosts {
		if !host.Is4() || !accept(host) {
			continue
		}
		if _, exists := queries[host]; exists {
			continue
		}

		packet, question, ok := ptrQuery(host)
		if !ok {
			continue
		}
		destination := net.UDPAddrFromAddrPort(netip.AddrPortFrom(host, uint16(port)))
		if _, err := conn.WriteTo(packet, destination); err == nil {
			queries[host] = question
		}
	}
	if len(queries) == 0 {
		return names
	}

	deadline := time.Now().Add(timeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := conn.SetReadDeadline(deadline); err != nil {
		return names
	}
	stopCancel := context.AfterFunc(ctx, func() {
		_ = conn.SetReadDeadline(time.Now())
	})
	defer stopCancel()

	packet := make([]byte, 65535)
	for {
		n, source, err := conn.ReadFrom(packet)
		if err != nil {
			return names
		}
		udpSource, ok := source.(*net.UDPAddr)
		if !ok {
			continue
		}
		sourceIP, ok := netip.AddrFromSlice(udpSource.IP)
		if !ok {
			continue
		}
		sourceIP = sourceIP.Unmap()
		expectedQuestion, queried := queries[sourceIP]
		if !queried {
			continue
		}
		responseQuestion, target, valid := ParsePTRResponse(packet[:n])
		if !valid || responseQuestion != expectedQuestion {
			continue
		}
		target = strings.TrimSuffix(target, ".local")
		if target == "" {
			continue
		}
		if _, exists := names[sourceIP]; !exists {
			names[sourceIP] = target
		}
		if len(names) == len(queries) {
			return names
		}
	}
}

func ptrQuery(host netip.Addr) ([]byte, string, bool) {
	question, ok := reverseQuestion(host)
	if !ok {
		return nil, "", false
	}

	packet := make([]byte, 12, 12+len(question)+6)
	if _, err := rand.Read(packet[:2]); err != nil {
		return nil, "", false
	}
	binary.BigEndian.PutUint16(packet[4:6], 1)
	for _, label := range strings.Split(question, ".") {
		packet = append(packet, byte(len(label)))
		packet = append(packet, label...)
	}
	packet = append(packet, 0, 0, 12, 0x80, 1)
	return packet, question, true
}

func reverseQuestion(host netip.Addr) (string, bool) {
	if !host.Is4() {
		return "", false
	}
	octets := host.As4()
	return strconv.Itoa(int(octets[3])) + "." +
		strconv.Itoa(int(octets[2])) + "." +
		strconv.Itoa(int(octets[1])) + "." +
		strconv.Itoa(int(octets[0])) + ".in-addr.arpa", true
}

func parseDNSName(packet []byte, start int) (string, int, bool) {
	if start < 0 || start >= len(packet) {
		return "", 0, false
	}

	position := start
	next := -1
	jumps := 0
	visited := make(map[int]struct{})
	labels := make([]string, 0, 4)
	wireLength := 1
	for {
		if position >= len(packet) {
			return "", 0, false
		}
		length := int(packet[position])
		switch length & 0xc0 {
		case 0xc0:
			if position+1 >= len(packet) {
				return "", 0, false
			}
			pointer := (length&0x3f)<<8 | int(packet[position+1])
			if pointer >= len(packet) {
				return "", 0, false
			}
			jumps++
			if jumps > 16 {
				return "", 0, false
			}
			if _, exists := visited[pointer]; exists {
				return "", 0, false
			}
			visited[pointer] = struct{}{}
			if next == -1 {
				next = position + 2
			}
			position = pointer
		case 0:
			position++
			if length == 0 {
				if next == -1 {
					next = position
				}
				return strings.ToLower(strings.Join(labels, ".")), next, true
			}
			if position+length > len(packet) {
				return "", 0, false
			}
			wireLength += length + 1
			if wireLength > 255 {
				return "", 0, false
			}
			labels = append(labels, string(packet[position:position+length]))
			position += length
		default:
			return "", 0, false
		}
	}
}
