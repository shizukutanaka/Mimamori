package report

import (
	"net/netip"
	"sort"
)

type Level string

const (
	Caution Level = "caution"
	Danger  Level = "danger"
)

type Finding struct {
	ID      string `json:"id"`
	Level   Level  `json:"level"`
	Points  int    `json:"points"`
	IP      string `json:"ip,omitempty"`
	Device  string `json:"device,omitempty"`
	Port    int    `json:"port,omitempty"`
	Message string `json:"message"`
	Fix     string `json:"fix"`
}

func Score(findings []Finding) int {
	score := 100
	for _, finding := range findings {
		score -= finding.Points
	}
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func Status(findings []Finding) string {
	for _, finding := range findings {
		if finding.Level == Danger {
			return "danger"
		}
	}
	if len(findings) > 0 {
		return "caution"
	}
	return "safe"
}

func Sort(findings []Finding) {
	sort.SliceStable(findings, func(i, j int) bool {
		left, right := findings[i], findings[j]
		if left.Level != right.Level {
			return left.Level == Danger
		}
		if left.Points != right.Points {
			return left.Points > right.Points
		}
		leftIP, leftErr := netip.ParseAddr(left.IP)
		rightIP, rightErr := netip.ParseAddr(right.IP)
		if leftErr == nil && rightErr == nil && leftIP != rightIP {
			return leftIP.Less(rightIP)
		}
		if left.IP != right.IP {
			return left.IP < right.IP
		}
		if left.ID != right.ID {
			return left.ID < right.ID
		}
		if left.Port != right.Port {
			return left.Port < right.Port
		}
		return left.Device < right.Device
	})
}
