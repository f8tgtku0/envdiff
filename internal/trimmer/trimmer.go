// Package trimmer removes unused or redundant keys from an env map
// by comparing it against a reference set of known/expected keys.
package trimmer

import "sort"

// Options controls Trim behaviour.
type Options struct {
	// DryRun reports what would be removed without modifying the map.
	DryRun bool
	// IgnorePrefix skips keys that start with any of these prefixes.
	IgnorePrefix []string
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{}
}

// Result holds the outcome of a Trim operation.
type Result struct {
	// Removed lists keys that were (or would be) removed.
	Removed []string
	// Kept is the resulting env map (nil when DryRun is true).
	Kept map[string]string
}

// Trim removes keys from env that are not present in reference.
// When DryRun is true the original env map is not mutated and Kept is nil.
func Trim(env map[string]string, reference map[string]string, opts Options) Result {
	var removed []string

	for key := range env {
		if ignoredByPrefix(key, opts.IgnorePrefix) {
			continue
		}
		if _, ok := reference[key]; !ok {
			removed = append(removed, key)
		}
	}

	sort.Strings(removed)

	if opts.DryRun {
		return Result{Removed: removed, Kept: nil}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for _, key := range removed {
		delete(out, key)
	}

	return Result{Removed: removed, Kept: out}
}

func ignoredByPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
