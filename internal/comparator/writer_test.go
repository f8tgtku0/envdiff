package comparator_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/comparator"
)

func TestWriteText_Clean(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"X": "1"},
		"b": {"X": "1"},
	}
	results, _ := comparator.Compare(envs, comparator.Options{RequireAll: true})
	var buf bytes.Buffer
	if err := comparator.WriteText(results, []string{"a", "b"}, &buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "All keys are consistent") {
		t.Error("expected clean message for consistent envs")
	}
}

func TestWriteText_ShowsMissing(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"FOO": "bar"},
		"b": {},
	}
	opts := comparator.DefaultOptions()
	opts.OnlyInconsistent = true
	results, _ := comparator.Compare(envs, opts)
	var buf bytes.Buffer
	_ = comparator.WriteText(results, []string{"a", "b"}, &buf)
	if !strings.Contains(buf.String(), "<missing>") {
		t.Error("expected <missing> placeholder in output")
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	envs := map[string]map[string]string{
		"x": {"K": "v1"},
		"y": {"K": "v2"},
	}
	results, _ := comparator.Compare(envs, comparator.DefaultOptions())
	var buf bytes.Buffer
	if err := comparator.WriteJSON(results, &buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if _, ok := rows[0]["key"]; !ok {
		t.Error("expected 'key' field in JSON row")
	}
	if _, ok := rows[0]["consistent"]; !ok {
		t.Error("expected 'consistent' field in JSON row")
	}
}

func TestWriteJSON_InconsistentFlag(t *testing.T) {
	envs := map[string]map[string]string{
		"p": {"SECRET": "abc"},
		"q": {"SECRET": "xyz"},
	}
	results, _ := comparator.Compare(envs, comparator.DefaultOptions())
	var buf bytes.Buffer
	_ = comparator.WriteJSON(results, &buf)

	var rows []map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &rows)
	if rows[0]["consistent"].(bool) {
		t.Error("expected consistent=false for mismatched SECRET")
	}
}
