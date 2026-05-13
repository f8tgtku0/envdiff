package tagger

import (
	"fmt"
	"regexp"
	"sort"
)

// Tag represents a label applied to a key.
type Tag struct {
	Key   string
	Label string
}

// Options controls tagging behaviour.
type Options struct {
	// Rules maps a label to a regexp pattern matched against key names.
	Rules map[string]string
	// AllowMultiple permits a key to receive more than one tag.
	AllowMultiple bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Rules: map[string]string{
			"secret":   `(?i)(secret|token|password|key|pwd)`,
			"database": `(?i)(db_|database|dsn|postgres|mysql|mongo)`,
			"feature":  `(?i)(feature_|flag_|ff_)`,
		},
		AllowMultiple: false,
	}
}

// Result holds the tagged output for a single env map.
type Result struct {
	// Tags maps each key to its assigned label(s).
	Tags map[string][]string
	// Untagged contains keys that matched no rule.
	Untagged []string
}

// Tag applies label rules to env and returns a Result.
func Tag(env map[string]string, opts Options) (*Result, error) {
	compiled := make(map[string]*regexp.Regexp, len(opts.Rules))
	for label, pattern := range opts.Rules {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("tagger: invalid pattern for label %q: %w", label, err)
		}
		compiled[label] = re
	}

	res := &Result{
		Tags: make(map[string][]string),
	}

	keys := sortedKeys(env)
	for _, k := range keys {
		var matched []string
		for label, re := range compiled {
			if re.MatchString(k) {
				matched = append(matched, label)
				if !opts.AllowMultiple {
					break
				}
			}
		}
		if len(matched) == 0 {
			res.Untagged = append(res.Untagged, k)
		} else {
			sort.Strings(matched)
			res.Tags[k] = matched
		}
	}
	return res, nil
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
