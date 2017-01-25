package nightscout

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type Interval struct {
	Start  time.Time
	Finish time.Time
}

func Gaps(since time.Time, gapDuration time.Duration) ([]Interval, error) {
	now := time.Now()
	window := now.Sub(since)
	log.Printf("retrieving Nightscout records from last %v", window)
	dateString := since.Format(DateStringLayout)
	// 2 entries per minute should be plenty.
	count := 2 * int(window/time.Minute)
	rest := fmt.Sprintf("entries.json?find[dateString][$gte]=%s&count=%d", dateString, count)
	var entries []EntryTime
	// Suppress verbose output for this.
	v := Verbose()
	SetVerbose(false)
	err := Get(rest, &entries)
	SetVerbose(v)
	if err != nil {
		return nil, err
	}
	// Sort entries in reverse chronological order,
	// even though they're currently already returned that way.
	sort.Sort(mostRecentFirst(entries))
	// Use current time to end any ongoing gap.
	times := []time.Time{now}
	// Convert millisecond Date field to Unix time.
	for _, r := range entries {
		sec := r.Date / 1000
		nsec := (r.Date % 1000) * 1000000
		times = append(times, time.Unix(sec, nsec))
	}
	// Use cutoff time to precede any ongoing gap.
	times = append(times, since)
	log.Printf("looking for gaps in %d Nightscout records", len(times))
	return findGaps(times, gapDuration), nil
}

func findGaps(entries []time.Time, gapDuration time.Duration) []Interval {
	var gaps []Interval
	for i := 0; i < len(entries)-1; i++ {
		cur := entries[i]
		prev := entries[i+1]
		if prev.IsZero() || prev.Equal(time.Unix(0, 0)) {
			continue
		}
		delta := cur.Sub(prev)
		if delta >= gapDuration {
			gaps = append(gaps, Interval{Start: prev, Finish: cur})
		}
	}
	return gaps
}

// Implement sort.Interface for reverse chronological order.
type mostRecentFirst []EntryTime

func (v mostRecentFirst) Len() int           { return len(v) }
func (v mostRecentFirst) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v mostRecentFirst) Less(i, j int) bool { return v[i].Date > v[j].Date }
