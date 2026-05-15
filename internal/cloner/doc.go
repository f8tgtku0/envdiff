// Package cloner copies env maps with optional key-prefix transformations.
//
// # Overview
//
// Clone duplicates a source env map into a destination, supporting:
//
//   - AddPrefix  — prepend a string to every output key
//   - StripPrefix — remove a leading string from every key before output
//   - StrictPrefix — treat missing prefix as a hard error
//   - Overwrite — control whether existing destination keys are replaced
//
// # Example
//
//	opts := cloner.DefaultOptions()
//	opts.StripPrefix = "APP_"
//	opts.AddPrefix  = "PROD_"
//	out, err := cloner.Clone(src, nil, opts)
package cloner
