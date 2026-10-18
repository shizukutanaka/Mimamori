package netbios

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"net"
	"net/netip"
	"strings"
	"time"
)

const encodedWildcardName = "CKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

func ParseNodeStatus(b []byte) (string, bool) {
	if !validNodeStatusHeader(b) {
		return "", false
	}

	offset := 12
	questionCount := int(binary.BigEndian.Uint16(b[4:6]))
	for i := 0; i < questionCount; i++ {
		next, ok := skipDNSName(b, offset)
		if !ok || next > len(b)-4 {
			return "", false
		}
		offset = next + 4
	}
	offset, ok := skipDNSName(b, offset)
	if !ok || offset > len(b)-10 {
		return "", false
	}
	recordType := binary.BigEndian.Uint16(b[offset : offset+2])
	class := binary.BigEndian.Uint16(b[offset+2 : offset+4])
	if recordType != 0x0021 || class != 1 {
		return "", false
	}
	dataLength := int(binary.BigEndian.Uint16(b[offset+8 : offset+10]))
	dataStart := offset + 10
	if dataLength < 1 || dataLength > len(b)-dataStart {
		return "", false
	}

	data := b[dataStart : dataStart+dataLength]
	nameCount := int(data[0])
	if nameCount > (len(data)-1)/18 {
		return "", false
	}
	for i := 0; i < nameCount; i++ {
		entry := data[1+i*18 : 1+(i+1)*18]
		if entry[15] != 0x00 || binary.BigEndian.Uint16(entry[16:18])&0x8000 != 0 {
			continue
		}
		name := strings.TrimRight(string(entry[:15]), " \x00")
		if name == "" || !printableASCII(name) {
			continue
		}
		return name, true
	}
	return "", false
}

func Names(ctx context.Context, hosts []netip.Addr, timeout time.Duration) map[netip.Addr]string {
	return names(ctx, hosts, 137, func(host netip.Addr) bool {
		return host.Is4() && host.IsPrivate()
	}, timeout)
}

func names(
	ctx context.Context,
	hosts []netip.Addr,
	port int,
	accept func(netip.Addr) bool,
	timeout time.Duration,
) map[netip.Addr]string {
	result := make(map[netip.Addr]string)
	if ctx == nil || ctx.Err() != nil || port < 1 || port > 65535 || timeout <= 0 || accept == nil {
		return result
	}

	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return result
	}
	defer conn.Close()

	deadline := time.Now().Add(timeout)
	queries := make(map[netip.Addr]uint16)
	for _, host := range hosts {
		if !host.Is4() || !accept(host) {
			continue
		}
		if _, exists := queries[host]; exists {
			continue
		}

		packet, ok := nodeStatusRequest()
		if !ok {
			continue
		}
		transactionID := binary.BigEndian.Uint16(packet[:2])
		destination := net.UDPAddrFromAddrPort(netip.AddrPortFrom(host, uint16(port)))
		if _, err := conn.WriteTo(packet, destination); err == nil {
			queries[host] = transactionID
		}
	}
	if len(queries) == 0 {
		return result
	}

	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := conn.SetReadDeadline(deadline); err != nil {
		return result
	}
	stopCancel := context.AfterFunc(ctx, func() {
		_ = conn.SetReadDeadline(time.Now())
	})
	defer stopCancel()

	answered := make(map[netip.Addr]struct{}, len(queries))
	packet := make([]byte, 65535)
	for {
		n, source, err := conn.ReadFrom(packet)
		if err != nil {
			return result
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
		transactionID, queried := queries[sourceIP]
		if !queried || n < 12 || binary.BigEndian.Uint16(packet[:2]) != transactionID {
			continue
		}
		if !validNodeStatusHeader(packet[:n]) {
			continue
		}

		answered[sourceIP] = struct{}{}
		if name, ok := ParseNodeStatus(packet[:n]); ok {
			if _, exists := result[sourceIP]; !exists {
				result[sourceIP] = name
			}
		}
		if len(answered) == len(queries) {
			return result
		}
	}
}

func nodeStatusRequest() ([]byte, bool) {
	packet := make([]byte, 12, 50)
	if _, err := rand.Read(packet[:2]); err != nil {
		return nil, false
	}
	binary.BigEndian.PutUint16(packet[4:6], 1)
	packet = append(packet, 0x20)
	packet = append(packet, encodedWildcardName...)
	packet = append(packet, 0, 0, 0x21, 0, 1)
	return packet, true
}

func validNodeStatusHeader(b []byte) bool {
	if len(b) < 12 {
		return false
	}
	flags := binary.BigEndian.Uint16(b[2:4])
	return flags&0x8000 != 0 &&
		flags&0x000f == 0 &&
		binary.BigEndian.Uint16(b[6:8]) >= 1
}

func skipDNSName(b []byte, offset int) (int, bool) {
	if offset < 0 || offset >= len(b) {
		return 0, false
	}
	for {
		if offset >= len(b) {
			return 0, false
		}
		length := int(b[offset])
		switch length & 0xc0 {
		case 0xc0:
			if offset+1 >= len(b) {
				return 0, false
			}
			pointer := (length&0x3f)<<8 | int(b[offset+1])
			if pointer >= len(b) {
				return 0, false
			}
			return offset + 2, true
		case 0:
			offset++
			if length == 0 {
				return offset, true
			}
			if offset+length > len(b) {
				return 0, false
			}
			offset += length
		default:
			return 0, false
		}
	}
}

func printableASCII(name string) bool {
	for i := 0; i < len(name); i++ {
		if name[i] < 0x20 || name[i] > 0x7e {
			return false
		}
	}
	return true
}
