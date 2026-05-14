package scoper

import "sort"

// Scope represents a named collection of env keys grouped under a logical scope label.
type Scope struct {
	Name string
	Keys []string
}

// Result holds all scopes extracted from an environment map.
type Result struct {
	Scopes []Scope
	Unscoped []string
}

// Options controls how scoping is performed.
type Options struct {
	// Rules maps a scope name to a list of key prefixes that belong to it.
	Rules map[string][]string
	// IncludeUnscoped controls whether keys that match no rule are collected.
	IncludeUnscoped bool
}

// DefaultOptions returns sensible defaults with no rules.
func DefaultOptions() Options {
	return Options{
		Rules:           map[string][]string{},
		IncludeUnscoped: true,
	}
}

// Scope partitions the env map into named scopes based on prefix rules.
func Scope(env map[string]string, opts Options) (Result, error) {
	if env == nil {
		return Result{}, nil
	}

	assigned := make(map[string]bool)
	scopeMap := make(map[string][]string)

	for scopeName, prefixes := range opts.Rules {
		for key := range env {
			for _, prefix := range prefixes {
				if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
					scopeMap[scopeName] = append(scopeMap[scopeName], key)
					assigned[key] = true
					break
				}
			}
		}
	}

	var scopes []Scope
	scopeNames := sortedStringKeys(scopeMap)
	for _, name := range scopeNames {
		keys := scopeMap[name]
		sort.Strings(keys)
		scopes = append(scopes, Scope{Name: name, Keys: keys})
	}

	var unscoped []string
	if opts.IncludeUnscoped {
		for key := range env {
			if !assigned[key] {
				unscoped = append(unscoped, key)
			}
		}
		sort.Strings(unscoped)
	}

	return Result{Scopes: scopes, Unscoped: unscoped}, nil
}

func sortedStringKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
