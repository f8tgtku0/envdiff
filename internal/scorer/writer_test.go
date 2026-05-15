package scorer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/scorer"
)

func TestWriteText_Clean(t *testing.T) {
	var buf bytes.Buffer
	r := scorer.Result{Score: 100, Grade: "A"}
	if err := scorer.WriteText(&buf, r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "100") {
		t.Fatalf("expected score in output, got: %s", buf.String())
	}
}

func TestWriteText_WithDeductions(t *testing.T) {
	var buf bytes.Buffer
	r := scorer.Result{
		Score: 80,
		Grade: "B",
		Deductions: []scorer.Deduction{
			{Key: "API_KEY", Reason: "missing", Points: 10},
			{Key: "DB_PORT", Reason: "empty value", Points: 3},
		},
	}
	if err := scorer.WriteText(&buf, r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Fatalf("expected API_KEY in output")
	}
	if !strings.Contains(out, "missing") {
		t.Fatalf("expected reason 'missing' in output")
	}
}

func TestWriteText_NilWriter(t *testing.T) {
	err := scorer.WriteText(nil, scorer.Result{})
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	r := scorer.Result{
		Score: 75,
		Grade: "B",
		Deductions: []scorer.Deduction{
			{Key: "LOG_LEVEL", Reason: "empty value", Points: 3},
		},
	}
	if err := scorer.WriteJSON(&buf, r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"Score\"") {
		t.Fatalf("expected JSON key 'Score', got: %s", out)
	}
	if !strings.Contains(out, "\"Grade\"") {
		t.Fatalf("expected JSON key 'Grade'")
	}
}

func TestWriteJSON_NilWriter(t *testing.T) {
	err := scorer.WriteJSON(nil, scorer.Result{})
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}
