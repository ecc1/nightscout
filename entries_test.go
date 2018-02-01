package nightscout

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

var (
	D = Entries{
		{Date: 0},
		{Date: 1},
		{Date: 2},
		{Date: 3},
		{Date: 4},
		{Date: 5},
		{Date: 6},
		{Date: 7},
		{Date: 8},
		{Date: 9}}

	dups  = append(D, D...).Sort()
	evens = Entries{D[0], D[2], D[4], D[6], D[8]}
	odds  = Entries{D[1], D[3], D[5], D[7], D[9]}
)

type P Entries

func (v P) String() string {
	if len(v) == 0 {
		return "[]"
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[%d", v[0].Date))
	for i := 1; i < len(v); i++ {
		buf.WriteString(fmt.Sprintf(" %d", v[i].Date))
	}
	buf.WriteString("]")
	return buf.String()
}

func TestSortEntries(t *testing.T) {
	cases := []struct {
		unsorted Entries
		sorted   Entries
	}{
		{D, D},
		{D[:0], D[:0]},
		{D[0:1], D[0:1]},
		{Entries{D[3], D[2], D[1]}, D[1:4]},
	}
	for _, c := range cases {
		v := c.unsorted.Sort()
		if !reflect.DeepEqual(v, c.sorted) {
			t.Errorf("Sort(%v) == %v, want %v", P(c.unsorted), P(v), P(c.sorted))
		}
	}
}

func TestSortReverseEntries(t *testing.T) {
	cases := []struct {
		unsorted Entries
		sorted   Entries
	}{
		{D[:0], D[:0]},
		{D[0:1], D[0:1]},
		{D[1:4], Entries{D[3], D[2], D[1]}},
	}
	for _, c := range cases {
		v := c.unsorted.SortReverse()
		if !reflect.DeepEqual(v, c.sorted) {
			t.Errorf("SortReverse(%v) == %v, want %v", P(c.unsorted), P(v), P(c.sorted))
		}
	}
}

func TestMergeEntries(t *testing.T) {
	cases := []struct {
		a, b   Entries
		merged Entries
	}{
		{D[:0], D[:0], D[:0]},
		{D, nil, D},
		{nil, D, D},
		{D, D, D},
		{dups, nil, D},
		{dups, D, D},
		{nil, dups, D},
		{D, dups, D},
		{D[:1], D[1:], D},
		{D[1:], D[:1], D},
		{D[:5], D[5:], D},
		{D[5:], D[:5], D},
		{D[:9], D[9:], D},
		{D[9:], D[:9], D},
		{evens, odds, D},
		{odds, evens, D},
	}
	for _, c := range cases {
		// Make a copy since sorting is done in place.
		v := MergeEntries(c.a, c.b)
		if !reflect.DeepEqual(v, c.merged) {
			t.Errorf("MergeEntries(%v, %v) == %v, want %v", P(c.a), P(c.b), P(v), P(c.merged))
		}
	}
}
