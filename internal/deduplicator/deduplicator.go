package deduplicator

import "sort"

// Options controls the behaviour of the Deduplicate function.
type Options struct {
	// PreferLast keeps the last-seen value when duplicates are found across
	// ordered sources; when false the first value wins.
	PreferLast bool

	// ReportOnly causes Deduplicate to return the merged map unchanged but
	// still populate the Report with every conflict it found.
	ReportOnly bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		PreferLast: false,
		ReportOnly: false,
	}
}

// Conflict records a key that appeared in more than one source with
// different values.
type Conflict struct {
	Key    string
	Values []string // one entry per source, in order
}

// Report summarises what Deduplicate found.
type Report struct {
	Conflicts []Conflict
}

// HasConflicts returns true when at least one conflicting key was found.
func (r Report) HasConflicts() bool { return len(r.Conflicts) > 0 }

// Deduplicate merges multiple env maps into a single map according to opts.
// Sources are processed in the order supplied; the first (or last, depending
// on opts.PreferLast) non-empty value for a key is kept.
// A Report is always returned describing any keys that had conflicting values.
func Deduplicate(sources []map[string]string, opts Options) (map[string]string, Report) {
	type entry struct {
		value  string
		values []string
	}

	tracked := make(map[string]*entry)
	order := []string{}

	for _, src := range sources {
		for k, v := range src {
			if e, exists := tracked[k]; exists {
				e.values = append(e.values, v)
				if opts.PreferLast {
					e.value = v
				}
			} else {
				tracked[k] = &entry{value: v, values: []string{v}}
				order = append(order, k)
			}
		}
	}

	result := make(map[string]string, len(tracked))
	var conflicts []Conflict

	for _, k := range order {
		e := tracked[k]
		if !opts.ReportOnly {
			result[k] = e.value
		}
		if hasVariation(e.values) {
			conflicts = append(conflicts, Conflict{Key: k, Values: e.values})
		}
	}

	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	if opts.ReportOnly {
		// Return the first-source merged map without modifications.
		for _, k := range order {
			result[k] = tracked[k].values[0]
		}
	}

	return result, Report{Conflicts: conflicts}
}

func hasVariation(values []string) bool {
	for i := 1; i < len(values); i++ {
		if values[i] != values[0] {
			return true
		}
	}
	return false
}
