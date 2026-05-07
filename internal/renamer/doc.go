// Package renamer provides utilities for renaming keys in an environment map.
//
// It accepts a RenameMap (old key -> new key) and applies the renames to a
// copy of the source env, leaving the original untouched. The caller can
// control conflict behaviour (when a target key already exists) via Options.
//
// Typical usage:
//
//	out, result, err := renamer.Rename(env, renamer.RenameMap{
//		"DB_HOST": "DATABASE_HOST",
//	}, renamer.DefaultOptions())
//
// Keys not present in the source env are recorded in Result.Skipped.
// Conflicts are recorded in Result.Conflicts; by default a conflict returns
// an error unless Options.OverwriteConflicts is true.
package renamer
