// Package promoter copies selected keys from one environment map into another,
// optionally overwriting existing values and filtering by prefix.
package promoter

import (
	"errors"
	"fmt"
	"strings"
)

// Options controls how promotion behaves.
type Options struct {
	// Overwrite allows keys already present in the destination to be replaced.
	Overwrite bool

	// Prefix restricts promotion to keys that start with the given string.
	// An empty string means all keys are eligible.
	Prefix string

	// Keys is an explicit allow-list of key names to promote.
	// When non-empty, only these keys are considered (Prefix is still applied).
	Keys []string
}

// DefaultOptions returns a safe default configuration.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
	}
}

// Result holds the outcome of a promotion operation.
type Result struct {
	// Promoted lists keys that were copied into the destination.
	Promoted []string
	// Skipped lists keys that were skipped because they already existed.
	Skipped []string
	// Ignored lists keys excluded by the prefix or key filter.
	Ignored []string
}

// Promote copies eligible keys from src into dst according to opts.
// dst must not be nil; src may be empty but not nil.
func Promote(src, dst map[string]string, opts Options) (*Result, error) {
	if src == nil {
		return nil, errors.New("promoter: src must not be nil")
	}
	if dst == nil {
		return nil, errors.New("promoter: dst must not be nil")
	}

	allowSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowSet[k] = true
	}

	res := &Result{}

	for k, v := range src {
		// Filter by prefix.
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			res.Ignored = append(res.Ignored, k)
			continue
		}
		// Filter by explicit key list.
		if len(allowSet) > 0 && !allowSet[k] {
			res.Ignored = append(res.Ignored, k)
			continue
		}

		if _, exists := dst[k]; exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, k)
			continue
		}

		dst[k] = v
		res.Promoted = append(res.Promoted, fmt.Sprintf("%s", k))
	}

	sortStrings(res.Promoted)
	sortStrings(res.Skipped)
	sortStrings(res.Ignored)

	return res, nil
}
