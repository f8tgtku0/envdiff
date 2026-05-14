// Package sanitizer provides utilities for cleaning and normalising env var
// values by stripping control characters, trimming whitespace, and optionally
// enforcing a maximum value length.
package sanitizer

import (
	"errors"
	"strings"
	"unicode"
)

// Options controls the behaviour of Sanitize.
type Options struct {
	// TrimSpace removes leading and trailing whitespace from every value.
	TrimSpace bool

	// StripControl removes non-printable / control characters from values.
	StripControl bool

	// MaxValueLen truncates values that exceed this length. Zero means no limit.
	MaxValueLen int
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:    true,
		StripControl: true,
		MaxValueLen:  0,
	}
}

// Result holds the sanitized env map and a record of which keys were mutated.
type Result struct {
	Env     map[string]string
	Changed []string
}

// Sanitize applies the given options to env, returning a new map and a Result
// describing every key whose value was altered. The original map is never
// mutated.
func Sanitize(env map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, errors.New("sanitizer: env must not be nil")
	}

	out := make(map[string]string, len(env))
	var changed []string

	for k, v := range env {
		original := v

		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}

		if opts.StripControl {
			v = strings.Map(func(r rune) rune {
				if unicode.IsControl(r) {
					return -1
				}
				return r
			}, v)
		}

		if opts.MaxValueLen > 0 && len(v) > opts.MaxValueLen {
			v = v[:opts.MaxValueLen]
		}

		out[k] = v
		if v != original {
			changed = append(changed, k)
		}
	}

	return &Result{Env: out, Changed: changed}, nil
}
