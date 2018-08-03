package nightscout

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"time"
)

// Implement sort.Interface for reverse chronological order.
// If date fields are the same, compare type fields.

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// After returns true if x comes after y chronologically,
// using the Type field to break ties.
func (x Entry) After(y Entry) bool {
	return x.Date > y.Date || (x.Date == y.Date && x.Type > y.Type)
}

func (e Entries) Less(i, j int) bool {
	return e[i].After(e[j])
}

// Sort sorts the entries into reverse chronological order, in place.
func (e Entries) Sort() {
	sort.Sort(e)
}

// Reverse reverses the entries, in place.
func (e Entries) Reverse() {
	for i, j := 0, len(e)-1; i < len(e)/2; i, j = i+1, j-1 {
		e[i], e[j] = e[j], e[i]
	}
}

// TrimAfter returns the entries that are more recent than the specified time.
// The entries must be in reverse chronological order.
func (e Entries) TrimAfter(cutoff time.Time) Entries {
	d := Date(cutoff)
	n := sort.Search(len(e), func(i int) bool {
		return e[i].Date <= d
	})
	return e[:n]
}

// MergeEntries merges entries that are already in reverse chronological order.
// Duplicates are removed.
func MergeEntries(u, v Entries) Entries {
	m := make(Entries, 0, len(u)+len(v))
	prev := Entry{Date: -1, Type: "invalid"}
	i := 0
	j := 0
	for i < len(u) || j < len(v) {
		if j == len(v) || (i < len(u) && u[i].After(v[j])) {
			m.addEntry(u, &i, &prev)
			continue
		}
		if i == len(u) || (j < len(v) && v[j].After(u[i])) {
			m.addEntry(v, &j, &prev)
			continue
		}
		m.addEntry(u, &i, &prev)
		m.addEntry(v, &j, &prev)
	}
	return m
}

func (e *Entries) addEntry(v Entries, i *int, prev *Entry) {
	cur := v[*i]
	if cur != *prev {
		*e = append(*e, cur)
		*prev = cur
	}
	*i = *i + 1
}

// Write writes entries in JSON format to an io.Writer.
func (e Entries) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(e)
}

// Print writes entries in JSON format on os.Stdout.
func (e Entries) Print() {
	_ = e.Write(os.Stdout)
}

// Save writes entries in JSON format to a file.
func (e Entries) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return e.Write(f)
}

// ReadEntries reads entries in JSON format from a file.
func ReadEntries(file string) (Entries, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	d := json.NewDecoder(f)
	var contents Entries
	err = d.Decode(&contents)
	return contents, err
}
