package nightscout

import (
	"fmt"
	"log"
	"sort"
	"time"
)

// Gap represents a gap.
type Gap struct {
	Start  time.Time
	Finish time.Time
}

// Gaps finds gaps in Nightscout entries since the given time
// that are longer than the specified duration.
func Gaps(since time.Time, gapDuration time.Duration) ([]Gap, error) {
	now := time.Now()
	window := now.Sub(since)
	log.Printf("retrieving Nightscout records from last %v", window)
	dateString := since.Format(DateStringLayout)
	// 2 entries per minute should be plenty.
	count := 2 * int(window/time.Minute)
	rest := fmt.Sprintf("entries.json?find[dateString][$gte]=%s&count=%d", dateString, count)
	var entries mostRecentFirst
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
	sort.Sort(entries)
	// Use current time to end any ongoing gap.
	times := []time.Time{now}
	// Convert Date fields to time.Time values.
	for _, e := range entries {
		times = append(times, e.Time())
	}
	// Use cutoff time to precede any ongoing gap.
	times = append(times, since)
	log.Printf("looking for gaps in %d Nightscout records", len(times))
	return findGaps(times, gapDuration), nil
}

func findGaps(entries []time.Time, gapDuration time.Duration) []Gap {
	var gaps []Gap
	for i := 0; i < len(entries)-1; i++ {
		cur := entries[i]
		prev := entries[i+1]
		if prev.IsZero() || prev.Equal(time.Unix(0, 0)) {
			continue
		}
		delta := cur.Sub(prev)
		if delta >= gapDuration {
			gaps = append(gaps, Gap{Start: prev, Finish: cur})
		}
	}
	return gaps
}

const (
	edgeMargin = 1 * time.Minute
)

// Missing returns the Entry values that fall within the given gaps.
func Missing(entries Entries, gaps []Gap) Entries {
	var missing Entries
	i := 0
	for _, g := range gaps {
		// Skip over entries that lie outside the gap.
		for i < len(entries) {
			t := entries[i].Time()
			if t.Before(g.Finish) {
				break
			}
			i++
		}
		// Add entries that fall within the gap
		// (by a margin of at least edgeMargin to avoid duplicates).
		for i < len(entries) {
			e := entries[i]
			t := e.Time()
			if t.Before(g.Start) {
				break
			}
			if t.Sub(g.Start) >= edgeMargin && g.Finish.Sub(t) >= edgeMargin {
				missing = append(missing, e)
			}
			i++
		}
	}
	return missing
}
