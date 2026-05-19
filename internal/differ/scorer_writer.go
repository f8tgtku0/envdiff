package differ

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteSimilarityText writes a human-readable similarity report to w.
func WriteSimilarityText(w io.Writer, r SimilarityReport) error {
	if w == nil {
		return fmt.Errorf("differ: writer must not be nil")
	}

	grade := similarityGrade(r.Score)
	fmt.Fprintf(w, "Similarity Score : %.1f%%  [%s]\n", r.Score*100, grade)
	fmt.Fprintf(w, "Total Keys       : %d\n", r.TotalKeys)
	if r.MissingLeft > 0 {
		fmt.Fprintf(w, "Missing in left  : %d\n", r.MissingLeft)
	}
	if r.MissingRight > 0 {
		fmt.Fprintf(w, "Missing in right : %d\n", r.MissingRight)
	}
	if r.Mismatched > 0 {
		fmt.Fprintf(w, "Mismatched       : %d\n", r.Mismatched)
	}
	if r.MissingLeft == 0 && r.MissingRight == 0 && r.Mismatched == 0 {
		fmt.Fprintln(w, "No differences found.")
	}
	return nil
}

// WriteSimilarityJSON writes a JSON-encoded similarity report to w.
func WriteSimilarityJSON(w io.Writer, r SimilarityReport) error {
	if w == nil {
		return fmt.Errorf("differ: writer must not be nil")
	}
	type payload struct {
		Score        float64  `json:"score"`
		Grade        string   `json:"grade"`
		TotalKeys    int      `json:"total_keys"`
		MissingLeft  int      `json:"missing_left"`
		MissingRight int      `json:"missing_right"`
		Mismatched   int      `json:"mismatched"`
		Keys         []string `json:"keys"`
	}
	p := payload{
		Score:        r.Score,
		Grade:        similarityGrade(r.Score),
		TotalKeys:    r.TotalKeys,
		MissingLeft:  r.MissingLeft,
		MissingRight: r.MissingRight,
		Mismatched:   r.Mismatched,
		Keys:         r.SortedKeys,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func similarityGrade(score float64) string {
	switch {
	case score >= 0.95:
		return "excellent"
	case score >= 0.75:
		return "good"
	case score >= 0.50:
		return "fair"
	default:
		return "poor"
	}
}
