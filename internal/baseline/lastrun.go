package baseline

import "time"

func (s Store) LastRun() (time.Time, bool) {
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
