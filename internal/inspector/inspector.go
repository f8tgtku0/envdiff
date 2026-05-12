package inspector

import "sort"

// KeyInfo holds metadata about a single key across all environments.
type KeyInfo struct {
	Key        string
	Envs       []string
	Values     map[string]string // env name -> value
	Consistent bool             // true when all envs share the same value
}

// Result is the output of an Inspect call.
type Result struct {
	Keys []KeyInfo
}

// Options configures the inspector.
type Options struct {
	// OnlyInconsistent filters out keys that are identical across all envs.
	OnlyInconsistent bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{OnlyInconsistent: false}
}

// Inspect analyses a set of named environments and returns per-key metadata.
// The envs map is keyed by environment name (e.g. "staging", "prod").
func Inspect(envs map[string]map[string]string, opts Options) Result {
	keySet := make(map[string]struct{})
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var infos []KeyInfo
	for _, k := range keys {
		info := KeyInfo{
			Key:    k,
			Values: make(map[string]string),
		}
		for _, name := range envNames {
			if v, ok := envs[name][k]; ok {
				info.Envs = append(info.Envs, name)
				info.Values[name] = v
			}
		}
		info.Consistent = isConsistent(info.Values)
		if opts.OnlyInconsistent && info.Consistent {
			continue
		}
		infos = append(infos, info)
	}

	return Result{Keys: infos}
}

func isConsistent(values map[string]string) bool {
	var first string
	set := false
	for _, v := range values {
		if !set {
			first = v
			set = true
			continue
		}
		if v != first {
			return false
		}
	}
	return true
}
