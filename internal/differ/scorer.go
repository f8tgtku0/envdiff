package differ

import (
	"sort"
)

// SimilarityOptions controls how environment similarity is scored.
type SimilarityOptions struct {
	// WeightMissing is the penalty applied per missing key (0.0–1.0).
	WeightMissing float64
	// WeightMismatch is the penalty applied per mismatched value (0.0–1.0).
	WeightMismatch float64
}

// DefaultSimilarityOptions returns sensible defaults.
func DefaultSimilarityOptions() SimilarityOptions {
	return SimilarityOptions{
		WeightMissing:  1.0,
		WeightMismatch: 0.5,
	}
}

// SimilarityReport holds the computed similarity score between two environments.
type SimilarityReport struct {
	// Score is a value in [0.0, 1.0] where 1.0 means identical.
	Score float64
	// TotalKeys is the union of all keys across both environments.
	TotalKeys int
	// MissingLeft is the count of keys present in right but absent in left.
	MissingLeft int
	// MissingRight is the count of keys present in left but absent in right.
	MissingRight int
	// Mismatched is the count of keys present in both but with different values.
	Mismatched int
	// SortedKeys is the ordered union of all keys.
	SortedKeys []string
}

// Similarity computes a similarity score between two env maps using the
// provided options. It returns a SimilarityReport with the score and breakdown.
func Similarity(left, right map[string]string, opts SimilarityOptions) SimilarityReport {
	union := make(map[string]struct{})
	for k := range left {
		union[k] = struct{}{}
	}
	for k := range right {
		union[k] = struct{}{}
	}

	keys := make([]string, 0, len(union))
	for k := range union {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var missingLeft, missingRight, mismatched int

	for _, k := range keys {
		lv, inLeft := left[k]
		rv, inRight := right[k]
		switch {
		case inLeft && !inRight:
			missingRight++
		case !inLeft && inRight:
			missingLeft++
		case lv != rv:
			mismatched++
		}
	}

	total := len(keys)
	var score float64
	if total > 0 {
		penalty := opts.WeightMissing*float64(missingLeft+missingRight) +
			opts.WeightMismatch*float64(mismatched)
		effective := float64(total) + opts.WeightMissing*float64(missingLeft+missingRight)
		score = 1.0 - penalty/effective
		if score < 0 {
			score = 0
		}
	} else {
		score = 1.0
	}

	return SimilarityReport{
		Score:        score,
		TotalKeys:    total,
		MissingLeft:  missingLeft,
		MissingRight: missingRight,
		Mismatched:   mismatched,
		SortedKeys:   keys,
	}
}
