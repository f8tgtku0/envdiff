package reporter

import (
	"encoding/json"
	"io"
)

type jsonMismatch struct {
	Key      string `json:"key"`
	LeftVal  string `json:"left_value"`
	RightVal string `json:"right_value"`
}

type jsonReport struct {
	Left          string         `json:"left"`
	Right         string         `json:"right"`
	MissingInLeft []string       `json:"missing_in_left"`
	MissingInRight []string      `json:"missing_in_right"`
	Mismatched    []jsonMismatch `json:"mismatched"`
	Clean         bool           `json:"clean"`
}

func (r *Report) writeJSON(w io.Writer) error {
	missingLeft := sortedKeys(r.Diff.MissingInLeft)
	missingRight := sortedKeys(r.Diff.MissingInRight)

	mismatchedKeys := sortedKeys(r.Diff.Mismatched)
	mismatches := make([]jsonMismatch, 0, len(mismatchedKeys))
	for _, k := range mismatchedKeys {
		m := r.Diff.Mismatched[k]
		mismatches = append(mismatches, jsonMismatch{
			Key:      k,
			LeftVal:  m.LeftVal,
			RightVal: m.RightVal,
		})
	}

	if missingLeft == nil {
		missingLeft = []string{}
	}
	if missingRight == nil {
		missingRight = []string{}
	}
	if mismatches == nil {
		mismatches = []jsonMismatch{}
	}

	out := jsonReport{
		Left:           r.Left,
		Right:          r.Right,
		MissingInLeft:  missingLeft,
		MissingInRight: missingRight,
		Mismatched:     mismatches,
		Clean:          r.Diff.IsClean(),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
