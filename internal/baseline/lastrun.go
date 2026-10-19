package baseline

import "time"

func (s Store) LastRun() (time.Time, bool) {
	if !s.LastCheck.IsZero() {
		return s.LastCheck, true
	}
	var last time.Time
	found := false
	for _, seen := range s.Devices {
		if !found || seen.LastSeen.After(last) {
			last = seen.LastSeen
			found = true
		}
	}
	return last, found
}
