// Package pinner provides functionality for locking (pinning) specific
// environment variable keys to fixed values.
//
// Pin walks a set of desired pin mappings and applies them to an env map,
// recording whether each key was pinned, skipped due to an existing value,
// or added fresh. The caller controls overwrite behaviour and whether a
// missing key should be treated as an error via Options.
//
// Example usage:
//
//	env := map[string]string{"APP_ENV": "staging", "DEBUG": "true"}
//	pins := map[string]string{"APP_ENV": "production"}
//	result, err := pinner.Pin(env, pins, pinner.DefaultOptions())
package pinner
