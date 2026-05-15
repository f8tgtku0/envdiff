// Package digester computes and compares content hashes for env maps,
// allowing callers to detect whether an environment has changed between
// two points in time without performing a full field-by-field diff.
package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// Algorithm selects the hash function used when computing a digest.
type Algorithm string

const (
	SHA256 Algorithm = "sha256"
)

// Options controls Digest behaviour.
type Options struct {
	// Algorithm is the hash algorithm to use. Defaults to SHA256.
	Algorithm Algorithm
	// IncludeKeys, when true, mixes key names into the hash so that a
	// rename is detected even when all values remain the same.
	IncludeKeys bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		Algorithm:   SHA256,
		IncludeKeys: true,
	}
}

// Result holds the digest of a single env map.
type Result struct {
	Hex       string
	Algorithm Algorithm
	KeyCount  int
}

// Digest computes a deterministic hash of env using opts.
// Keys are sorted before hashing to ensure stability.
func Digest(env map[string]string, opts Options) (Result, error) {
	if env == nil {
		return Result{}, fmt.Errorf("digester: env must not be nil")
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		if opts.IncludeKeys {
			fmt.Fprintf(h, "%s=%s\n", k, env[k])
		} else {
			fmt.Fprintf(h, "%s\n", env[k])
		}
	}

	return Result{
		Hex:       hex.EncodeToString(h.Sum(nil)),
		Algorithm: opts.Algorithm,
		KeyCount:  len(keys),
	}, nil
}

// Equal returns true when two Results share the same digest hex.
func Equal(a, b Result) bool {
	return a.Hex == b.Hex
}
