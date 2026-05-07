package renamer

import "fmt"

// RenameMap is a mapping from old key names to new key names.
type RenameMap = map[string]string

// Result holds the outcome of a rename operation.
type Result struct {
	// Renamed contains keys that were successfully renamed: old -> new.
	Renamed map[string]string
	// Skipped contains old key names that were not found in the source env.
	Skipped []string
	// Conflicts contains new key names that already existed in the source env.
	Conflicts []string
}

// Options controls the behaviour of Rename.
type Options struct {
	// OverwriteConflicts allows renaming even when the target key already exists.
	OverwriteConflicts bool
}

// DefaultOptions returns sensible defaults for Rename.
func DefaultOptions() Options {
	return Options{OverwriteConflicts: false}
}

// Rename applies a RenameMap to env, returning a new map with keys renamed
// according to the mapping. The original env is never mutated.
func Rename(env map[string]string, renames RenameMap, opts Options) (map[string]string, Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	result := Result{
		Renamed: make(map[string]string),
	}

	for oldKey, newKey := range renames {
		val, exists := out[oldKey]
		if !exists {
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}

		if _, conflict := out[newKey]; conflict && oldKey != newKey {
			if !opts.OverwriteConflicts {
				result.Conflicts = append(result.Conflicts, newKey)
				return nil, result, fmt.Errorf("rename conflict: target key %q already exists", newKey)
			}
			result.Conflicts = append(result.Conflicts, newKey)
		}

		delete(out, oldKey)
		out[newKey] = val
		result.Renamed[oldKey] = newKey
	}

	return out, result, nil
}
