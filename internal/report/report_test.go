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

func TestScoreClampsWithoutIntegerOverflow(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1

	if got := Score([]Finding{{Points: maxInt}, {Points: maxInt}}); got != 0 {
		t.Errorf("Score() with large deductions = %d, want 0", got)
	}
	if got := Score([]Finding{{Points: minInt}}); got != 100 {
		t.Errorf("Score() with a large negative deduction = %d, want 100", got)
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
		if got := Status(100, test.findings); got != test.want {
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
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.15z"},
		{ID: "b", Level: Caution, Points: 8, IP: "192.168.1.2"},
	}
	want := []Finding{
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21, Device: "A"},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21, Device: "Z"},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 5900},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.10"},
		{ID: "a", Level: Caution, Points: 8, IP: "192.168.1.15z"},
		{ID: "b", Level: Caution, Points: 8, IP: "192.168.1.2"},
		{ID: "z", Level: Caution, Points: 8, IP: "192.168.1.2", Port: 21},
	}

	for _, input := range [][]Finding{
		findings,
		{findings[6], findings[5], findings[4], findings[3], findings[2], findings[1], findings[0]},
	} {
		got := append([]Finding(nil), input...)
		Sort(got)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Sort() = %#v, want %#v", got, want)
		}
	}
}
