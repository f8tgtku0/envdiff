// Package resolver provides utilities for resolving a canonical environment
// map from multiple .env files, applying precedence rules so that later
// files override earlier ones unless a key is already locked.
package resolver

import "fmt"

// Options controls how resolution behaves.
type Options struct {
	// Overwrite allows later files to overwrite keys from earlier files.
	// When false the first value seen for a key wins (left-priority).
	Overwrite bool

	// Strict causes Resolve to return an error if the same key appears in
	// more than one source with a different value.
	Strict bool
}

// DefaultOptions returns an Options value with sensible defaults:
// later files override earlier ones and strict mode is off.
func DefaultOptions() Options {
	return Options{Overwrite: true, Strict: false}
}

// Source pairs a label (e.g. a filename) with its parsed key/value map.
type Source struct {
	Label string
	Env   map[string]string
}

// Resolve merges multiple Sources into a single map according to opts.
// It also returns a Provenance map that records which label each key came from.
func Resolve(sources []Source, opts Options) (env map[string]string, provenance map[string]string, err error) {
	env = make(map[string]string)
	provenance = make(map[string]string)

	for _, src := range sources {
		for k, v := range src.Env {
			existing, seen := env[k]
			if !seen {
				env[k] = v
				provenance[k] = src.Label
				continue
			}

			if opts.Strict && existing != v {
				return nil, nil, fmt.Errorf(
					"resolver: key %q has conflicting values in %q and %q",
					k, provenance[k], src.Label,
				)
			}

			if opts.Overwrite {
				env[k] = v
				provenance[k] = src.Label
			}
		}
	}

	return env, provenance, nil
}
