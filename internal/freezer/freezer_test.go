package freezer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/freezer"
)

func TestFreeze_NoViolations(t *testing.T) {
	frozen := map[string]string{"A": "1", "B": "2"}
	live := map[string]string{"A": "1", "B": "2"}
	r, err := freezer.Freeze(frozen, live, freezer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Clean() {
		t.Fatalf("expected clean result, got %+v", r.Violations)
	}
}

func TestFreeze_RemovedKey(t *testing.T) {
	frozen := map[string]string{"A": "1", "B": "2"}
	live := map[string]string{"A": "1"}
	r, err := freezer.Freeze(frozen, live, freezer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Clean() {
		t.Fatal("expected violations for removed key")
	}
	if r.Violations[0].Key != "B" {
		t.Errorf("expected key B, got %s", r.Violations[0].Key)
	}
}

func TestFreeze_ChangedValue(t *testing.T) {
	frozen := map[string]string{"A": "original"}
	live := map[string]string{"A": "changed"}
	r, err := freezer.Freeze(frozen, live, freezer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Clean() {
		t.Fatal("expected violation for changed value")
	}
	if !strings.Contains(r.Violations[0].Reason, "value changed") {
		t.Errorf("unexpected reason: %s", r.Violations[0].Reason)
	}
}

func TestFreeze_UnexpectedKeyBlocked(t *testing.T) {
	frozen := map[string]string{"A": "1"}
	live := map[string]string{"A": "1", "EXTRA": "x"}
	r, err := freezer.Freeze(frozen, live, freezer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Clean() {
		t.Fatal("expected violation for unexpected key")
	}
}

func TestFreeze_UnexpectedKeyAllowed(t *testing.T) {
	frozen := map[string]string{"A": "1"}
	live := map[string]string{"A": "1", "EXTRA": "x"}
	opts := freezer.DefaultOptions()
	opts.AllowExpand = true
	r, err := freezer.Freeze(frozen, live, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Clean() {
		t.Fatalf("expected clean result, got %+v", r.Violations)
	}
}

func TestFreeze_IgnoreKeys(t *testing.T) {
	frozen := map[string]string{"A": "1", "SKIP": "old"}
	live := map[string]string{"A": "1", "SKIP": "new"}
	opts := freezer.DefaultOptions()
	opts.IgnoreKeys = []string{"SKIP"}
	r, err := freezer.Freeze(frozen, live, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Clean() {
		t.Fatalf("expected clean result, got %+v", r.Violations)
	}
}

func TestFreeze_NilFrozenReturnsError(t *testing.T) {
	_, err := freezer.Freeze(nil, map[string]string{}, freezer.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil frozen")
	}
}

func TestWriteText_Clean(t *testing.T) {
	r := freezer.Result{}
	var buf bytes.Buffer
	if err := freezer.WriteText(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWriteJSON_WithViolations(t *testing.T) {
	r := freezer.Result{Violations: []freezer.Violation{{Key: "X", Reason: "key removed from live env"}}}
	var buf bytes.Buffer
	if err := freezer.WriteJSON(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"clean": false`) {
		t.Errorf("expected clean:false in JSON output: %s", buf.String())
	}
}
