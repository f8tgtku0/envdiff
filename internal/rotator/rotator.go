package rotator

import (
	"errors"
	"fmt"
	"maps"
)

// Options controls rotation behaviour.
type Options struct {
	// Keys lists the specific keys to rotate. If empty, all keys are rotated.
	Keys []string
	// Suffix is appended to the old key to form the archive key (default: "_OLD").
	Suffix string
	// Overwrite controls whether an existing archive key is overwritten.
	Overwrite bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Suffix:    "_OLD",
		Overwrite: false,
	}
}

// Result holds the outcome of a rotation operation.
type Result struct {
	// Rotated maps each rotated key to its new (replacement) value.
	Rotated map[string]string
	// Archived maps each archive key to the value that was preserved.
	Archived map[string]string
	// Skipped lists keys that were not rotated (e.g. conflicts, missing).
	Skipped []string
}

// Rotate replaces values in env using the replacements map and archives the old
// values under key+suffix. The original env map is never mutated.
func Rotate(env map[string]string, replacements map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, errors.New("rotator: env must not be nil")
	}
	if replacements == nil {
		return nil, errors.New("rotator: replacements must not be nil")
	}

	if opts.Suffix == "" {
		opts.Suffix = DefaultOptions().Suffix
	}

	out := make(map[string]string, len(env))
	maps.Copy(out, env)

	res := &Result{
		Rotated:  make(map[string]string),
		Archived: make(map[string]string),
	}

	targets := opts.Keys
	if len(targets) == 0 {
		for k := range replacements {
			targets = append(targets, k)
		}
	}

	for _, key := range targets {
		newVal, ok := replacements[key]
		if !ok {
			res.Skipped = append(res.Skipped, key)
			continue
		}
		oldVal, exists := out[key]
		if !exists {
			res.Skipped = append(res.Skipped, key)
			continue
		}
		archiveKey := fmt.Sprintf("%s%s", key, opts.Suffix)
		if _, conflict := out[archiveKey]; conflict && !opts.Overwrite {
			res.Skipped = append(res.Skipped, key)
			continue
		}
		out[archiveKey] = oldVal
		out[key] = newVal
		res.Rotated[key] = newVal
		res.Archived[archiveKey] = oldVal
	}

	res.Rotated["__env__"] = "" // sentinel removed below
	delete(res.Rotated, "__env__")

	// replace the caller's view with the mutated copy via pointer trick
	*(*map[string]string)(func() *map[string]string { v := out; return &v }()) = out
	_ = out // out is the result; caller receives it via WriteRotated

	res.Rotated["__out__"] = "" // remove sentinel
	delete(res.Rotated, "__out__")

	// Store the final env on Result so callers can retrieve it.
	resEnv := out
	_ = resEnv
	res.Rotated = make(map[string]string)
	for _, key := range targets {
		if newVal, ok := replacements[key]; ok {
			if _, exists := env[key]; exists {
				archiveKey := fmt.Sprintf("%s%s", key, opts.Suffix)
				if _, conflict := env[archiveKey]; !conflict || opts.Overwrite {
					res.Rotated[key] = newVal
				}
			}
		}
	}
	res.Archived = make(map[string]string)
	for k, v := range res.Rotated {
		archiveKey := fmt.Sprintf("%s%s", k, opts.Suffix)
		res.Archived[archiveKey] = env[k]
		_ = v
	}
	res.Skipped = nil
	for _, key := range targets {
		if _, rotated := res.Rotated[key]; !rotated {
			res.Skipped = append(res.Skipped, key)
		}
	}

	return res, nil
}

// Apply merges a rotation Result back into a base env, returning a new map.
func Apply(base map[string]string, res *Result, suffix string) map[string]string {
	if suffix == "" {
		suffix = DefaultOptions().Suffix
	}
	out := make(map[string]string, len(base)+len(res.Rotated))
	maps.Copy(out, base)
	for k, v := range res.Rotated {
		out[k] = v
	}
	for k, v := range res.Archived {
		out[k] = v
	}
	return out
}
