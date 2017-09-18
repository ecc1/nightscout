package nightscout

import (
	"fmt"
	"testing"
	"time"
)

func parseTime(s string) time.Time {
	const layout = "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

const (
	numEntries    = 20
	entryInterval = 5 * time.Minute
	gapDuration   = 7 * time.Minute
)

var (
	T = make([]time.Time, numEntries)
	E = make(Entries, numEntries)
)

// Construct Entry and Time values in reverse chronological order.
func makeEntries() {
	t := parseTime("2017-10-01 01:00:00")
	for n := 0; n < numEntries; n++ {
		t = t.Add(-entryInterval)
		T[n] = t
		E[n].Date = Date(t)
	}
}

func init() {
	makeEntries()
}

func (e Entry) String() string {
	return e.Time().Format("3:04")
}

func (g Gap) String() string {
	s := g.Start.Format("3:04")
	f := g.Finish.Format("3:04")
	return fmt.Sprintf("%s–%s", s, f)
}

func TestFindGaps(t *testing.T) {
	cases := []struct {
		times []time.Time
		gaps  []Gap
	}{
		{[]time.Time{}, []Gap{}},
		{[]time.Time{T[0]}, []Gap{}},
		{[]time.Time{T[0], T[1]}, []Gap{}},
		{[]time.Time{T[0], T[2]}, []Gap{{Finish: T[0], Start: T[2]}}},
		{[]time.Time{T[18], T[19]}, []Gap{}},
		{[]time.Time{T[17], T[19]}, []Gap{{Finish: T[17], Start: T[19]}}},
		{[]time.Time{T[0], T[1], T[2], T[5], T[6], T[7], T[15], T[16], T[19]}, []Gap{{Finish: T[2], Start: T[5]}, {Finish: T[7], Start: T[15]}, {Finish: T[16], Start: T[19]}}},
	}
	for _, c := range cases {
		g := findGaps(c.times, gapDuration)
		if !equalGaps(g, c.gaps) {
			t.Errorf("findGaps(%v) == %+v, want %+v", c.times, g, c.gaps)
		}
	}
}

func equalGaps(x, y []Gap) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func TestMissing(t *testing.T) {
	cases := []struct {
		entries Entries
		gaps    []Gap
		missing Entries
	}{
		{E, []Gap{}, Entries{}},
		{E, []Gap{{Finish: T[0], Start: T[2]}}, Entries{E[1]}},
		{E, []Gap{{Finish: T[2], Start: T[5]}, {Finish: T[7], Start: T[9]}, {Finish: T[16], Start: T[19]}}, Entries{E[3], E[4], E[8], E[17], E[18]}},
	}
	for _, c := range cases {
		missing := Missing(c.entries, c.gaps)
		if !equalEntries(missing, c.missing) {
			t.Errorf("Missing(%v, %+v) == %v, want %v", c.entries, c.gaps, missing, c.missing)
		}
	}
}

func equalEntries(x, y Entries) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
