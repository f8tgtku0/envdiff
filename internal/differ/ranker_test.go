package differ

import (
	"testing"
)

func makeRankResult(missingRight, missingLeft []string, mismatched map[string][2]string) *Result {
	r := newResult()
	r.MissingInRight = append(r.MissingInRight, missingRight...)
	r.MissingInLeft = append(r.MissingInLeft, missingLeft...)
	for k, v := range mismatched {
		r.Mismatched[k] = v
	}
	return r
}

func TestRank_EmptyResults(t *testing.T) {
	entries := Rank(DefaultRankOptions())
	if len(entries) != 0 {
		t.Fatalf("expected empty rank, got %d entries", len(entries))
	}
}

func TestRank_NilResultSkipped(t *testing.T) {
	entries := Rank(DefaultRankOptions(), nil, nil)
	if len(entries) != 0 {
		t.Fatalf("expected empty rank, got %d entries", len(entries))
	}
}

func TestRank_MismatchScoresHigher(t *testing.T) {
	r := makeRankResult(
		[]string{"MISSING_KEY"},
		nil,
		map[string][2]string{"MISMATCH_KEY": {"a", "b"}},
	)
	entriesopts := DefaultRankOptions()
	entries := Rank(entriesopts, r)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "MISMATCH_KEY" {
		t.Errorf("expected MISMATCH_KEY first, got %s", entries[0].Key)
	}
	if entries[0].Score != 2 {
		t.Errorf("expected score 2, got %d", entries[0].Score)
	}
	if entries[1].Score != 1 {
		t.Errorf("expected score 1 for missing key, got %d", entries[1].Score)
	}
}

func TestRank_AccumulatesAcrossResults(t *testing.T) {
	r1 := makeRankResult([]string{"KEY_A"}, nil, nil)
	r2 := makeRankResult([]string{"KEY_A"}, nil, nil)
	entriesopts := DefaultRankOptions()
	entries := Rank(entriesopts, r1, r2)

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Score != 2 {
		t.Errorf("expected accumulated score 2, got %d", entries[0].Score)
	}
}

func TestRank_SortedByScoreThenKey(t *testing.T) {
	r := makeRankResult(
		[]string{"ZEBRA", "ALPHA"},
		nil,
		nil,
	)
	entries := Rank(DefaultRankOptions(), r)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// Same score — should be alphabetical
	if entries[0].Key != "ALPHA" {
		t.Errorf("expected ALPHA first (alpha sort), got %s", entries[0].Key)
	}
}

func TestRank_CustomWeights(t *testing.T) {
	opts := RankOptions{MissingWeight: 5, MismatchWeight: 1}
	r := makeRankResult(
		[]string{"MISSING_KEY"},
		nil,
		map[string][2]string{"MISMATCH_KEY": {"a", "b"}},
	)
	entries := Rank(opts, r)

	if entries[0].Key != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY first with custom weights, got %s", entries[0].Key)
	}
	if entries[0].Score != 5 {
		t.Errorf("expected score 5, got %d", entries[0].Score)
	}
}
