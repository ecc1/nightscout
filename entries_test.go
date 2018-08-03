package nightscout

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var (
	dups  = append(E, E...).sorted()
	evens = E.skipping(1)
	odds  = E[1:].skipping(1)
)

func (e Entries) sorted() Entries {
	v := make(Entries, len(e))
	copy(v, e)
	v.Sort()
	return v
}

func (e Entries) reversed() Entries {
	v := make(Entries, len(e))
	for i, x := range e {
		v[len(e)-i-1] = x
	}
	return v
}

func (e Entries) skipping(n int) Entries {
	v := make(Entries, 0, len(e)/(n+1))
	for i := 0; i < len(e); i += n {
		v = append(v, e[i])
	}
	return v
}

// Custom String method for printing entries.
func (e Entries) String() string {
	if len(e) == 0 {
		return "[]"
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[%s", entryName(e[0])))
	for _, x := range e[1:] {
		buf.WriteString(fmt.Sprintf(" %s", entryName(x)))
	}
	buf.WriteString("]")
	return buf.String()
}

// Use a symbolic name if e is an element of E, otherwise use its date.
func entryName(e Entry) string {
	for i, x := range E {
		if x == e {
			return fmt.Sprintf("E%d", i)
		}
	}
	return e.String()
}

func TestSortEntries(t *testing.T) {
	cases := []struct {
		unsorted Entries
		sorted   Entries
	}{
		{E, E},
		{E[:0], E[:0]},
		{E[0:1], E[0:1]},
		{Entries{E[3], E[2], E[1]}, E[1:4]},
		{E.reversed(), E},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			// Make a copy since sorting is done in place.
			v := make(Entries, len(c.unsorted))
			copy(v, c.unsorted)
			v.Sort()
			if !reflect.DeepEqual(v, c.sorted) {
				t.Errorf("Sort(%v) == %v, want %v", c.unsorted, v, c.sorted)
			}
		})
	}
}

func TestTrimEntries(t *testing.T) {
	cases := []struct {
		untrimmed Entries
		cutoff    time.Time
		trimmed   Entries
	}{
		{E, parseTime("2017-09-30 23:59:59"), E},
		{E, parseTime("2017-10-01 00:00:00"), E[:19]},
		{E, parseTime("2017-10-01 01:00:00"), E[:7]},
		{E, parseTime("2017-10-01 01:20:00"), E[:3]},
		{E, parseTime("2017-10-01 01:35:00"), E[:0]},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			v := c.untrimmed.TrimAfter(c.cutoff)
			if !reflect.DeepEqual(v, c.trimmed) {
				t.Errorf("TrimAfter(%v, %v) == %v, want %v", c.untrimmed, c.cutoff, v, c.trimmed)
			}
		})
	}
}

func TestMergeEntries(t *testing.T) {
	cases := []struct {
		a, b   Entries
		merged Entries
	}{
		{E[:0], E[:0], E[:0]},
		{E, nil, E},
		{nil, E, E},
		{E, E, E},
		{dups, nil, E},
		{dups, E, E},
		{nil, dups, E},
		{E, dups, E},
		{E[:1], E[1:], E},
		{E[1:], E[:1], E},
		{E[:5], E[5:], E},
		{E[5:], E[:5], E},
		{E[:9], E[9:], E},
		{E[9:], E[:9], E},
		{evens, odds, E},
		{odds, evens, E},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			v := MergeEntries(c.a, c.b)
			if !reflect.DeepEqual(v, c.merged) {
				t.Errorf("MergeEntries(%v, %v) == %v, want %v", c.a, c.b, v, c.merged)
			}
		})
	}
}
