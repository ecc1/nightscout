package nightscout

import (
	"fmt"
	"testing"
	"time"
)

const (
	gapDuration = 7 * time.Minute

	TestTimeLayout = "2006-01-02 15:04:05"
)

func parseTime(s string) time.Time {
	t, err := time.ParseInLocation(TestTimeLayout, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

var (
	T = []time.Time{
		parseTime("2017-10-01 01:35:00"),
		parseTime("2017-10-01 01:30:00"),
		parseTime("2017-10-01 01:25:00"),
		parseTime("2017-10-01 01:20:00"),
		parseTime("2017-10-01 01:15:00"),
		parseTime("2017-10-01 01:10:00"),
		parseTime("2017-10-01 01:05:00"),
		parseTime("2017-10-01 01:00:00"),
		parseTime("2017-10-01 00:55:00"),
		parseTime("2017-10-01 00:50:00"),
		parseTime("2017-10-01 00:45:00"),
		parseTime("2017-10-01 00:40:00"),
		parseTime("2017-10-01 00:35:00"),
		parseTime("2017-10-01 00:30:00"),
		parseTime("2017-10-01 00:25:00"),
		parseTime("2017-10-01 00:20:00"),
		parseTime("2017-10-01 00:15:00"),
		parseTime("2017-10-01 00:10:00"),
		parseTime("2017-10-01 00:05:00"),
		parseTime("2017-10-01 00:00:00"),
	}

	E = Entries{
		{Date: Date(T[0])},
		{Date: Date(T[1])},
		{Date: Date(T[2])},
		{Date: Date(T[3])},
		{Date: Date(T[4])},
		{Date: Date(T[5])},
		{Date: Date(T[6])},
		{Date: Date(T[7])},
		{Date: Date(T[8])},
		{Date: Date(T[9])},
		{Date: Date(T[10])},
		{Date: Date(T[11])},
		{Date: Date(T[12])},
		{Date: Date(T[13])},
		{Date: Date(T[14])},
		{Date: Date(T[15])},
		{Date: Date(T[16])},
		{Date: Date(T[17])},
		{Date: Date(T[18])},
		{Date: Date(T[19])},
	}
)

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
		{nil, nil},
		{[]time.Time{T[0]}, nil},
		{[]time.Time{T[0], T[1]}, nil},
		{[]time.Time{T[0], T[2]}, []Gap{{Finish: T[0], Start: T[2]}}},
		{[]time.Time{T[18], T[19]}, nil},
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
