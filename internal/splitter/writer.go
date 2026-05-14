package splitter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable representation of the split result to w.
func WriteText(r *Result, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("splitter: writer must not be nil")
	}
	if r == nil {
		_, err := fmt.Fprintln(w, "(no result)")
		return err
	}

	prefixes := sortedBucketKeys(r.Buckets)
	for _, p := range prefixes {
		label := p
		if label == "" {
			label = "(unmatched)"
		}
		fmt.Fprintf(w, "[%s]\n", label) //nolint:errcheck
		bucket := r.Buckets[p]
		keys := make([]string, 0, len(bucket))
		for k := range bucket {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "  %s=%s\n", k, bucket[k]) //nolint:errcheck
		}
	}
	return nil
}

// WriteJSON writes the split result as a JSON object to w.
func WriteJSON(r *Result, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("splitter: writer must not be nil")
	}
	if r == nil {
		_, err := fmt.Fprintln(w, "{}")
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r.Buckets)
}

func sortedBucketKeys(m map[string]map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
