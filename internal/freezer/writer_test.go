package freezer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/freezer"
)

func TestWriteText_ShowsViolations(t *testing.T) {
	r := freezer.Result{
		Violations: []freezer.Violation{
			{Key: "DB_PASS", Reason: "value changed (frozen=\"old\", live=\"new\")"},
			{Key: "API_KEY", Reason: "key removed from live env"},
		},
	}
	var buf bytes.Buffer
	if err := freezer.WriteText(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output: %s", out)
	}
	if !strings.Contains(out, "2 violation") {
		t.Errorf("expected violation count in output: %s", out)
	}
}

func TestWriteText_NilWriter(t *testing.T) {
	r := freezer.Result{}
	if err := freezer.WriteText(r, nil); err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestWriteJSON_Clean(t *testing.T) {
	r := freezer.Result{}
	var buf bytes.Buffer
	if err := freezer.WriteJSON(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"clean": true`) {
		t.Errorf("expected clean:true in output: %s", out)
	}
	if !strings.Contains(out, `"violations"`) {
		t.Errorf("expected violations key in output: %s", out)
	}
}

func TestWriteJSON_NilWriter(t *testing.T) {
	r := freezer.Result{}
	if err := freezer.WriteJSON(r, nil); err == nil {
		t.Fatal("expected error for nil writer")
	}
}
