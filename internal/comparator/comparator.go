package comparator

import (
	"fmt"
	"sort"
)

// Result holds the outcome of a multi-env key comparison.
type Result struct {
	// Key is the environment variable name.
	Key string
	// Values maps each environment label to its value.
	Values map[string]string
	// Consistent is true when all envs that define the key share the same value.
	Consistent bool
	// Missing is the list of env labels that do not define the key at all.
	Missing []string
}

// Options controls comparator behaviour.
type Options struct {
	// OnlyInconsistent skips keys that are fully consistent across all envs.
	OnlyInconsistent bool
	// RequireAll reports a key as inconsistent when any env is missing it.
	RequireAll bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		OnlyInconsistent: false,
		RequireAll:       true,
	}
}

// Compare performs a value-level comparison of a shared key across multiple
// named environments. envs maps an environment label (e.g. "staging") to its
// parsed key→value map.
func Compare(envs map[string]map[string]string, opts Options) ([]Result, error) {
	if len(envs) == 0 {
		return nil, fmt.Errorf("comparator: at least one environment is required")
	}

	// Collect the union of all keys.
	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	labels := sortedLabels(envs)

	var results []Result
	for key := range keySet {
		res := buildResult(key, labels, envs, opts)
		if opts.OnlyInconsistent && res.Consistent {
			continue
		}
		results = append(results, res)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results, nil
}

func buildResult(key string, labels []string, envs map[string]map[string]string, opts Options) Result {
	values := make(map[string]string, len(labels))
	var missing []string

	for _, label := range labels {
		v, ok := envs[label][key]
		if !ok {
			missing = append(missing, label)
			continue
		}
		values[label] = v
	}

	consistent := isConsistent(values, missing, opts)
	return Result{
		Key:        key,
		Values:     values,
		Consistent: consistent,
		Missing:    missing,
	}
}

func isConsistent(values map[string]string, missing []string, opts Options) bool {
	if opts.RequireAll && len(missing) > 0 {
		return false
	}
	var ref string
	first := true
	for _, v := range values {
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return false
		}
	}
	return true
}

func sortedLabels(envs map[string]map[string]string) []string {
	labels := make([]string, 0, len(envs))
	for l := range envs {
		labels = append(labels, l)
	}
	sort.Strings(labels)
	return labels
}
