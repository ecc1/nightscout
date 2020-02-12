package nightscout

import (
	"fmt"
	"time"
)

func (w Website) XDripTime() (time.Time, error) {
	var p struct {
		Status []struct {
			Now int64 `json:"now"` // Unix time in milliseconds
		} `json:"status"`
	}
	err := w.Get("pebble", &p)
	if err != nil {
		return time.Time{}, err
	}
	if len(p.Status) != 1 {
		return time.Time{}, fmt.Errorf("unexpected Status length (%d) in xDrip pebble response", len(p.Status))
	}
	return msecsToTime(p.Status[0].Now), nil
}

func (w Website) XDripEntries() (Entries, error) {
	var entries Entries
	err := w.Get("sgv.json", &entries)
	return entries, err
}
