// Package stripper removes keys from an env map based on prefix, suffix,
// pattern, or an explicit key list, returning a cleaned copy.
package stripper

import (
	"fmt"
	"regexp"
	"strings"
)

// Options controls which keys are stripped from the env map.
type Options struct {
	// Prefixes removes any key that starts with one of these strings.
	Prefixes []string
	// Suffixes removes any key that ends with one of these strings.
	Suffixes []string
	// Patterns removes any key whose name matches one of these regular expressions.
	Patterns []string
	// Keys removes exactly these keys (case-sensitive).
	Keys []string
	// DryRun returns a report without mutating the returned map.
	DryRun bool
}

// DefaultOptions returns an Options with no rules configured.
func DefaultOptions() Options { return Options{} }

// Result holds the output of a Strip operation.
type Result struct {
	Env     map[string]string
	Stripped []string
}

// Strip applies the given options to env and returns a Result.
// The original map is never mutated.
func Strip(env map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, fmt.Errorf("stripper: env must not be nil")
	}

	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("stripper: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	explicit := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = true
	}

	out := make(map[string]string, len(env))
	var stripped []string

	for k, v := range env {
		if shouldStrip(k, opts.Prefixes, opts.Suffixes, compiled, explicit) {
			stripped = append(stripped, k)
			continue
		}
		out[k] = v
	}

	if opts.DryRun {
		out = copyMap(env)
	}

	return &Result{Env: out, Stripped: stripped}, nil
}

func shouldStrip(key string, prefixes, suffixes []string, patterns []*regexp.Regexp, explicit map[string]bool) bool {
	if explicit[key] {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, s := range suffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
