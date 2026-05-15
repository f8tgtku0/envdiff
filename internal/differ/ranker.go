package differ

import (
	"sort"
)

// RankEntry represents a key ranked by its "difference score" across
// multiple comparison results. A higher score means the key is more
// frequently mismatched or missing.
type RankEntry struct {
	Key   string
	Score int
}

// RankOptions controls how Rank weights different diff types.
type RankOptions struct {
	// MissingWeight is added to the score each time a key is absent
	// from one side of a comparison. Defaults to 1.
	MissingWeight int

	// MismatchWeight is added to the score each time a key is present
	// on both sides but has differing values. Defaults to 2.
	MismatchWeight int
}

// DefaultRankOptions returns sensible defaults for RankOptions.
func DefaultRankOptions() RankOptions {
	return RankOptions{
		MissingWeight:  1,
		MismatchWeight: 2,
	}
}

// Rank aggregates one or more Result values and returns a slice of
// RankEntry sorted descending by score (most problematic keys first).
// Keys that appear only in clean results receive a score of zero and
// are omitted from the output unless every result is clean.
func Rank(opts RankOptions, results ...*Result) []RankEntry {
	if opts.MissingWeight == 0 {
		opts.MissingWeight = 1
	}
	if opts.MismatchWeight == 0 {
		opts.MismatchWeight = 2
	}

	scores := make(map[string]int)

	for _, r := range results {
		if r == nil {
			continue
		}
		for _, k := range r.MissingInRight {
			scores[k] += opts.MissingWeight
		}
		for _, k := range r.MissingInLeft {
			scores[k] += opts.MissingWeight
		}
		for k := range r.Mismatched {
			scores[k] += opts.MismatchWeight
		}
	}

	entries := make([]RankEntry, 0, len(scores))
	for k, s := range scores {
		if s > 0 {
			entries = append(entries, RankEntry{Key: k, Score: s})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Key < entries[j].Key
	})

	return entries
}
