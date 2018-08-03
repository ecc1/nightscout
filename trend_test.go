package nightscout

import (
	"testing"
	"time"
)

func TestTrend(t *testing.T) {
	cases := []struct {
		trend   string
		entries Entries
	}{
		{"DoubleUp", sgvEntries(126, 108, 93, 79)},
		{"SingleUp", sgvEntries(108, 93, 79, 77)},
		{"FortyFiveUp", sgvEntries(141, 135, 126, 121)},
		{"Flat", sgvEntries(102, 98, 97, 99)},
		{"FortyFiveDown", sgvEntries(82, 87, 93, 99)},
		{"SingleDown", sgvEntries(117, 117, 129, 147)},
		{"DoubleDown", sgvEntries(117, 129, 147, 164)},
	}
	for _, c := range cases {
		t.Run(c.trend, func(t *testing.T) {
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
