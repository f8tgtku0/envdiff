package masker

import (
	"fmt"
	"strings"
)

// Options controls how values are masked in an env map.
type Options struct {
	// Keys is the explicit list of keys whose values should be masked.
	Keys []string
	// Prefixes masks any key that starts with one of these prefixes.
	Prefixes []string
	// Mask is the string used to replace sensitive values. Defaults to "***".
	Mask string
	// ShowLength appends the original value length hint, e.g. "***[8]".
	ShowLength bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Prefixes:   []string{"SECRET_", "PRIVATE_", "TOKEN_", "PASSWORD_", "PASS_"},
		Mask:       "***",
		ShowLength: false,
	}
}

// Mask returns a new map with sensitive values replaced by the mask string.
// The original map is never mutated.
func Mask(env map[string]string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, fmt.Errorf("masker: env map must not be nil")
	}
	if opts.Mask == "" {
		opts.Mask = "***"
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if shouldMask(k, keySet, opts.Prefixes) {
			out[k] = buildMask(opts.Mask, v, opts.ShowLength)
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func shouldMask(key string, keySet map[string]struct{}, prefixes []string) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

func buildMask(mask, original string, showLength bool) string {
	if showLength {
		return fmt.Sprintf("%s[%d]", mask, len(original))
	}
	return mask
}
