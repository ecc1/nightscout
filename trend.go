package nightscout

import (
	"time"
)

const (
	trendEntries = 4
	trendWindow  = trendEntries * 5 * time.Minute
)

// Trend returns a string describing the glucose trend,
// assuming the entries are in reverse chronological order.
func Trend(entries Entries) string {
	cur := entries[0]
	if cur.Type != SGVType {
		return ""
	}
	history := getHistory(entries)
	if len(history) == 1 {
		return ""
	}
	slope := FindLine(history).Slope
	if slope > 3 {
		return "DoubleUp"
	}
	if slope > 2 {
		return "SingleUp"
	}
	if slope > 1 {
		return "FortyFiveUp"
	}
	if slope >= -1 {
		return "Flat"
	}
	if slope >= -2 {
		return "FortyFiveDown"
	}
	if slope >= -3 {
		return "SingleDown"
	}
	return "DoubleDown"
}

func getHistory(entries Entries) Entries {
	history := make(Entries, 0, trendEntries)
	history = append(history, entries[0])
	for _, e := range entries[1:] {
		if len(history) == trendEntries {
			break
		}
		if e.Type != SGVType {
			continue
		}
		d := history[len(history)-1]
		if d.Time().Sub(e.Time()) > trendWindow {
			break
		}
		history = append(history, e)
	}
	return history
}

// These functions implement the Points interface for FindLine.

// X returns the x-coordinate of entry i.
func (e Entries) X(i int) float64 {
	return float64(e[i].Time().UnixNano()) / float64(time.Minute)
}

// Y returns the y-coordinate of entry i.
func (e Entries) Y(i int) float64 {
	return float64(e[i].SGV)
}
