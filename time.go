package nightscout

import (
	"time"
)

func msecsToTime(n int64) time.Time {
	sec := n / 1000
	nsec := (n % 1000) * 1000000
	return time.Unix(sec, nsec)
}

// Time returns the time.Time value corresponding to the Date field.
func (e Entry) Time() time.Time {
	return msecsToTime(e.Date)
}

// Time returns the time.Time value corresponding to the Date field.
func (e EntryTime) Time() time.Time {
	return msecsToTime(e.Date)
}

// Date converts a time.Time value to a Nightscout date value in milliseconds.
func Date(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

type (
	// EntryTime is used to unmarshal just the Date field of an Entry.
	EntryTime struct {
		Date int64 `json:"date"` // Unix time in milliseconds
	}

	// EntryTimes represents a sequence of entry times.
	EntryTimes []EntryTime
)

// Implement sort.Interface for reverse chronological order.

func (v EntryTimes) Len() int {
	return len(v)
}

func (v EntryTimes) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v EntryTimes) Less(i, j int) bool {
	return v[i].Date > v[j].Date
}
