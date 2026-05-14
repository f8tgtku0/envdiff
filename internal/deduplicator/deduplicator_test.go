package deduplicator_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduplicator"
)

func TestDeduplicate_NoConflicts(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"C": "3"}

	out, rep := deduplicator.Deduplicate([]map[string]string{a, b}, deduplicator.DefaultOptions())

	if rep.HasConflicts() {
		t.Fatalf("expected no conflicts, got %v", rep.Conflicts)
	}
	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected merged map: %v", out)
	}
}

func TestDeduplicate_PreferFirst(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	opts := deduplicator.DefaultOptions() // PreferLast == false
	out, rep := deduplicator.Deduplicate([]map[string]string{a, b}, opts)

	if out["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", out["KEY"])
	}
	if !rep.HasConflicts() {
		t.Fatal("expected a conflict to be reported")
	}
	if rep.Conflicts[0].Key != "KEY" {
		t.Errorf("unexpected conflict key: %s", rep.Conflicts[0].Key)
	}
}

func TestDeduplicate_PreferLast(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	opts := deduplicator.Options{PreferLast: true}
	out, rep := deduplicator.Deduplicate([]map[string]string{a, b}, opts)

	if out["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", out["KEY"])
	}
	if !rep.HasConflicts() {
		t.Fatal("expected conflict in report")
	}
}

func TestDeduplicate_SameValueNoConflict(t *testing.T) {
	a := map[string]string{"X": "same"}
	b := map[string]string{"X": "same"}

	_, rep := deduplicator.Deduplicate([]map[string]string{a, b}, deduplicator.DefaultOptions())

	if rep.HasConflicts() {
		t.Errorf("identical values should not be reported as conflicts")
	}
}

func TestDeduplicate_ReportOnly(t *testing.T) {
	a := map[string]string{"K": "alpha"}
	b := map[string]string{"K": "beta"}

	opts := deduplicator.Options{ReportOnly: true}
	out, rep := deduplicator.Deduplicate([]map[string]string{a, b}, opts)

	// Result should still contain the first value.
	if out["K"] != "alpha" {
		t.Errorf("report-only should preserve first value, got %q", out["K"])
	}
	if !rep.HasConflicts() {
		t.Fatal("expected conflict in report")
	}
}

func TestDeduplicate_EmptySources(t *testing.T) {
	out, rep := deduplicator.Deduplicate([]map[string]string{}, deduplicator.DefaultOptions())

	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
	if rep.HasConflicts() {
		t.Error("expected no conflicts for empty input")
	}
}

func TestDeduplicate_ConflictsSorted(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "x", "M": "p"}
	b := map[string]string{"Z": "2", "A": "y", "M": "q"}

	_, rep := deduplicator.Deduplicate([]map[string]string{a, b}, deduplicator.DefaultOptions())

	if len(rep.Conflicts) != 3 {
		t.Fatalf("expected 3 conflicts, got %d", len(rep.Conflicts))
	}
	keys := []string{rep.Conflicts[0].Key, rep.Conflicts[1].Key, rep.Conflicts[2].Key}
	if keys[0] > keys[1] || keys[1] > keys[2] {
		t.Errorf("conflicts not sorted: %v", keys)
	}
}
