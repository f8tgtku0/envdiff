package patcher

import (
	"fmt"
	"os"
	"strings"
)

// Options controls patching behaviour.
type Options struct {
	// Overwrite existing keys when applying a patch.
	Overwrite bool
	// DryRun reports what would change without writing anything.
	DryRun bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		DryRun:    false,
	}
}

// Change describes a single key-level modification.
type Change struct {
	Key      string
	OldValue string // empty when the key is new
	NewValue string
	Action   string // "add", "update", "skip"
}

// Apply merges patch into base according to opts.
// It returns the resulting env map and a list of changes made (or that would
// be made in dry-run mode).
func Apply(base, patch map[string]string, opts Options) (map[string]string, []Change, error) {
	if base == nil {
		return nil, nil, fmt.Errorf("patcher: base map must not be nil")
	}
	if patch == nil {
		return nil, nil, fmt.Errorf("patcher: patch map must not be nil")
	}

	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	var changes []Change

	for k, newVal := range patch {
		oldVal, exists := out[k]
		switch {
		case !exists:
			changes = append(changes, Change{Key: k, OldValue: "", NewValue: newVal, Action: "add"})
			if !opts.DryRun {
				out[k] = newVal
			}
		case opts.Overwrite:
			changes = append(changes, Change{Key: k, OldValue: oldVal, NewValue: newVal, Action: "update"})
			if !opts.DryRun {
				out[k] = newVal
			}
		default:
			changes = append(changes, Change{Key: k, OldValue: oldVal, NewValue: newVal, Action: "skip"})
		}
	}

	return out, changes, nil
}

// WritePatched serialises result to path in KEY=VALUE format.
func WritePatched(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		if strings.ContainsAny(v, " \t") {
			v = `"` + v + `"`
		}
		sb.WriteString(k + "=" + v + "\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
