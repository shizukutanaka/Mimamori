package baseline

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ErrInvalidJSON = errors.New("invalid device history JSON")

var syncDirFunc = syncDir

type Seen struct {
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	LastIP    string    `json:"last_ip"`
	Name      string    `json:"name,omitempty"`
}

type Store struct {
	Version   int             `json:"version"`
	Devices   map[string]Seen `json:"devices"`
	LastCheck time.Time       `json:"last_check"`
}

func emptyStore() Store {
	return Store{Version: 1, Devices: make(map[string]Seen)}
}

func Load(path string) (Store, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return emptyStore(), false, nil
	}
	if err != nil {
		return Store{}, false, err
	}
	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return Store{}, true, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if store.Devices == nil {
		store.Devices = make(map[string]Seen)
	}
	if store.Version == 0 {
		store.Version = 1
	}
	return store, true, nil
}

func (s *Store) Observe(mac, ip, name string, now time.Time) (isNew bool) {
	if s.Devices == nil {
		s.Devices = make(map[string]Seen)
	}
	key := strings.ToLower(strings.TrimSpace(mac))
	seen, exists := s.Devices[key]
	if !exists {
		s.Devices[key] = Seen{
			FirstSeen: now,
			LastSeen:  now,
			LastIP:    ip,
			Name:      name,
		}
		return true
	}
	seen.LastSeen = now
	seen.LastIP = ip
	if name != "" {
		seen.Name = name
	}
	s.Devices[key] = seen
	return false
}

func Save(path string, store Store) error {
	directory := filepath.Dir(path)
	var createdDirectories []string
	for current := directory; current != filepath.Dir(current); current = filepath.Dir(current) {
		_, err := os.Stat(current)
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			break
		}
		createdDirectories = append(createdDirectories, current)
	}
	if err := os.MkdirAll(directory, 0700); err != nil {
		return err
	}
	for _, created := range createdDirectories {
		syncDirFunc(filepath.Dir(created))
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(directory, ".devices-*.tmp")
	if err != nil {
		return err
	}
	tempPath := file.Name()
	defer os.Remove(tempPath)
	if err := file.Chmod(0600); err != nil {
		_ = file.Close()
		return err
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := renameReplace(tempPath, path); err != nil {
		return err
	}
	syncDirFunc(directory)
	return nil
}

func Quarantine(path string) (string, error) {
	quarantinePath := path + ".corrupt"
	for index := 0; index <= 99; index++ {
		if index > 0 {
			quarantinePath = fmt.Sprintf("%s.corrupt.%d", path, index)
		}
		_, err := os.Lstat(quarantinePath)
		if err == nil {
			continue
		}
		if !errors.Is(err, os.ErrNotExist) {
			return quarantinePath, err
		}
		if err := os.Rename(path, quarantinePath); err != nil {
			return quarantinePath, err
		}
		return quarantinePath, nil
	}
	return quarantinePath, fmt.Errorf("all quarantine backup names are in use")
}

func DefaultPath() (string, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDirectory, "Mimamori", "devices.json"), nil
}
