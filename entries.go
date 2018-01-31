package nightscout

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"time"
)

// Chronological implements sort.Interface for chronological order.
type Chronological Entries

func (v Chronological) Len() int {
	return len(v)
}

func (v Chronological) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Chronological) Less(i, j int) bool {
	return v[i].Date < v[j].Date
}

// Sort returns a copy of the entries in chronological order.
func (e Entries) Sort() Entries {
	v := make(Entries, len(e))
	copy(v, e)
	sort.Sort(Chronological(v))
	return v
}

// SortReverse returns a copy of the entries in reverse chronological order.
func (e Entries) SortReverse() Entries {
	v := make(Entries, len(e))
	copy(v, e)
	sort.Sort(sort.Reverse(Chronological(v)))
	return v
}

// TrimAfter returns the entries that are more recent than the specified time.
// The entries must be in chronological order.
func (e Entries) TrimAfter(cutoff time.Time) Entries {
	d := Date(cutoff)
	n := sort.Search(len(e), func(i int) bool {
		return e[i].Date > d
	})
	return e[n:]
}

// MergeEntries merges entries that are already in chronological order.
func MergeEntries(u, v Entries) Entries {
	m := make(Entries, 0, len(u)+len(v))
	i := 0
	j := 0
	for {
		if i == len(u) {
			m = append(m, v[j:]...)
			break
		}
		if j == len(v) {
			m = append(m, u[i:]...)
			break
		}
		if u[i].Date < v[j].Date {
			m = append(m, u[i])
			i++
			continue
		}
		if v[j].Date < u[i].Date {
			m = append(m, v[j])
			j++
			continue
		}
		m = append(m, u[i])
		i++
		m = append(m, v[j])
		j++
	}
	return m
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

// Save writes entries in JSON format to a file,
// which is first renamed with a "~" suffix.
func (e Entries) Save(file string) error {
	err := os.Rename(file, file+"~")
	if err != nil {
		return err
	}
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
