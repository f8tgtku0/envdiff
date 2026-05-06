package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Merged["FOO"] != "1" || res.Merged["BAR"] != "2" || res.Merged["BAZ"] != "3" {
		t.Errorf("unexpected merged map: %v", res.Merged)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "from-a" {
		t.Errorf("expected 'from-a', got %q", res.Merged["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "from-b" {
		t.Errorf("expected 'from-b', got %q", res.Merged["KEY"])
	}
}

func TestMerge_StrategyError(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	_, err := merger.Merge([]map[string]string{a, b}, merger.StrategyError)
	if err == nil {
		t.Fatal("expected an error for conflicting key, got nil")
	}
}

func TestMerge_SameValueNoConflict(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyError)
	if err != nil {
		t.Fatalf("unexpected error for identical values: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts for identical values, got %d", len(res.Conflicts))
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res, err := merger.Merge([]map[string]string{}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 0 {
		t.Errorf("expected empty map, got %v", res.Merged)
	}
}
