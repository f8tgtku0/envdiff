// Package linter provides style and convention checks for .env file entries.
// It inspects parsed key-value maps and reports issues such as lowercase keys,
// keys with whitespace, empty values, and duplicate keys across files.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single linting problem found in an env map.
type Issue struct {
	Key     string
	Message string
}

// String returns a human-readable representation of the issue.
func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s", i.Key, i.Message)
}

// Options controls which lint rules are enabled.
type Options struct {
	RequireUppercase bool // keys must be ALL_CAPS
	DisallowEmpty    bool // values must not be empty
	DisallowSpaces   bool // keys must not contain spaces
}

// DefaultOptions returns a sensible default set of lint rules.
func DefaultOptions() Options {
	return Options{
		RequireUppercase: true,
		DisallowEmpty:    true,
		DisallowSpaces:   true,
	}
}

// Lint checks the provided env map against the configured options and returns
// any issues found. The returned slice is empty when no issues are detected.
func Lint(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for k, v := range env {
		if opts.DisallowSpaces && strings.ContainsFunc(k, unicode.IsSpace) {
			issues = append(issues, Issue{
				Key:     k,
				Message: "key contains whitespace",
			})
		}

		if opts.RequireUppercase && k != strings.ToUpper(k) {
			issues = append(issues, Issue{
				Key:     k,
				Message: "key is not uppercase",
			})
		}

		if opts.DisallowEmpty && strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{
				Key:     k,
				Message: "value is empty",
			})
		}
	}

	return issues
}
