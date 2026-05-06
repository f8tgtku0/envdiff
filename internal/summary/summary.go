// Package summary provides aggregation and statistics over diff results
// produced by the differ package.
package summary

import (
	"fmt"
	"strings"

	"github.com/yourorg/envdiff/internal/differ"
)

// Stats holds aggregated counts from a diff result.
type Stats struct {
	TotalKeys   int
	Clean       int
	MissingLeft int
	MissingRight int
	Mismatched  int
}

// Compute derives Stats from a differ.Result.
func Compute(r differ.Result) Stats {
	return Stats{
		TotalKeys:    r.TotalKeys(),
		Clean:        r.CleanCount(),
		MissingLeft:  len(r.MissingInLeft),
		MissingRight: len(r.MissingInRight),
		Mismatched:   len(r.Mismatched),
	}
}

// HasDiff returns true when any discrepancies exist.
func (s Stats) HasDiff() bool {
	return s.MissingLeft > 0 || s.MissingRight > 0 || s.Mismatched > 0
}

// String returns a human-readable one-line summary.
func (s Stats) String() string {
	parts := []string{
		fmt.Sprintf("total=%d", s.TotalKeys),
		fmt.Sprintf("clean=%d", s.Clean),
		fmt.Sprintf("missing_left=%d", s.MissingLeft),
		fmt.Sprintf("missing_right=%d", s.MissingRight),
		fmt.Sprintf("mismatched=%d", s.Mismatched),
	}
	return strings.Join(parts, " ")
}
