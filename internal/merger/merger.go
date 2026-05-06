// Package merger provides functionality to merge multiple .env maps
// into a single resolved map, with configurable conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error when a key conflict is detected.
	StrategyError
)

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key    string
	Values []string // values in the order they were encountered
}

// Result holds the merged environment map and any conflicts that were detected.
type Result struct {
	Merged    map[string]string
	Conflicts []Conflict
}

// Merge combines multiple env maps according to the given Strategy.
// Sources are processed in the order they are provided.
func Merge(sources []map[string]string, s Strategy) (*Result, error) {
	merged := make(map[string]string)
	conflictMap := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			existing, exists := merged[k]
			if !exists {
				merged[k] = v
				continue
			}

			if existing == v {
				continue
			}

			// Record conflict values lazily.
			if _, seen := conflictMap[k]; !seen {
				conflictMap[k] = []string{existing}
			}
			conflictMap[k] = append(conflictMap[k], v)

			switch s {
			case StrategyFirst:
				// keep existing — do nothing
			case StrategyLast:
				merged[k] = v
			case StrategyError:
				return nil, fmt.Errorf("merger: conflict on key %q: %q vs %q", k, existing, v)
			}
		}
	}

	var conflicts []Conflict
	for k, vals := range conflictMap {
		conflicts = append(conflicts, Conflict{Key: k, Values: vals})
	}

	return &Result{Merged: merged, Conflicts: conflicts}, nil
}
