package baseline

import "time"

func (s *Store) Prune(now time.Time, maxAge time.Duration, prunable func(mac string) bool) int {
	if s == nil || prunable == nil {
		return 0
	}
	pruned := 0
	for mac, seen := range s.Devices {
		if prunable(mac) && now.Sub(seen.LastSeen) > maxAge {
			delete(s.Devices, mac)
			pruned++
		}
	}
	return pruned
}
