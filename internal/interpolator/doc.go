// Package interpolator resolves variable references within .env file values.
//
// It supports both $VAR and ${VAR} syntax, expanding references using the
// keys present in the same environment map.  Chained references (where a
// value itself references another variable) are fully resolved in a single
// pass.
//
// Two modes are available:
//
//   - Lenient (default): unresolved or circular references are left as-is.
//   - Strict: unresolved or circular references return an error immediately.
//
// The input map is never mutated; Interpolate always returns a new map.
package interpolator
