package compactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/compactor"
)

func TestCompact_EmptyLayers(t *testing.T) {
	res, err := compactor.Compact(nil, compactor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Compacted) != 0 {
		t.Errorf("expected empty compacted map, got %v", res.Compacted)
	}
}

func TestCompact_SingleLayer(t *testing.T) {
	layer := map[string]string{"A": "1", "B": "2"}
	res, err := compactor.Compact([]map[string]string{layer}, compactor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Compacted["A"] != "1" || res.Compacted["B"] != "2" {
		t.Errorf("unexpected compacted values: %v", res.Compacted)
	}
	if len(res.Overridden) != 0 {
		t.Errorf("expected no overrides for single layer, got %v", res.Overridden)
	}
}

func TestCompact_LaterLayerWins(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	override := map[string]string{"HOST": "prod.db"}
	res, err := compactor.Compact([]map[string]string{base, override}, compactor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Compacted["HOST"] != "prod.db" {
		t.Errorf("expected HOST=prod.db, got %q", res.Compacted["HOST"])
	}
	if res.Compacted["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", res.Compacted["PORT"])
	}
	if res.Overridden["HOST"] != 1 {
		t.Errorf("expected HOST overridden once, got %d", res.Overridden["HOST"])
	}
}

func TestCompact_SkipEmpty(t *testing.T) {
	base := map[string]string{"A": "hello", "B": "", "C": "world"}
	opts := compactor.DefaultOptions()
	opts.SkipEmpty = true
	res, err := compactor.Compact([]map[string]string{base}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Compacted["B"]; ok {
		t.Error("expected B to be dropped")
	}
	if len(res.Dropped) != 1 || res.Dropped[0] != "B" {
		t.Errorf("expected Dropped=[B], got %v", res.Dropped)
	}
	if res.Compacted["A"] != "hello" || res.Compacted["C"] != "world" {
		t.Errorf("unexpected compacted values: %v", res.Compacted)
	}
}

func TestCompact_NilLayerSkipped(t *testing.T) {
	base := map[string]string{"X": "1"}
	res, err := compactor.Compact([]map[string]string{base, nil}, compactor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Compacted["X"] != "1" {
		t.Errorf("expected X=1, got %q", res.Compacted["X"])
	}
}

func TestCompact_MultipleOverrides(t *testing.T) {
	l1 := map[string]string{"K": "v1"}
	l2 := map[string]string{"K": "v2"}
	l3 := map[string]string{"K": "v3"}
	res, _ := compactor.Compact([]map[string]string{l1, l2, l3}, compactor.DefaultOptions())
	if res.Compacted["K"] != "v3" {
		t.Errorf("expected K=v3, got %q", res.Compacted["K"])
	}
	if res.Overridden["K"] != 2 {
		t.Errorf("expected K overridden twice, got %d", res.Overridden["K"])
	}
}
