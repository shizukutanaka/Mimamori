//go:build !windows

package wifi

func currentWLAN() (Security, bool, error) {
	return Unknown, false, nil
}
