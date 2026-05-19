package differ_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func TestSimilarity_Identical(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	r := differ.Similarity(env, env, differ.DefaultSimilarityOptions())
	if r.Score != 1.0 {
		t.Fatalf("expected score 1.0, got %.4f", r.Score)
	}
	if r.Mismatched != 0 || r.MissingLeft != 0 || r.MissingRight != 0 {
		t.Fatalf("expected no diffs, got %+v", r)
	}
}

func TestSimilarity_AllMissing(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"B": "2"}
	r := differ.Similarity(left, right, differ.DefaultSimilarityOptions())
	if r.Score >= 1.0 {
		t.Fatalf("expected score < 1.0, got %.4f", r.Score)
	}
	if r.MissingLeft != 1 || r.MissingRight != 1 {
		t.Fatalf("expected 1 missing each side, got %+v", r)
	}
}

func TestSimilarity_Mismatched(t *testing.T) {
	left := map[string]string{"A": "old", "B": "same"}
	right := map[string]string{"A": "new", "B": "same"}
	r := differ.Similarity(left, right, differ.DefaultSimilarityOptions())
	if r.Mismatched != 1 {
		t.Fatalf("expected 1 mismatch, got %d", r.Mismatched)
	}
	if r.Score >= 1.0 {
		t.Fatalf("expected score < 1.0, got %.4f", r.Score)
	}
}

func TestSimilarity_Empty(t *testing.T) {
	r := differ.Similarity(nil, nil, differ.DefaultSimilarityOptions())
	if r.Score != 1.0 {
		t.Fatalf("expected score 1.0 for empty envs, got %.4f", r.Score)
	}
	if r.TotalKeys != 0 {
		t.Fatalf("expected 0 total keys, got %d", r.TotalKeys)
	}
}

func TestSimilarity_SortedKeys(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "2", "M": "3"}
	right := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := differ.Similarity(left, right, differ.DefaultSimilarityOptions())
	for i := 1; i < len(r.SortedKeys); i++ {
		if r.SortedKeys[i] < r.SortedKeys[i-1] {
			t.Fatalf("keys not sorted: %v", r.SortedKeys)
		}
	}
}

func TestWriteSimilarityText_Clean(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	r := differ.Similarity(env, env, differ.DefaultSimilarityOptions())
	var buf bytes.Buffer
	if err := differ.WriteSimilarityText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "100.0%") {
		t.Fatalf("expected 100.0%% in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Fatalf("expected clean message, got: %s", buf.String())
	}
}

func TestWriteSimilarityJSON_ValidStructure(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "2"}
	r := differ.Similarity(left, right, differ.DefaultSimilarityOptions())
	var buf bytes.Buffer
	if err := differ.WriteSimilarityJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"score"`) {
		t.Fatalf("expected JSON with score field, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), `"grade"`) {
		t.Fatalf("expected JSON with grade field, got: %s", buf.String())
	}
}

func TestWriteSimilarityText_NilWriter(t *testing.T) {
	r := differ.SimilarityReport{}
	if err := differ.WriteSimilarityText(nil, r); err == nil {
		t.Fatal("expected error for nil writer")
	}
}
