package summary_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
	"github.com/yourorg/envdiff/internal/summary"
)

func makeResult(missingLeft, missingRight, mismatched []string) differ.Result {
	r := differ.Result{
		MissingInLeft:  make(map[string]string),
		MissingInRight: make(map[string]string),
		Mismatched:     make(map[string][2]string),
		Clean:          make(map[string]string),
	}
	for _, k := range missingLeft {
		r.MissingInLeft[k] = "val"
	}
	for _, k := range missingRight {
		r.MissingInRight[k] = "val"
	}
	for _, k := range mismatched {
		r.Mismatched[k] = [2]string{"a", "b"}
	}
	r.Clean["SHARED"] = "same"
	return r
}

func TestCompute_Clean(t *testing.T) {
	r := makeResult(nil, nil, nil)
	s := summary.Compute(r)
	if s.HasDiff() {
		t.Errorf("expected no diff, got %s", s)
	}
	if s.Clean != 1 {
		t.Errorf("expected 1 clean key, got %d", s.Clean)
	}
}

func TestCompute_WithDiffs(t *testing.T) {
	r := makeResult([]string{"A"}, []string{"B", "C"}, []string{"D"})
	s := summary.Compute(r)
	if !s.HasDiff() {
		t.Error("expected diff to be detected")
	}
	if s.MissingLeft != 1 {
		t.Errorf("expected MissingLeft=1, got %d", s.MissingLeft)
	}
	if s.MissingRight != 2 {
		t.Errorf("expected MissingRight=2, got %d", s.MissingRight)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
}

func TestStats_String(t *testing.T) {
	r := makeResult([]string{"X"}, nil, []string{"Y"})
	s := summary.Compute(r)
	out := s.String()
	for _, want := range []string{"total=", "clean=", "missing_left=", "missing_right=", "mismatched="} {
		if !strings.Contains(out, want) {
			t.Errorf("String() missing field %q in %q", want, out)
		}
	}
}
