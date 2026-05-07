package patcher

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// WriteReport prints a human-readable summary of changes to w.
func WriteReport(w io.Writer, changes []Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "No changes.")
		return err
	}

	// Sort for deterministic output.
	sorted := make([]Change, len(changes))
	copy(sorted, changes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "ACTION\tKEY\tOLD\tNEW")
	_, _ = fmt.Fprintln(tw, "------\t---\t---\t---")
	for _, c := range sorted {
		old := c.OldValue
		if old == "" {
			old = "<none>"
		}
		_, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", c.Action, c.Key, old, c.NewValue)
		if err != nil {
			return err
		}
	}
	return tw.Flush()
}

// Stats summarises a slice of changes by action type.
type Stats struct {
	Added   int
	Updated int
	Skipped int
}

// Summarise counts changes by action.
func Summarise(changes []Change) Stats {
	var s Stats
	for _, c := range changes {
		switch c.Action {
		case "add":
			s.Added++
		case "update":
			s.Updated++
		case "skip":
			s.Skipped++
		}
	}
	return s
}

// String returns a one-line summary.
func (s Stats) String() string {
	return fmt.Sprintf("added=%d updated=%d skipped=%d", s.Added, s.Updated, s.Skipped)
}
