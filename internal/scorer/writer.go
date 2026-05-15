package scorer

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes a human-readable score report to w.
func WriteText(w io.Writer, r Result) error {
	if w == nil {
		return fmt.Errorf("scorer: writer is nil")
	}
	if r.Score == 0 && len(r.Deductions) == 0 {
		_, err := fmt.Fprintln(w, "Score: 100 / 100  Grade: A  (no issues)")
		return err
	}

	_, err := fmt.Fprintf(w, "Score: %.0f / 100  Grade: %s\n", r.Score, r.Grade)
	if err != nil {
		return err
	}
	if len(r.Deductions) == 0 {
		_, err = fmt.Fprintln(w, "  No deductions.")
		return err
	}
	_, err = fmt.Fprintln(w, "Deductions:")
	if err != nil {
		return err
	}
	for _, d := range r.Deductions {
		_, err = fmt.Fprintf(w, "  %-30s  %-20s  -%.0f pts\n", d.Key, d.Reason, d.Points)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes the result as compact JSON to w.
func WriteJSON(w io.Writer, r Result) error {
	if w == nil {
		return fmt.Errorf("scorer: writer is nil")
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
