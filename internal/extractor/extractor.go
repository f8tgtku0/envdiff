// Package extractor provides utilities for extracting a subset of keys
// from an environment map based on prefix, suffix, or explicit key lists.
package extractor

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how extraction is performed.
type Options struct {
	// Prefixes filters keys that start with any of these strings.
	Prefixes []string
	// Suffixes filters keys that end with any of these strings.
	Suffixes []string
	// Keys is an explicit list of keys to extract.
	Keys []string
	// StripPrefix removes the matched prefix from extracted keys when set.
	StripPrefix bool
}

// DefaultOptions returns an Options with no filters applied.
func DefaultOptions() Options {
	return Options{}
}

// Result holds the extracted key/value pairs and metadata.
type Result struct {
	Env      map[string]string
	Extracted []string // sorted list of extracted key names (post-strip)
}

// Extract returns a new map containing only the keys from env that match
// the criteria defined in opts. If no filter criteria are set, all keys
// are returned.
func Extract(env map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, fmt.Errorf("extractor: env must not be nil")
	}

	out := make(map[string]string)

	for k, v := range env {
		if !matches(k, opts) {
			continue
		}
		outKey := k
		if opts.StripPrefix {
			outKey = stripFirstPrefix(k, opts.Prefixes)
		}
		out[outKey] = v
	}

	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &Result{Env: out, Extracted: keys}, nil
}

func matches(key string, opts Options) bool {
	hasFilter := len(opts.Prefixes) > 0 || len(opts.Suffixes) > 0 || len(opts.Keys) > 0
	if !hasFilter {
		return true
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, s := range opts.Suffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	for _, k := range opts.Keys {
		if key == k {
			return true
		}
	}
	return false
}

func stripFirstPrefix(key string, prefixes []string) string {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return strings.TrimPrefix(key, p)
		}
	}
	return key
}
