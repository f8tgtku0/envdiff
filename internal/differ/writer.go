package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable diff report to w.
// It returns an error if writing fails.
func WriteText(r *Result, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("differ: writer must not be nil")
	}
	if r == nil {
		return fmt.Errorf("differ: result must not be nil")
	}

	if len(r.MissingInRight) == 0 && len(r.MissingInLeft) == 0 && len(r.Mismatched) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	if len(r.MissingInRight) > 0 {
		fmt.Fprintln(w, "Missing in right:")
		keys := sortedStringSlice(r.MissingInRight)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.MissingInLeft) > 0 {
		fmt.Fprintln(w, "Missing in left:")
		keys := sortedStringSlice(r.MissingInLeft)
		for _, k := range keys {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(r.Mismatched) > 0 {
		fmt.Fprintln(w, "Mismatched values:")
		keys := make([]string, 0, len(r.Mismatched))
		for k := range r.Mismatched {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			p := r.Mismatched[k]
			fmt.Fprintf(w, "  ~ %s: %q != %q\n", k, p.Left, p.Right)
		}
	}

	return nil
}

// WriteJSON writes the diff result as JSON to w.
func WriteJSON(r *Result, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("differ: writer must not be nil")
	}
	if r == nil {
		return fmt.Errorf("differ: result must not be nil")
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func sortedStringSlice(ss []string) []string {
	out := make([]string, len(ss))
	copy(out, ss)
	sort.Strings(out)
	return out
}
