package highlighter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/highlighter"
)

func TestHighlight_AllOk(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	lines := highlighter.Highlight(left, right, highlighter.DefaultOptions())
	for _, l := range lines {
		if l.Status != "ok" {
			t.Errorf("expected ok for key %s, got %s", l.Key, l.Status)
		}
	}
}

func TestHighlight_MissingRight(t *testing.T) {
	left := map[string]string{"ONLY_LEFT": "val"}
	right := map[string]string{}
	lines := highlighter.Highlight(left, right, highlighter.DefaultOptions())
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0].Status != "missing_right" {
		t.Errorf("expected missing_right, got %s", lines[0].Status)
	}
}

func TestHighlight_MissingLeft(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{"ONLY_RIGHT": "val"}
	lines := highlighter.Highlight(left, right, highlighter.DefaultOptions())
	if len(lines) != 1 || lines[0].Status != "missing_left" {
		t.Errorf("expected missing_left status")
	}
}

func TestHighlight_Mismatch(t *testing.T) {
	left := map[string]string{"KEY": "foo"}
	right := map[string]string{"KEY": "bar"}
	lines := highlighter.Highlight(left, right, highlighter.DefaultOptions())
	if len(lines) != 1 || lines[0].Status != "mismatch" {
		t.Errorf("expected mismatch status")
	}
	if lines[0].Left != "foo" || lines[0].Right != "bar" {
		t.Errorf("unexpected left/right values")
	}
}

func TestHighlight_SortedOutput(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "2", "M": "3"}
	right := map[string]string{"Z": "1", "A": "2", "M": "3"}
	lines := highlighter.Highlight(left, right, highlighter.DefaultOptions())
	keys := make([]string, len(lines))
	for i, l := range lines {
		keys[i] = l.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}

func TestWriteText_NilWriter(t *testing.T) {
	lines := []highlighter.Line{{Key: "K", Status: "ok", Left: "v"}}
	err := highlighter.WriteText(nil, lines, highlighter.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil writer")
	}
}

func TestWriteText_ContainsPrefixes(t *testing.T) {
	opts := highlighter.Options{UseColor: false, MissingPrefix: "- ", MismatchPrefix: "~ "}
	lines := []highlighter.Line{
		{Key: "GONE", Left: "x", Right: "", Status: "missing_right"},
		{Key: "DIFF", Left: "a", Right: "b", Status: "mismatch"},
		{Key: "SAME", Left: "c", Right: "c", Status: "ok"},
	}
	var buf bytes.Buffer
	if err := highlighter.WriteText(&buf, lines, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- GONE") {
		t.Errorf("expected missing prefix in output")
	}
	if !strings.Contains(out, "~ DIFF") {
		t.Errorf("expected mismatch prefix in output")
	}
	if !strings.Contains(out, "  SAME") {
		t.Errorf("expected ok line in output")
	}
}

func TestSummary_Mixed(t *testing.T) {
	lines := []highlighter.Line{
		{Status: "ok"}, {Status: "ok"},
		{Status: "mismatch"},
		{Status: "missing_right"},
	}
	s := highlighter.Summary(lines)
	if !strings.Contains(s, "2 ok") {
		t.Errorf("expected '2 ok' in summary, got: %s", s)
	}
	if !strings.Contains(s, "1 missing") {
		t.Errorf("expected '1 missing' in summary, got: %s", s)
	}
	if !strings.Contains(s, "1 mismatched") {
		t.Errorf("expected '1 mismatched' in summary, got: %s", s)
	}
}

func TestSummary_AllOk(t *testing.T) {
	lines := []highlighter.Line{{Status: "ok"}, {Status: "ok"}}
	s := highlighter.Summary(lines)
	if strings.Contains(s, "missing") || strings.Contains(s, "mismatch") {
		t.Errorf("unexpected diff info in clean summary: %s", s)
	}
}
