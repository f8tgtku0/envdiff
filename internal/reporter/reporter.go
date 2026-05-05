// Package reporter formats and outputs the diff results from envdiff comparisons.
package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/differ"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the diff result and metadata for rendering.
type Report struct {
	Left  string
	Right string
	Diff  differ.DiffResult
}

// NewReport creates a new Report with the given file labels and diff result.
func NewReport(left, right string, diff differ.DiffResult) *Report {
	return &Report{
		Left:  left,
		Right: right,
		Diff:  diff,
	}
}

// Write renders the report in the specified format to the given writer.
func (r *Report) Write(w io.Writer, format Format) error {
	switch format {
	case FormatJSON:
		return r.writeJSON(w)
	default:
		return r.writeText(w)
	}
}

func (r *Report) writeText(w io.Writer) error {
	if r.Diff.IsClean() {
		_, err := fmt.Fprintf(w, "No differences found between %s and %s.\n", r.Left, r.Right)
		return err
	}

	fmt.Fprintf(w, "Comparing %s <-> %s\n", r.Left, r.Right)
	fmt.Fprintln(w, "")

	if len(r.Diff.MissingInRight) > 0 {
		keys := sortedKeys(r.Diff.MissingInRight)
		fmt.Fprintf(w, "Missing in %s:\n", r.Right)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
		fmt.Fprintln(w, "")
	}

	if len(r.Diff.MissingInLeft) > 0 {
		keys := sortedKeys(r.Diff.MissingInLeft)
		fmt.Fprintf(w, "Missing in %s:\n", r.Left)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
		fmt.Fprintln(w, "")
	}

	if len(r.Diff.Mismatched) > 0 {
		keys := sortedKeys(r.Diff.Mismatched)
		fmt.Fprintln(w, "Mismatched values:")
		for _, k := range keys {
			m := r.Diff.Mismatched[k]
			fmt.Fprintf(w, "  ~ %s: %q (%s) vs %q (%s)\n", k, m.LeftVal, r.Left, m.RightVal, r.Right)
		}
	}

	return nil
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
