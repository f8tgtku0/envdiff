package comparator

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteText writes a human-readable comparison report to w.
func WriteText(results []Result, labels []string, w io.Writer) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "All keys are consistent across environments.")
		return err
	}

	sortedLabels := make([]string, len(labels))
	copy(sortedLabels, labels)
	sort.Strings(sortedLabels)

	header := fmt.Sprintf("%-30s  %s\n", "KEY", strings.Join(sortedLabels, "  "))
	if _, err := fmt.Fprint(w, header); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, strings.Repeat("-", len(header))); err != nil {
		return err
	}

	for _, r := range results {
		mark := "✓"
		if !r.Consistent {
			mark = "✗"
		}
		var cols []string
		for _, l := range sortedLabels {
			v, ok := r.Values[l]
			if !ok {
				v = "<missing>"
			}
			cols = append(cols, v)
		}
		line := fmt.Sprintf("%s %-28s  %s\n", mark, r.Key, strings.Join(cols, "  "))
		if _, err := fmt.Fprint(w, line); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes the comparison results as a JSON array to w.
func WriteJSON(results []Result, w io.Writer) error {
	type jsonRow struct {
		Key        string            `json:"key"`
		Consistent bool              `json:"consistent"`
		Missing    []string          `json:"missing,omitempty"`
		Values     map[string]string `json:"values"`
	}

	rows := make([]jsonRow, len(results))
	for i, r := range results {
		rows[i] = jsonRow{
			Key:        r.Key,
			Consistent: r.Consistent,
			Missing:    r.Missing,
			Values:     r.Values,
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
