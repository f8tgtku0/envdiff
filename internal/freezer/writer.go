package freezer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// WriteText writes a human-readable report of violations to w.
func WriteText(r Result, w io.Writer) error {
	if w == nil {
		return errors.New("freezer: writer must not be nil")
	}
	if r.Clean() {
		_, err := fmt.Fprintln(w, "frozen env: no violations found")
		return err
	}
	fmt.Fprintf(w, "frozen env: %d violation(s) found\n", len(r.Violations))
	for _, v := range r.Violations {
		fmt.Fprintf(w, "  [%s] %s\n", v.Key, v.Reason)
	}
	return nil
}

// WriteJSON writes the result as a JSON object to w.
func WriteJSON(r Result, w io.Writer) error {
	if w == nil {
		return errors.New("freezer: writer must not be nil")
	}
	type jsonViolation struct {
		Key    string `json:"key"`
		Reason string `json:"reason"`
	}
	type jsonResult struct {
		Clean      bool            `json:"clean"`
		Violations []jsonViolation `json:"violations"`
	}
	out := jsonResult{Clean: r.Clean(), Violations: []jsonViolation{}}
	for _, v := range r.Violations {
		out.Violations = append(out.Violations, jsonViolation{Key: v.Key, Reason: v.Reason})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
