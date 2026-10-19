package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func FuzzLoad(f *testing.F) {
	f.Add([]byte(`{"version":1,"devices":{"aa:bb:cc:dd:ee:ff":{"first_seen":"2025-01-02T03:04:05Z","last_seen":"2025-01-02T04:05:06Z","last_ip":"192.168.1.2","name":"Camera"},"11:22:33:44:55:66":{"first_seen":"2025-02-03T04:05:06Z","last_seen":"2025-02-03T05:06:07Z","last_ip":"192.168.1.3","name":"Laptop"}}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"devices":null}`))
	f.Add([]byte(`{"version":1,"devices":`))
	f.Add([]byte("{\"version\":1,\"devices\":{},\"extra\":\"\xff\"}"))
	f.Add([]byte(`{"version":999999999999999999999999999999999999999999,"devices":{}}`))
	f.Add([]byte(`{"version":1,"devices":{"aa:bb:cc:dd:ee:ff":{"first_seen":"not-a-time","last_seen":"2025-01-02T04:05:06Z","last_ip":"192.168.1.2"}}}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		dir := t.TempDir()
		path := filepath.Join(dir, "devices.json")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}

		store, _, err := Load(path)
		if err != nil {
			return
		}

		roundTripPath := filepath.Join(dir, "round-trip", "devices.json")
		if err := Save(roundTripPath, store); err != nil {
			t.Fatal(err)
		}
		roundTrip, _, err := Load(roundTripPath)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(store, roundTrip) {
			t.Fatalf("Load/Save round trip = %#v, want %#v", roundTrip, store)
		}
	})
}
