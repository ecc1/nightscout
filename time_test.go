package nightscout

import (
	"os"
	"testing"
	"time"
)

// Force timezone to match test data.
func init() {
	os.Setenv("TZ", "America/New_York")
}

func TestTime(t *testing.T) {
	cases := []struct {
		d int64
		t time.Time
	}{
		{1230786000000, parseTime("2009-01-01 00:00:00")},
		{1249871143000, parseTime("2009-08-09 22:25:43")},
		{1298091985000, parseTime("2011-02-19 00:06:25")},
		{1469042740000, parseTime("2016-07-20 15:25:40")},
		{1480132712000, parseTime("2016-11-25 22:58:32")},
		{1505625229000, parseTime("2017-09-17 01:13:49")},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tv := msecsToTime(c.d)
			if !tv.Equal(c.t) {
				t.Errorf("msecsToTime(%v) == %v, want %v", c.d, tv, c.t)
			}
			d := Date(c.t)
			if d != c.d {
				t.Errorf("Date(%v) == %v, want %v", c.t, d, c.d)
			}
		})
	}
}

var layouts = []string{
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
}

func parseTime(s string) time.Time {
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			return t
		}
	}
	panic(err)
}
