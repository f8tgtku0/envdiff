package annotator

import (
	"fmt"
	"io"
	"sort"
)

// WriteText writes the annotated env to w in dotenv format, placing each
// annotation comment on the line immediately above its key.
func WriteText(r *Result, w io.Writer) error {
	if r == nil {
		return fmt.Errorf("annotator: result must not be nil")
	}

	keys := make([]string, 0, len(r.Env))
	for k := range r.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if ann, ok := r.Annotations[k]; ok {
			if _, err := fmt.Fprintln(w, ann.Comment); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, r.Env[k]); err != nil {
			return err
		}
	}

	return nil
}

// WriteSummary writes a short human-readable summary of the annotation run.
func WriteSummary(r *Result, w io.Writer) error {
	if r == nil {
		return fmt.Errorf("annotator: result must not be nil")
	}

	annotated := len(r.Annotations)
	skipped := len(r.Skipped)

	_, err := fmt.Fprintf(w,
		"annotated: %d key(s), skipped: %d key(s)\n",
		annotated, skipped,
	)
	return err
}
