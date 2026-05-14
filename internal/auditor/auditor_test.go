package auditor

import (
	"testing"
)

func TestAudit_Added(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "1", "B": "2"}

	report, err := Audit(before, after, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	e := report.Entries[0]
	if e.Key != "B" || e.Action != "added" || e.NewValue != "2" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestAudit_Removed(t *testing.T) {
	before := map[string]string{"A": "1", "B": "2"}
	after := map[string]string{"A": "1"}

	report, err := Audit(before, after, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	if report.Entries[0].Action != "removed" {
		t.Errorf("expected removed, got %s", report.Entries[0].Action)
	}
}

func TestAudit_Changed(t *testing.T) {
	before := map[string]string{"KEY": "old"}
	after := map[string]string{"KEY": "new"}

	report, err := Audit(before, after, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 1 || report.Entries[0].Action != "changed" {
		t.Fatalf("expected 1 changed entry")
	}
	e := report.Entries[0]
	if e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("unexpected values: %+v", e)
	}
}

func TestAudit_IncludeUnchanged(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "1"}

	opts := DefaultOptions()
	opts.IncludeUnchanged = true

	report, err := Audit(before, after, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 1 || report.Entries[0].Action != "unchanged" {
		t.Fatalf("expected 1 unchanged entry, got %d", len(report.Entries))
	}
}

func TestAudit_RedactValues(t *testing.T) {
	before := map[string]string{"SECRET": "old-secret"}
	after := map[string]string{"SECRET": "new-secret"}

	opts := DefaultOptions()
	opts.RedactValues = true

	report, err := Audit(before, after, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	e := report.Entries[0]
	if e.OldValue != "***" || e.NewValue != "***" {
		t.Errorf("expected redacted values, got old=%q new=%q", e.OldValue, e.NewValue)
	}
}

func TestAudit_NilBeforeReturnsError(t *testing.T) {
	_, err := Audit(nil, map[string]string{}, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil before")
	}
}

func TestAudit_NilAfterReturnsError(t *testing.T) {
	_, err := Audit(map[string]string{}, nil, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil after")
	}
}
