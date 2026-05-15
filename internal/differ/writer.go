package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable diff report to w.
func WriteText(r *Result, w io.Writer) error {
	if r == nil || w == nil {
		return fmt.Errorf("differ: nil argument")
	}

	if r.Clean() {
		fmt.Fprintln(w, "No differences found.")
		return nil
	}

	if len(r.MissingInRight) > 0 {
		fmt.Fprintln(w, "Missing in right:")
		for _, k := range sortedStringSlice(r.MissingInRight) {
			fmt.Fprintf(w, "  - %s=%s\n", k, r.MissingInRight[k])
		}
	}

	if len(r.MissingInLeft) > 0 {
		fmt.Fprintln(w, "Missing in left:")
		for _, k := range sortedStringSlice(r.MissingInLeft) {
			fmt.Fprintf(w, "  - %s=%s\n", k, r.MissingInLeft[k])
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

// WriteJSON writes the diff result as a JSON object to w.
func WriteJSON(r *Result, w io.Writer) error {
	if r == nil || w == nil {
		return fmt.Errorf("differ: nil argument")
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func sortedStringSlice(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
