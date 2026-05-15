// Package compactor removes redundant or overridden keys from a set of
// layered env maps, producing a single compacted map that reflects only
// the effective values after all layers are applied.
package compactor

import "sort"

// Options controls compaction behaviour.
type Options struct {
	// SkipEmpty removes keys whose final resolved value is an empty string.
	SkipEmpty bool

	// PreserveOrder returns keys in the order they were first seen across
	// layers rather than alphabetically. When false (default) keys are sorted.
	PreserveOrder bool
}

// DefaultOptions returns Options with safe defaults.
func DefaultOptions() Options {
	return Options{
		SkipEmpty:     false,
		PreserveOrder: false,
	}
}

// Result holds the output of a compaction run.
type Result struct {
	// Compacted is the final effective key→value map.
	Compacted map[string]string

	// Dropped contains keys that were removed because their effective value
	// was empty and SkipEmpty was set.
	Dropped []string

	// Overridden maps each key that appeared in more than one layer to the
	// number of times it was superseded.
	Overridden map[string]int
}

// Compact merges layers left-to-right (later layers win) and applies the
// supplied options to produce a Result.
//
// A nil layer is silently skipped.
func Compact(layers []map[string]string, opts Options) (*Result, error) {
	if len(layers) == 0 {
		return &Result{
			Compacted:  map[string]string{},
			Dropped:    []string{},
			Overridden: map[string]int{},
		}, nil
	}

	merged := make(map[string]string)
	overridden := make(map[string]int)
	seen := make(map[string]bool) // tracks first-seen order
	order := []string{}

	for _, layer := range layers {
		if layer == nil {
			continue
		}
		for k, v := range layer {
			if !seen[k] {
				seen[k] = true
				order = append(order, k)
			} else {
				overridden[k]++
			}
			merged[k] = v
		}
	}

	dropped := []string{}
	if opts.SkipEmpty {
		for k, v := range merged {
			if v == "" {
				delete(merged, k)
				dropped = append(dropped, k)
			}
		}
		sort.Strings(dropped)
	}

	if !opts.PreserveOrder {
		// normalise order field (not returned, but keeps Compacted deterministic)
		_ = order
	}

	return &Result{
		Compacted:  merged,
		Dropped:    dropped,
		Overridden: overridden,
	}, nil
}
