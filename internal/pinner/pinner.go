package pinner

import (
	"errors"
	"fmt"
	"sort"
)

// Options controls the behaviour of Pin.
type Options struct {
	// Overwrite allows an existing pinned value to be replaced.
	Overwrite bool
	// StrictMissing returns an error if a key to pin does not exist in env.
	StrictMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite:     false,
		StrictMissing: false,
	}
}

// PinEntry records the outcome for a single key.
type PinEntry struct {
	Key      string
	OldValue string
	NewValue string
	Pinned   bool
	Skipped  bool
	Reason   string
}

// Result holds all pin outcomes.
type Result struct {
	Entries []PinEntry
}

// Pinned returns entries that were successfully pinned.
func (r *Result) Pinned() []PinEntry {
	var out []PinEntry
	for _, e := range r.Entries {
		if e.Pinned {
			out = append(out, e)
		}
	}
	return out
}

// Skipped returns entries that were not changed.
func (r *Result) Skipped() []PinEntry {
	var out []PinEntry
	for _, e := range r.Entries {
		if e.Skipped {
			out = append(out, e)
		}
	}
	return out
}

// Pin locks specific keys in env to the values provided in pins.
// The env map is mutated in place and a Result is returned.
func Pin(env map[string]string, pins map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, errors.New("pinner: env must not be nil")
	}
	if pins == nil {
		return &Result{}, nil
	}

	keys := sortedKeys(pins)
	result := &Result{}

	for _, k := range keys {
		newVal := pins[k]
		oldVal, exists := env[k]

		if !exists && opts.StrictMissing {
			return nil, fmt.Errorf("pinner: key %q not found in env", k)
		}

		if !exists {
			env[k] = newVal
			result.Entries = append(result.Entries, PinEntry{
				Key:      k,
				OldValue: "",
				NewValue: newVal,
				Pinned:   true,
				Reason:   "added",
			})
			continue
		}

		if oldVal == newVal {
			result.Entries = append(result.Entries, PinEntry{
				Key:      k,
				OldValue: oldVal,
				NewValue: newVal,
				Skipped:  true,
				Reason:   "already matches",
			})
			continue
		}

		if !opts.Overwrite {
			result.Entries = append(result.Entries, PinEntry{
				Key:      k,
				OldValue: oldVal,
				NewValue: newVal,
				Skipped:  true,
				Reason:   "overwrite disabled",
			})
			continue
		}

		env[k] = newVal
		result.Entries = append(result.Entries, PinEntry{
			Key:      k,
			OldValue: oldVal,
			NewValue: newVal,
			Pinned:   true,
			Reason:   "overwritten",
		})
	}

	return result, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
