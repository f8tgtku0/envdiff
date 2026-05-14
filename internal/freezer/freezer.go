package freezer

import (
	"errors"
	"fmt"
	"sort"
)

// Options configures the Freeze operation.
type Options struct {
	// AllowExpand permits keys to be added to a frozen env without error.
	AllowExpand bool
	// IgnoreKeys is a set of keys to skip during enforcement.
	IgnoreKeys []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		AllowExpand: false,
	}
}

// Violation describes a single frozen-env breach.
type Violation struct {
	Key    string
	Reason string
}

// Result holds the outcome of a Freeze enforcement check.
type Result struct {
	Violations []Violation
}

// Clean returns true when no violations were found.
func (r Result) Clean() bool { return len(r.Violations) == 0 }

// Freeze compares a live env map against a previously frozen snapshot and
// returns any violations (removed keys, changed values, or unexpected additions
// when AllowExpand is false).
func Freeze(frozen, live map[string]string, opts Options) (Result, error) {
	if frozen == nil {
		return Result{}, errors.New("freezer: frozen snapshot must not be nil")
	}
	if live == nil {
		return Result{}, errors.New("freezer: live env must not be nil")
	}

	ignored := make(map[string]bool, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignored[k] = true
	}

	var violations []Violation

	for _, k := range sortedKeys(frozen) {
		if ignored[k] {
			continue
		}
		lv, ok := live[k]
		if !ok {
			violations = append(violations, Violation{Key: k, Reason: "key removed from live env"})
			continue
		}
		if lv != frozen[k] {
			violations = append(violations, Violation{
				Key:    k,
				Reason: fmt.Sprintf("value changed (frozen=%q, live=%q)", frozen[k], lv),
			})
		}
	}

	if !opts.AllowExpand {
		for _, k := range sortedKeys(live) {
			if ignored[k] {
				continue
			}
			if _, ok := frozen[k]; !ok {
				violations = append(violations, Violation{Key: k, Reason: "unexpected key added to live env"})
			}
		}
	}

	return Result{Violations: violations}, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
