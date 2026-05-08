// Package sorter provides utilities for sorting and reordering
// key-value pairs from .env files alphabetically or by custom order.
package sorter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls how sorting is performed.
type Options struct {
	// Reverse sorts keys in descending order when true.
	Reverse bool
	// GroupPrefixes groups keys sharing the same prefix together before sorting.
	GroupPrefixes bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		Reverse:       false,
		GroupPrefixes: false,
	}
}

// Sort returns a new map with the same contents as env (unchanged) and
// writes the sorted key=value lines to w. The original map is not mutated.
func Sort(env map[string]string, opts Options, w io.Writer) error {
	if env == nil {
		return fmt.Errorf("sorter: env map must not be nil")
	}
	if w == nil {
		return fmt.Errorf("sorter: writer must not be nil")
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	if opts.GroupPrefixes {
		sort.Slice(keys, func(i, j int) bool {
			pi := prefix(keys[i])
			pj := prefix(keys[j])
			if pi != pj {
				if opts.Reverse {
					return pi > pj
				}
				return pi < pj
			}
			if opts.Reverse {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	} else {
		sort.Slice(keys, func(i, j int) bool {
			if opts.Reverse {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	}

	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, env[k]); err != nil {
			return fmt.Errorf("sorter: write error: %w", err)
		}
	}
	return nil
}

// prefix returns the portion of a key before the first underscore,
// or the whole key if no underscore is present.
func prefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
