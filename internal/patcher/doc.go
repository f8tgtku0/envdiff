// Package patcher applies a patch env map on top of a base env map,
// producing a merged result together with a structured change log.
//
// Basic usage:
//
//	base  := map[string]string{"HOST": "localhost", "PORT": "5432"}
//	patch := map[string]string{"PORT": "6543", "DEBUG": "true"}
//
//	opts := patcher.DefaultOptions()
//	opts.Overwrite = true
//
//	result, changes, err := patcher.Apply(base, patch, opts)
//
// The returned changes slice describes every key that was added, updated,
// or skipped. Use WriteReport to render a human-readable diff and Summarise
// to obtain aggregate counts.
package patcher
