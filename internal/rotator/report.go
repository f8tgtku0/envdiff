package rotator

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteReport writes a human-readable rotation report to w.
func WriteReport(w io.Writer, res *Result) error {
	if res == nil {
		return fmt.Errorf("rotator: result must not be nil")
	}

	rotatedKeys := sortedKeys(res.Rotated)
	skipped := append([]string(nil), res.Skipped...)
	sort.Strings(skipped)

	if len(rotatedKeys) == 0 && len(skipped) == 0 {
		_, err := fmt.Fprintln(w, "No keys rotated.")
		return err
	}

	var sb strings.Builder
	if len(rotatedKeys) > 0 {
		sb.WriteString(fmt.Sprintf("Rotated (%d):\n", len(rotatedKeys)))
		for _, k := range rotatedKeys {
			archiveKey := k + "_OLD"
			if _, ok := res.Archived[archiveKey]; ok {
				sb.WriteString(fmt.Sprintf("  %-30s  (archived as %s)\n", k, archiveKey))
			} else {
				sb.WriteString(fmt.Sprintf("  %s\n", k))
			}
		}
	}

	if len(skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped (%d):\n", len(skipped)))
		for _, k := range skipped {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}

// Summarise returns a one-line summary string.
func Summarise(res *Result) string {
	if res == nil {
		return "rotator: nil result"
	}
	return fmt.Sprintf("%d rotated, %d archived, %d skipped",
		len(res.Rotated), len(res.Archived), len(res.Skipped))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
