package splitter

import (
	"fmt"
	"sort"
)

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		PrefixSep: "_",
		KeepPrefix: false,
	}
}

// Options controls how Split partitions an env map.
type Options struct {
	// PrefixSep is the separator between prefix and key name (default "_").
	PrefixSep string
	// KeepPrefix retains the full original key name in the output bucket.
	KeepPrefix bool
	// Prefixes is the ordered list of prefixes to split on.
	// Keys that match no prefix land in the "" bucket.
	Prefixes []string
}

// Result holds the split buckets.
type Result struct {
	// Buckets maps prefix (or "") to the env vars that belong to it.
	Buckets map[string]map[string]string
}

// Split partitions env into buckets keyed by prefix.
// Keys not matching any declared prefix are placed in the unmatched bucket ("").
func Split(env map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, fmt.Errorf("splitter: env must not be nil")
	}
	if opts.PrefixSep == "" {
		opts.PrefixSep = "_"
	}

	buckets := make(map[string]map[string]string)
	// Pre-create declared buckets so they always appear in the result.
	for _, p := range opts.Prefixes {
		buckets[p] = make(map[string]string)
	}
	buckets[""] = make(map[string]string)

	keys := sortedKeys(env)
	for _, k := range keys {
		matched := false
		for _, p := range opts.Prefixes {
			token := p + opts.PrefixSep
			if len(k) > len(token) && k[:len(token)] == token {
				outKey := k
				if !opts.KeepPrefix {
					outKey = k[len(token):]
				}
				buckets[p][outKey] = env[k]
				matched = true
				break
			}
		}
		if !matched {
			buckets[""][k] = env[k]
		}
	}

	return &Result{Buckets: buckets}, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
