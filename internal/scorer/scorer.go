package scorer

import (
	"math"
	"sort"
	"strings"
)

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		MissingPenalty:   10.0,
		EmptyPenalty:     3.0,
		MismatchPenalty:  5.0,
		MaxScore:         100.0,
	}
}

// Options controls how the health score is computed.
type Options struct {
	MissingPenalty  float64
	EmptyPenalty    float64
	MismatchPenalty float64
	MaxScore        float64
}

// Result holds the computed score and a breakdown of deductions.
type Result struct {
	Score      float64
	Deductions []Deduction
	Grade      string
}

// Deduction describes a single penalty applied during scoring.
type Deduction struct {
	Key    string
	Reason string
	Points float64
}

// Score evaluates the health of an env map against a reference set of
// expected keys and returns a Result with a 0–100 score.
func Score(env map[string]string, expected []string, opts Options) Result {
	var deductions []Deduction
	total := 0.0

	for _, key := range expected {
		val, ok := env[key]
		if !ok {
			deductions = append(deductions, Deduction{Key: key, Reason: "missing", Points: opts.MissingPenalty})
			total += opts.MissingPenalty
			continue
		}
		if strings.TrimSpace(val) == "" {
			deductions = append(deductions, Deduction{Key: key, Reason: "empty value", Points: opts.EmptyPenalty})
			total += opts.EmptyPenalty
		}
	}

	sort.Slice(deductions, func(i, j int) bool {
		return deductions[i].Key < deductions[j].Key
	})

	raw := opts.MaxScore - total
	score := math.Max(0, math.Min(opts.MaxScore, raw))

	return Result{
		Score:      score,
		Deductions: deductions,
		Grade:      grade(score),
	}
}

// ScorePair computes a score that also penalises mismatched values between
// two env maps (e.g. staging vs production).
func ScorePair(left, right map[string]string, expected []string, opts Options) Result {
	res := Score(left, expected, opts)
	extra := 0.0

	for _, key := range expected {
		lv, lok := left[key]
		rv, rok := right[key]
		if lok && rok && lv != rv {
			d := Deduction{Key: key, Reason: "value mismatch", Points: opts.MismatchPenalty}
			res.Deductions = append(res.Deductions, d)
			extra += opts.MismatchPenalty
		}
	}

	res.Score = math.Max(0, math.Min(opts.MaxScore, res.Score-extra))
	res.Grade = grade(res.Score)

	sort.Slice(res.Deductions, func(i, j int) bool {
		return res.Deductions[i].Key < res.Deductions[j].Key
	})
	return res
}

func grade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
