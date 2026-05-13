// Package rotator provides utilities for rotating secret values inside .env
// environments.
//
// A rotation replaces a key's current value with a new one supplied by the
// caller, while preserving the old value under an archive key (by default
// KEY_OLD). This makes it straightforward to roll back a rotation if the new
// credential turns out to be invalid.
//
// Basic usage:
//
//	env := map[string]string{"DB_PASSWORD": "old"}
//	replacements := map[string]string{"DB_PASSWORD": "new"}
//	res, err := rotator.Rotate(env, replacements, rotator.DefaultOptions())
//	updated := rotator.Apply(env, res, "")
//
// The Result type records which keys were rotated, which archive entries were
// created, and which keys were skipped (e.g. due to archive-key conflicts).
package rotator
