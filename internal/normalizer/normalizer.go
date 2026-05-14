// Package normalizer provides utilities for normalising environment
// variable maps before comparison or export — trimming whitespace,
// canonicalising key casing, and removing blank values.
package normalizer

import (
	"strings"
)

// Options controls which normalisation steps are applied.
type Options struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool

	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool

	// RemoveEmpty drops entries whose value is empty after other transforms.
	RemoveEmpty bool
}

// DefaultOptions returns Options with all safe normalisation steps enabled.
func DefaultOptions() Options {
	return Options{
		TrimSpace:     true,
		UppercaseKeys: false,
		RemoveEmpty:   false,
	}
}

// Normalize applies the given Options to a copy of env and returns the
// normalised map. The original map is never mutated.
func Normalize(env map[string]string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, errNilEnv
	}

	out := make(map[string]string, len(env))

	for k, v := range env {
		key := k
		val := v

		if opts.TrimSpace {
			key = strings.TrimSpace(key)
			val = strings.TrimSpace(val)
		}

		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}

		if opts.RemoveEmpty && val == "" {
			continue
		}

		out[key] = val
	}

	return out, nil
}

// errNilEnv is returned when a nil map is passed to Normalize.
var errNilEnv = normalizerError("normalizer: env map must not be nil")

type normalizerError string

func (e normalizerError) Error() string { return string(e) }
