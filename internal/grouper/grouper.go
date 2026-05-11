package grouper

import (
	"sort"
	"strings"
)

// Options controls how keys are grouped.
type Options struct {
	// Delimiter separates the prefix from the rest of the key.
	// Defaults to "_".
	Delimiter string

	// MaxDepth is the number of delimiter-separated segments used to form
	// the group name. 0 or 1 both mean "first segment only".
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Delimiter: "_",
		MaxDepth:  1,
	}
}

// Group holds all keys that share a common prefix.
type Group struct {
	Prefix string
	Keys   []string
	Env    map[string]string
}

// Result is the output of a Group operation.
type Result struct {
	Groups   []*Group
	// Ungrouped contains keys that had no delimiter and therefore no prefix.
	Ungrouped []string
}

// GroupEnv partitions the env map into groups based on key prefixes.
func GroupEnv(env map[string]string, opts Options) Result {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}
	if opts.MaxDepth < 1 {
		opts.MaxDepth = 1
	}

	prefixMap := make(map[string]*Group)
	var ungrouped []string

	for k, v := range env {
		prefix := extractPrefix(k, opts.Delimiter, opts.MaxDepth)
		if prefix == "" {
			ungrouped = append(ungrouped, k)
			continue
		}
		g, ok := prefixMap[prefix]
		if !ok {
			g = &Group{
				Prefix: prefix,
				Env:    make(map[string]string),
			}
			prefixMap[prefix] = g
		}
		g.Keys = append(g.Keys, k)
		g.Env[k] = v
	}

	groups := make([]*Group, 0, len(prefixMap))
	for _, g := range prefixMap {
		sort.Strings(g.Keys)
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})
	sort.Strings(ungrouped)

	return Result{Groups: groups, Ungrouped: ungrouped}
}

func extractPrefix(key, delimiter string, depth int) string {
	parts := strings.SplitN(key, delimiter, depth+1)
	if len(parts) <= 1 {
		return ""
	}
	return strings.Join(parts[:depth], delimiter)
}
