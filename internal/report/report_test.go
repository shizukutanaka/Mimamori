package report

import (
	"reflect"
	"testing"
)

func TestScoreClamps(t *testing.T) {
	if got := Score([]Finding{{Points: 125}}); got != 0 {
		t.Fatalf("Score() = %d, want 0", got)
	}
	if got := Score([]Finding{{Points: -5}}); got != 100 {
		t.Fatalf("Score() = %d, want 100", got)
	}
	if got := Score([]Finding{{Points: 10}}); got != 90 {
		t.Fatalf("Score() = %d, want 90", got)
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		findings []Finding
		want     string
	}{
		{nil, "safe"},
		{[]Finding{{Level: Caution}}, "caution"},
		{[]Finding{{Level: Danger}}, "danger"},
	}
	for _, test := range tests {
		if got := Status(test.findings); got != test.want {
			t.Errorf("Status() = %q, want %q", got, test.want)
		}
	}
}

func TestSort(t *testing.T) {
	findings := []Finding{
		{ID: "caution-low", Level: Caution, Points: 5, IP: "192.168.1.1"},
		{ID: "danger-later", Level: Danger, Points: 25, IP: "192.168.1.10"},
		{ID: "danger-first", Level: Danger, Points: 25, IP: "192.168.1.2"},
		{ID: "caution-high", Level: Caution, Points: 8, IP: "192.168.1.3"},
	}
	Sort(findings)
	want := []string{"danger-first", "danger-later", "caution-high", "caution-low"}
	for i := range findings {
		if findings[i].ID != want[i] {
			t.Fatalf("Sort() order = %#v", findings)
		}
	}
}

func TestSortIsStableAndTotal(t *testing.T) {
	findings := []Finding{
		{ID: "z", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.10"},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 5900},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21, Device: "Z"},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21, Device: "A"},
		{ID: "b", Level: Caution, Points: 8, IP: "192.168.1.2"},
	}
	want := []Finding{
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21, Device: "A"},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21, Device: "Z"},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 5900},
		{ID: "b", Level: Caution, Points: 8, IP: "192.168.1.2"},
		{ID: "z", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.10"},
	}

	for _, input := range [][]Finding{
		findings,
		{findings[5], findings[4], findings[3], findings[2], findings[1], findings[0]},
	} {
		got := append([]Finding(nil), input...)
		Sort(got)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Sort() = %#v, want %#v", got, want)
		}
	}
}
