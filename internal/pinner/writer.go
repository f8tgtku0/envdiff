package pinner

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes a human-readable summary of pin results to w.
func WriteText(w io.Writer, r *Result) error {
	if w == nil {
		return fmt.Errorf("pinner: writer must not be nil")
	}
	if r == nil {
		_, err := fmt.Fprintln(w, "no pin result")
		return err
	}

	pinned := r.Pinned()
	skipped := r.Skipped()

	if len(pinned) == 0 && len(skipped) == 0 {
		_, err := fmt.Fprintln(w, "nothing to pin")
		return err
	}

	for _, e := range pinned {
		if e.OldValue == "" && e.Reason == "added" {
			_, err := fmt.Fprintf(w, "PINNED   %s = %q (added)\n", e.Key, e.NewValue)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(w, "PINNED   %s: %q -> %q\n", e.Key, e.OldValue, e.NewValue)
			if err != nil {
				return err
			}
		}
	}

	for _, e := range skipped {
		_, err := fmt.Fprintf(w, "SKIPPED  %s (%s)\n", e.Key, e.Reason)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, "\nSummary: %d pinned, %d skipped\n", len(pinned), len(skipped))
	return err
}

// WriteJSON writes pin results as a JSON array to w.
func WriteJSON(w io.Writer, r *Result) error {
	if w == nil {
		return fmt.Errorf("pinner: writer must not be nil")
	}
	if r == nil {
		_, err := fmt.Fprintln(w, "[]")
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r.Entries)
}
