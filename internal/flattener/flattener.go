// Package flattener collapses nested prefix hierarchies in an env map into
// a single-level map using a configurable separator.
package flattener

import (
	"errors"
	"sort"
	"strings"
)

// DefaultOptions returns a set of sensible defaults for Flatten.
func DefaultOptions() Options {
	return Options{
		Separator: "__",
		Lowercase:  false,
		StripDepth: 0,
	}
}

// Options controls the behaviour of Flatten.
type Options struct {
	// Separator is the string used to join prefix segments (default "__").
	Separator string

	// Lowercase converts all resulting keys to lower-case.
	Lowercase bool

	// StripDepth removes the first N prefix segments from each key.
	// A value of 0 keeps all segments.
	StripDepth int
}

// Flatten takes a map whose keys may contain Separator-delimited segments and
// returns a new map with those segments normalised.  When StripDepth > 0 the
// leading N segments are removed from every key.  If two keys collide after
// stripping, the last value (in sorted key order) wins.
//
// Flatten returns an error if env is nil or if Separator is empty.
func Flatten(env map[string]string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("flattener: env must not be nil")
	}
	if opts.Separator == "" {
		return nil, errors.New("flattener: separator must not be empty")
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make(map[string]string, len(env))
	for _, k := range keys {
		v := env[k]
		newKey := flatten(k, opts)
		out[newKey] = v
	}
	return out, nil
}

func flatten(key string, opts Options) string {
	segments := strings.Split(key, opts.Separator)
	if opts.StripDepth > 0 && opts.StripDepth < len(segments) {
		segments = segments[opts.StripDepth:]
	}
	result := strings.Join(segments, opts.Separator)
	if opts.Lowercase {
		result = strings.ToLower(result)
	}
	return result
}
