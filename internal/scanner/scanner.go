// Package scanner provides utilities for scanning directories and
// collecting .env files matching configurable glob patterns.
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options configures the behaviour of a directory scan.
type Options struct {
	// Patterns is a list of glob patterns to match (e.g. ".env*", "*.env").
	// Defaults to [".env*"] when empty.
	Patterns []string

	// Recursive controls whether sub-directories are traversed.
	Recursive bool

	// IgnoreDirs is a set of directory names to skip (e.g. ".git", "vendor").
	IgnoreDirs []string
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		Patterns:   []string{".env*", "*.env"},
		Recursive:  false,
		IgnoreDirs: []string{".git", "vendor", "node_modules"},
	}
}

// Scan walks root (or just its top level when opts.Recursive is false)
// and returns all file paths that match at least one of opts.Patterns,
// sorted lexicographically.
func Scan(root string, opts Options) ([]string, error) {
	if len(opts.Patterns) == 0 {
		opts.Patterns = DefaultOptions().Patterns
	}

	ignored := make(map[string]bool, len(opts.IgnoreDirs))
	for _, d := range opts.IgnoreDirs {
		ignored[d] = true
	}

	var results []string

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("scanner: accessing %q: %w", path, err)
		}

		if info.IsDir() {
			if path != root && (ignored[info.Name()] || !opts.Recursive) {
				return filepath.SkipDir
			}
			return nil
		}

		name := filepath.Base(path)
		for _, pattern := range opts.Patterns {
			matched, err := filepath.Match(pattern, name)
			if err != nil {
				return fmt.Errorf("scanner: bad pattern %q: %w", pattern, err)
			}
			if matched {
				results = append(results, filepath.ToSlash(path))
				break
			}
		}
		return nil
	}

	if err := filepath.Walk(root, walkFn); err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return strings.ToLower(results[i]) < strings.ToLower(results[j])
	})

	return results, nil
}
