package report

import (
	"math/big"
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
	var total, points big.Int
	for _, finding := range findings {
		points.SetInt64(int64(finding.Points))
		total.Add(&total, &points)
	}
	if total.Cmp(big.NewInt(100)) >= 0 {
		return 0
	}
	if total.Sign() <= 0 {
		return 100
	}
	return 100 - int(total.Int64())
}

func Status(score int, findings []Finding) string {
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
			leftRank, rightRank := levelRank(left.Level), levelRank(right.Level)
			if leftRank != rightRank {
				return leftRank < rightRank
			}
			return left.Level < right.Level
		}
		if left.Points != right.Points {
			return left.Points > right.Points
		}
		if left.ID != right.ID {
			return left.ID < right.ID
		}
		if order := compareIP(left.IP, right.IP); order != 0 {
			return order < 0
		}
		if left.Port != right.Port {
			return left.Port < right.Port
		}
		if left.Device != right.Device {
			return left.Device < right.Device
		}
		if left.Message != right.Message {
			return left.Message < right.Message
		}
		return left.Fix < right.Fix
	})
}

func levelRank(level Level) int {
	switch level {
	case Danger:
		return 0
	case Caution:
		return 1
	default:
		return 2
	}
}

func compareIP(left, right string) int {
	leftAddr, leftErr := netip.ParseAddr(left)
	rightAddr, rightErr := netip.ParseAddr(right)
	if leftErr == nil && rightErr == nil {
		if leftAddr != rightAddr {
			if leftAddr.Less(rightAddr) {
				return -1
			}
			return 1
		}
	} else if leftErr == nil {
		return -1
	} else if rightErr == nil {
		return 1
	}
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}
