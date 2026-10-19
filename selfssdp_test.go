package main

import (
	"net/netip"
	"reflect"
	"testing"

	"github.com/shizukutanaka/Mimamori/internal/ssdp"
)

func TestWithoutLocalResponses(t *testing.T) {
	localIP := netip.MustParseAddr("192.168.1.20")
	firstRemote := ssdp.Response{IP: netip.MustParseAddr("192.168.1.1"), Server: "gateway"}
	self := ssdp.Response{IP: localIP, Server: "local media server"}
	lastRemote := ssdp.Response{IP: netip.MustParseAddr("192.168.1.30"), Server: "other device"}
	responses := []ssdp.Response{firstRemote, self, lastRemote}

	got := withoutLocalResponses(responses, map[netip.Addr]struct{}{localIP: {}})
	want := []ssdp.Response{firstRemote, lastRemote}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("withoutLocalResponses() = %#v, want %#v", got, want)
	}
}

func TestWithoutLocalResponsesNilMap(t *testing.T) {
	responses := []ssdp.Response{
		{IP: netip.MustParseAddr("192.168.1.20"), Server: "local media server"},
		{IP: netip.MustParseAddr("192.168.1.1"), Server: "gateway"},
	}

	got := withoutLocalResponses(responses, nil)
	if !reflect.DeepEqual(got, responses) {
		t.Fatalf("withoutLocalResponses() = %#v, want %#v", got, responses)
	}
	if len(got) > 0 && &got[0] != &responses[0] {
		t.Fatal("withoutLocalResponses() copied responses with a nil local map")
	}
}

func TestCollectHostsExcludesLocalSSDPResponse(t *testing.T) {
	selfIP := netip.MustParseAddr("192.168.1.20")
	gatewayIP := netip.MustParseAddr("192.168.1.1")
	responses := []ssdp.Response{
		{IP: selfIP, Server: "local media server"},
		{IP: gatewayIP, Server: "gateway"},
	}
	filtered := withoutLocalResponses(responses, map[netip.Addr]struct{}{selfIP: {}})

	hosts := collectHosts(nil, gatewayIP, filtered, nil)
	for _, host := range hosts {
		if host == selfIP {
			t.Fatalf("collectHosts() includes local IP %s: %v", selfIP, hosts)
		}
	}
	if len(hosts) != 1 || hosts[0] != gatewayIP {
		t.Fatalf("collectHosts() = %v, want only gateway %s", hosts, gatewayIP)
	}
}
