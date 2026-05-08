package interpolator

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches ${VAR} and $VAR style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// Strict causes Interpolate to return an error when a referenced variable
	// is not found in the environment map.
	Strict bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{Strict: false}
}

// Interpolate resolves variable references inside env values using the
// provided env map as the source of truth.  It returns a new map with
// resolved values and leaves the original untouched.
//
// References can be written as $VAR or ${VAR}.  Self-referential or
// circular references are left unexpanded when Strict is false; in
// Strict mode they produce an error.
func Interpolate(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := resolve(k, v, env, opts, map[string]bool{})
		if err != nil {
			return nil, err
		}
		out[k] = resolved
	}
	return out, nil
}

// resolve expands a single value string, tracking visited keys to detect
// circular references.
func resolve(key, value string, env map[string]string, opts Options, visited map[string]bool) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		ref := extractName(match)
		if ref == key || visited[ref] {
			// circular or self-reference
			if opts.Strict {
				expandErr = fmt.Errorf("interpolator: circular reference detected for key %q", ref)
			}
			return match
		}
		val, ok := env[ref]
		if !ok {
			if opts.Strict {
				expandErr = fmt.Errorf("interpolator: undefined variable %q referenced in key %q", ref, key)
			}
			return match
		}
		nextVisited := copyVisited(visited)
		nextVisited[key] = true
		expanded, err := resolve(ref, val, env, opts, nextVisited)
		if err != nil {
			expandErr = err
			return match
		}
		return expanded
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}

func copyVisited(src map[string]bool) map[string]bool {
	dst := make(map[string]bool, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
