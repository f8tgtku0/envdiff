// Package templater provides functionality to generate a .env.template file
// from one or more parsed environment maps, replacing values with empty strings
// or placeholder comments to indicate required variables.
package templater

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls how the template is rendered.
type Options struct {
	// Placeholder is the value written for each key in the template.
	// Defaults to an empty string if not set.
	Placeholder string
	// CommentValues, when true, appends the original value as an inline comment.
	CommentValues bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Placeholder:   "",
		CommentValues: false,
	}
}

// Generate writes a .env template to w derived from the provided env map.
// Keys are sorted alphabetically. Values are replaced by opts.Placeholder.
func Generate(w io.Writer, env map[string]string, opts Options) error {
	if w == nil {
		return fmt.Errorf("templater: writer must not be nil")
	}

	keys := sortedKeys(env)

	for _, k := range keys {
		original := env[k]
		line := fmt.Sprintf("%s=%s", k, opts.Placeholder)
		if opts.CommentValues && strings.TrimSpace(original) != "" {
			line += fmt.Sprintf(" # was: %s", original)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("templater: write error: %w", err)
		}
	}
	return nil
}

// Merge combines multiple env maps into a single deduplicated key set.
// When a key appears in multiple maps the first non-empty value is kept
// solely for CommentValues rendering; the template output itself uses the placeholder.
func Merge(envs ...map[string]string) map[string]string {
	out := make(map[string]string)
	for _, env := range envs {
		for k, v := range env {
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
