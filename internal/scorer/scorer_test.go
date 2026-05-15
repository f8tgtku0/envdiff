package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/scorer"
)

var expected = []string{"DB_HOST", "DB_PORT", "API_KEY", "LOG_LEVEL"}

func fullEnv() map[string]string {
	return map[string]string{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
		"API_KEY":   "secret",
		"LOG_LEVEL": "info",
	}
}

func TestScore_Perfect(t *testing.T) {
	res := scorer.Score(fullEnv(), expected, scorer.DefaultOptions())
	if res.Score != 100 {
		t.Fatalf("expected 100, got %.0f", res.Score)
	}
	if res.Grade != "A" {
		t.Fatalf("expected grade A, got %s", res.Grade)
	}
	if len(res.Deductions) != 0 {
		t.Fatalf("expected no deductions, got %d", len(res.Deductions))
	}
}

func TestScore_MissingKey(t *testing.T) {
	env := fullEnv()
	delete(env, "API_KEY")
	opts := scorer.DefaultOptions()
	res := scorer.Score(env, expected, opts)
	if res.Score != 100-opts.MissingPenalty {
		t.Fatalf("unexpected score %.0f", res.Score)
	}
	if len(res.Deductions) != 1 || res.Deductions[0].Key != "API_KEY" {
		t.Fatalf("expected deduction for API_KEY")
	}
}

func TestScore_EmptyValue(t *testing.T) {
	env := fullEnv()
	env["LOG_LEVEL"] = "   "
	opts := scorer.DefaultOptions()
	res := scorer.Score(env, expected, opts)
	if res.Score != 100-opts.EmptyPenalty {
		t.Fatalf("unexpected score %.0f", res.Score)
	}
}

func TestScore_ClampedToZero(t *testing.T) {
	opts := scorer.DefaultOptions()
	opts.MissingPenalty = 40
	res := scorer.Score(map[string]string{}, expected, opts)
	if res.Score != 0 {
		t.Fatalf("expected 0, got %.0f", res.Score)
	}
	if res.Grade != "F" {
		t.Fatalf("expected grade F, got %s", res.Grade)
	}
}

func TestScorePair_Mismatch(t *testing.T) {
	left := fullEnv()
	right := fullEnv()
	right["DB_PORT"] = "3306"
	opts := scorer.DefaultOptions()
	res := scorer.ScorePair(left, right, expected, opts)
	if res.Score != 100-opts.MismatchPenalty {
		t.Fatalf("unexpected score %.0f", res.Score)
	}
	found := false
	for _, d := range res.Deductions {
		if d.Key == "DB_PORT" && d.Reason == "value mismatch" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected mismatch deduction for DB_PORT")
	}
}

func TestScorePair_NoMismatch(t *testing.T) {
	res := scorer.ScorePair(fullEnv(), fullEnv(), expected, scorer.DefaultOptions())
	if res.Score != 100 {
		t.Fatalf("expected 100, got %.0f", res.Score)
	}
}
