package patcher_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/patcher"
)

func TestWriteReport_NoChanges(t *testing.T) {
	var sb strings.Builder
	if err := patcher.WriteReport(&sb, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No changes") {
		t.Errorf("expected 'No changes', got %q", sb.String())
	}
}

func TestWriteReport_WithChanges(t *testing.T) {
	changes := []patcher.Change{
		{Key: "FOO", OldValue: "", NewValue: "bar", Action: "add"},
		{Key: "BAZ", OldValue: "old", NewValue: "new", Action: "update"},
		{Key: "QUX", OldValue: "keep", NewValue: "ignored", Action: "skip"},
	}
	var sb strings.Builder
	if err := patcher.WriteReport(&sb, changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, want := range []string{"add", "update", "skip", "FOO", "BAZ", "QUX"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in report output", want)
		}
	}
}

func TestSummarise(t *testing.T) {
	changes := []patcher.Change{
		{Action: "add"},
		{Action: "add"},
		{Action: "update"},
		{Action: "skip"},
	}
	s := patcher.Summarise(changes)
	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Updated != 1 {
		t.Errorf("expected Updated=1, got %d", s.Updated)
	}
	if s.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", s.Skipped)
	}
}

func TestStats_String(t *testing.T) {
	s := patcher.Stats{Added: 3, Updated: 1, Skipped: 2}
	got := s.String()
	if !strings.Contains(got, "added=3") || !strings.Contains(got, "updated=1") || !strings.Contains(got, "skipped=2") {
		t.Errorf("unexpected Stats.String output: %q", got)
	}
}
