// Package filter provides functionality to include or exclude environment
// variable keys based on prefix or pattern matching rules.
package filter

import "strings"

// Options holds the filtering configuration.
type Options struct {
	// IncludePrefixes, if non-empty, keeps only keys that start with one of
	// the given prefixes.
	IncludePrefixes []string

	// ExcludePrefixes removes keys that start with one of the given prefixes.
	// Exclusion is applied after inclusion filtering.
	ExcludePrefixes []string
}

// Apply returns a new map containing only the key/value pairs that pass the
// filter rules defined in opts. The original map is never modified.
func Apply(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))

	for k, v := range env {
		if !included(k, opts.IncludePrefixes) {
			continue
		}
		if excluded(k, opts.ExcludePrefixes) {
			continue
		}
		out[k] = v
	}

	return out
}

// included reports whether key passes the include-prefix filter.
// If no prefixes are specified every key is included.
func included(key string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// excluded reports whether key matches any of the exclude prefixes.
func excluded(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
