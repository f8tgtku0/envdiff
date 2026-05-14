package typecheck

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable report of type issues to w.
func WriteText(issues []Issue, w io.Writer) error {
	if len(issues) == 0 {
		_, err := fmt.Fprintln(w, "typecheck: all values match expected types")
		return err
	}

	sorted := make([]Issue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	_, err := fmt.Fprintf(w, "typecheck: %d issue(s) found\n", len(sorted))
	if err != nil {
		return err
	}
	for _, iss := range sorted {
		if _, err := fmt.Fprintf(w, "  [%s] %s\n", iss.Expected, iss.String()); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes a JSON array of type issues to w.
func WriteJSON(issues []Issue, w io.Writer) error {
	type jsonIssue struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		Expected string `json:"expected_type"`
		Reason   string `json:"reason"`
	}

	out := make([]jsonIssue, 0, len(issues))
	for _, iss := range issues {
		out = append(out, jsonIssue{
			Key:      iss.Key,
			Value:    iss.Value,
			Expected: string(iss.Expected),
			Reason:   iss.Reason,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
