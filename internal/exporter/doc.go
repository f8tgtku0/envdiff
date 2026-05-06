// Package exporter serialises a map of environment variables to various
// text formats suitable for use in shell scripts or .env files.
//
// Supported formats:
//
//   - dotenv  – plain KEY=VALUE pairs (default)
//   - shell   – same as dotenv, suitable for sourcing in POSIX shells
//   - export  – prefixes each line with "export" for direct shell evaluation
//
// Values that contain whitespace, hash characters, dollar signs, or newlines
// are automatically double-quoted.  Output keys are always written in
// lexicographic order for deterministic diffs.
package exporter
