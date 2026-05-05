package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reporter"
)

func TestWriteText_Clean(t *testing.T) {
	diff := differ.DiffResult{}
	r := reporter.NewReport(".env.dev", ".env.prod", diff)

	var buf bytes.Buffer
	if err := r.Write(&buf, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestWriteText_MissingAndMismatched(t *testing.T) {
	diff := differ.DiffResult{
		MissingInRight: map[string]string{"DB_HOST": "localhost"},
		MissingInLeft:  map[string]string{"NEW_KEY": "value"},
		Mismatched: map[string]differ.Mismatch{
			"LOG_LEVEL": {LeftVal: "debug", RightVal: "info"},
		},
	}
	r := reporter.NewReport(".env.dev", ".env.prod", diff)

	var buf bytes.Buffer
	if err := r.Write(&buf, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output")
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected LOG_LEVEL in output")
	}
	if !strings.Contains(out, "debug") || !strings.Contains(out, "info") {
		t.Errorf("expected mismatch values in output")
	}
}

func TestWriteJSON_Clean(t *testing.T) {
	diff := differ.DiffResult{}
	r := reporter.NewReport(".env.dev", ".env.prod", diff)

	var buf bytes.Buffer
	if err := r.Write(&buf, reporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["clean"] != true {
		t.Errorf("expected clean=true in JSON output")
	}
}

func TestWriteJSON_WithDiffs(t *testing.T) {
	diff := differ.DiffResult{
		MissingInRight: map[string]string{"SECRET_KEY": "abc123"},
		Mismatched: map[string]differ.Mismatch{
			"PORT": {LeftVal: "3000", RightVal: "8080"},
		},
	}
	r := reporter.NewReport(".env.staging", ".env.prod", diff)

	var buf bytes.Buffer
	if err := r.Write(&buf, reporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["clean"] != false {
		t.Errorf("expected clean=false")
	}
	missingRight, ok := out["missing_in_right"].([]interface{})
	if !ok || len(missingRight) != 1 {
		t.Errorf("expected 1 missing_in_right entry")
	}
}
