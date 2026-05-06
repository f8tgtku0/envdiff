// Package summary aggregates statistics from a differ.Result, providing
// counts of clean, missing, and mismatched keys across two .env files.
//
// Typical usage:
//
//	result := differ.Compare(left, right)
//	stats := summary.Compute(result)
//	if stats.HasDiff() {
//		fmt.Println(stats)
//	}
package summary
