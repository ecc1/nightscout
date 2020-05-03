package nightscout

import (
	"testing"
	"time"
)

func TestTrend(t *testing.T) {
	cases := []struct {
		entries Entries
		slope   float64
		trend   string
	}{
		{sgvEntries(126, 108, 93, 79), 3.12, "DoubleUp"},
		{sgvEntries(108, 93, 79, 77), 2.14, "SingleUp"},
		{sgvEntries(141, 135, 126, 121), 1.38, "FortyFiveUp"},
		{sgvEntries(102, 98, 97, 99), 0.2, "Flat"},
		{sgvEntries(82, 87, 93, 99), -1.14, "FortyFiveDown"},
		{sgvEntries(117, 117, 129, 147), -2.04, "SingleDown"},
		{sgvEntries(117, 129, 147, 164), -3.18, "DoubleDown"},
	}
	for _, c := range cases {
		t.Run(c.trend, func(t *testing.T) {
			slope := FindLine(c.entries).Slope
			if !closeEnough(slope, c.slope) {
				t.Errorf("Slope == %v, want %v", slope, c.slope)
			}
			trend := Trend(c.entries)
			if trend != c.trend {
				t.Errorf("Trend == %s, want %s", trend, c.trend)
			}
		})
	}
}

func sgvEntry(t time.Time, bg int) Entry {
	return Entry{
		Date:       Date(t),
		DateString: t.Format(DateStringLayout),
		Device:     Device(),
		Type:       SGVType,
		SGV:        bg,
	}
}

func sgvEntries(bgs ...int) Entries {
	t := parseTime("2018-06-30 12:00")
	entries := make(Entries, len(bgs))
	for i, bg := range bgs {
		entries[i] = sgvEntry(t, bg)
		t = t.Add(-5 * time.Minute)
	}
	return entries
}
