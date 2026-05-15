// Package cloner provides utilities for cloning and transforming env maps
// with optional key prefix injection or stripping.
package cloner

import "fmt"

// Options controls the behaviour of Clone.
type Options struct {
	// AddPrefix is prepended to every key in the output map.
	AddPrefix string
	// StripPrefix is removed from the start of every key before cloning.
	// Keys that do not carry the prefix are included unchanged unless
	// StrictPrefix is true.
	StripPrefix string
	// StrictPrefix causes Clone to return an error if StripPrefix is set
	// and a key is encountered that does not start with it.
	StrictPrefix bool
	// Overwrite controls whether existing keys in dst are replaced when
	// merging into a non-nil destination map.
	Overwrite bool
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{Overwrite: true}
}

// Clone copies src into a new map, optionally stripping and/or adding a key
// prefix. If dst is non-nil the cloned pairs are merged into it according to
// the Overwrite flag and the merged map is returned.
func Clone(src map[string]string, dst map[string]string, opts Options) (map[string]string, error) {
	if src == nil {
		return nil, fmt.Errorf("cloner: src must not be nil")
	}

	if dst == nil {
		dst = make(map[string]string, len(src))
	}

	for k, v := range src {
		outKey := k

		if opts.StripPrefix != "" {
			if len(k) >= len(opts.StripPrefix) && k[:len(opts.StripPrefix)] == opts.StripPrefix {
				outKey = k[len(opts.StripPrefix):]
			} else if opts.StrictPrefix {
				return nil, fmt.Errorf("cloner: key %q does not start with prefix %q", k, opts.StripPrefix)
			}
		}

		if opts.AddPrefix != "" {
			outKey = opts.AddPrefix + outKey
		}

		if _, exists := dst[outKey]; exists && !opts.Overwrite {
			continue
		}
		dst[outKey] = v
	}

	return dst, nil
}
