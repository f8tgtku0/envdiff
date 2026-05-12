// Package masker provides utilities for masking sensitive values in env maps.
//
// It replaces the values of keys that match explicit names or key prefixes with
// a configurable mask string (default "***"). The original map is never mutated;
// Mask always returns a new copy.
//
// Example usage:
//
//	opts := masker.DefaultOptions()
//	opts.ShowLength = true
//	masked, err := masker.Mask(env, opts)
//
The package is intentionally separate from the redactor package: redactor
operates on pattern-matched key names using regular expressions, while masker
offers a simpler prefix/key-list approach suited for quick CLI output sanitisation.
package masker
