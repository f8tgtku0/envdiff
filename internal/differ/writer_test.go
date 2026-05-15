package differ

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makeResult(missingRight, missingLeft map[string]string, mismatched map[string]Pair) *Result {
	r := newResult()
	for k, v := range missingRight {
		r.MissingInRight[k] = v
	}
	for k, v := range missingLeft {
		r.MissingInLeft[k] = v
	}
	for k, p := range mismatched {
		r.Mismatched[k] = p
	}
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
	r := makeResult(
		map[string]string{"ONLY_LEFT": "val"},
		map[string]string{"ONLY_RIGHT": "val2"},
		map[string]Pair{"SHARED": {Left: "a", Right: "b"}},
	)
	var buf bytes.Buffer
	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ONLY_LEFT") {
		t.Errorf("expected ONLY_LEFT in output")
	}
	if !strings.Contains(out, "ONLY_RIGHT") {
		t.Errorf("expected ONLY_RIGHT in output")
	}
	if !strings.Contains(out, "SHARED") {
		t.Errorf("expected SHARED in output")
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
	r := makeResult(
		map[string]string{"A": "1"},
		nil,
		nil,
	)
	var buf bytes.Buffer
	if err := WriteJSON(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["missing_in_right"]; !ok {
		t.Error("expected missing_in_right field in JSON")
	}
}

func TestWriteJSON_NilWriter(t *testing.T) {
	r := newResult()
	if err := WriteJSON(r, nil); err == nil {
		t.Error("expected error for nil writer")
	}
}
