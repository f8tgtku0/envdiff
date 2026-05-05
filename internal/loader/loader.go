// Package loader provides functionality for loading and merging multiple
// .env files, applying filters, and preparing data for diffing.
package loader

import (
	"fmt"

	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/parser"
)

// Options controls how files are loaded.
type Options struct {
	// IncludePrefixes limits keys to those matching any of these prefixes.
	IncludePrefixes []string
	// ExcludePrefixes omits keys matching any of these prefixes.
	ExcludePrefixes []string
}

// LoadedEnv holds the parsed and filtered contents of a single .env file.
type LoadedEnv struct {
	// Path is the original file path that was loaded.
	Path string
	// Vars contains the key/value pairs after filtering.
	Vars map[string]string
}

// Load reads a .env file from the given path, parses it, and applies any
// filters specified in opts. It returns a LoadedEnv or an error.
func Load(path string, opts Options) (*LoadedEnv, error) {
	vars, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("loader: failed to parse %q: %w", path, err)
	}

	filtered := filter.Apply(vars, filter.Options{
		IncludePrefixes: opts.IncludePrefixes,
		ExcludePrefixes: opts.ExcludePrefixes,
	})

	return &LoadedEnv{
		Path: path,
		Vars: filtered,
	}, nil
}

// LoadPair loads two .env files with the same options and returns them in
// order. It is a convenience wrapper around Load for the common diff case.
func LoadPair(leftPath, rightPath string, opts Options) (*LoadedEnv, *LoadedEnv, error) {
	left, err := Load(leftPath, opts)
	if err != nil {
		return nil, nil, err
	}

	right, err := Load(rightPath, opts)
	if err != nil {
		return nil, nil, err
	}

	return left, right, nil
}
