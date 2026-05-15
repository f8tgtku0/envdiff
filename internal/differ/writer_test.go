package differ

import (
	"bytes"
	"strings"
	"testing"
)

func makeResult() *Result {
	r := newResult()
	r.MissingInRight = []string{"DB_HOST", "API_KEY"}
	r.MissingInLeft = []string{"NEW_VAR"}
	r.Mismatched["PORT"] = ValuePair{Left: "8080", Right: "9090"}
	return r
}

func TestWriteText_Clean(t *testing.T) {
	r := newResult()
	var buf bytes.Buffer
	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestWriteText_ShowsDiffs(t *testing.T) {
	r := makeResult()
	var buf bytes.Buffer
	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Missing in right", "API_KEY", "DB_HOST", "Missing in left", "NEW_VAR", "Mismatched", "PORT", "8080", "9090"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteText_NilWriter(t *testing.T) {
	r := newResult()
	if err := WriteText(r, nil); err == nil {
		t.Error("expected error for nil writer")
	}
}

func TestWriteText_NilResult(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteText(nil, &buf); err == nil {
		t.Error("expected error for nil result")
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	r := makeResult()
	var buf bytes.Buffer
	if err := WriteJSON(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"missingInRight", "missingInLeft", "mismatched", "PORT"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got:\n%s", want, out)
		}
	}
}

func TestWriteJSON_NilWriter(t *testing.T) {
	r := newResult()
	if err := WriteJSON(r, nil); err == nil {
		t.Error("expected error for nil writer")
	}
}

func TestWriteText_SortedOutput(t *testing.T) {
	r := newResult()
	r.MissingInRight = []string{"Z_VAR", "A_VAR", "M_VAR"}
	var buf bytes.Buffer
	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	aIdx := strings.Index(out, "A_VAR")
	mIdx := strings.Index(out, "M_VAR")
	zIdx := strings.Index(out, "Z_VAR")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Errorf("expected sorted output, got:\n%s", out)
	}
}
