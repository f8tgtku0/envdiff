package aliaser

import "fmt"

// Options configures the Alias operation.
type Options struct {
	// Overwrite allows an alias target key to overwrite an existing key.
	Overwrite bool
}

// DefaultOptions returns sensible defaults for aliasing.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
	}
}

// AliasMap maps a new key name to an existing source key name.
// Example: {"DB_HOST": "DATABASE_HOST"} copies the value of DATABASE_HOST into DB_HOST.
type AliasMap map[string]string

// Result holds the outcome of an Alias operation.
type Result struct {
	Applied  []string // new keys that were written
	Skipped  []string // new keys skipped because target already existed
	Missing  []string // source keys that were not found in env
}

// Alias creates new keys in env by copying values from existing keys according
// to the provided AliasMap. The original source keys are preserved.
func Alias(env map[string]string, aliases AliasMap, opts Options) (map[string]string, Result, error) {
	if env == nil {
		return nil, Result{}, fmt.Errorf("aliaser: env must not be nil")
	}
	if len(aliases) == 0 {
		return copyMap(env), Result{}, nil
	}

	out := copyMap(env)
	var res Result

	for newKey, srcKey := range aliases {
		val, ok := env[srcKey]
		if !ok {
			res.Missing = append(res.Missing, srcKey)
			continue
		}
		if _, exists := out[newKey]; exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, newKey)
			continue
		}
		out[newKey] = val
		res.Applied = append(res.Applied, newKey)
	}

	return out, res, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
