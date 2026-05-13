package comparator_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/comparator"
)

func makeEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"PORT":  "3000",
			"DEBUG": "true",
			"HOST":  "localhost",
		},
		"staging": {
			"PORT":  "3000",
			"DEBUG": "false",
			"HOST":  "localhost",
		},
		"prod": {
			"PORT":  "8080",
			"DEBUG": "false",
		},
	}
}

func TestCompare_AllKeys(t *testing.T) {
	results, err := comparator.Compare(makeEnvs(), comparator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestCompare_OnlyInconsistent(t *testing.T) {
	opts := comparator.DefaultOptions()
	opts.OnlyInconsistent = true
	results, err := comparator.Compare(makeEnvs(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Consistent {
			t.Errorf("expected only inconsistent results, got consistent key %q", r.Key)
		}
	}
}

func TestCompare_MissingKey(t *testing.T) {
	results, err := comparator.Compare(makeEnvs(), comparator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "HOST" {
			if r.Consistent {
				t.Error("HOST should be inconsistent because prod is missing it")
			}
			if len(r.Missing) != 1 || r.Missing[0] != "prod" {
				t.Errorf("expected prod in Missing, got %v", r.Missing)
			}
			return
		}
	}
	t.Error("HOST key not found in results")
}

func TestCompare_EmptyEnvs(t *testing.T) {
	_, err := comparator.Compare(map[string]map[string]string{}, comparator.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty envs")
	}
}

func TestCompare_RequireAllFalse(t *testing.T) {
	opts := comparator.DefaultOptions()
	opts.RequireAll = false
	results, err := comparator.Compare(makeEnvs(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "HOST" && !r.Consistent {
			t.Error("HOST should be consistent when RequireAll is false and present values match")
		}
	}
}

func TestWriteText_Output(t *testing.T) {
	results, _ := comparator.Compare(makeEnvs(), comparator.DefaultOptions())
	var buf bytes.Buffer
	err := comparator.WriteText(results, []string{"dev", "staging", "prod"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in text output")
	}
}

func TestWriteJSON_Output(t *testing.T) {
	results, _ := comparator.Compare(makeEnvs(), comparator.DefaultOptions())
	var buf bytes.Buffer
	err := comparator.WriteJSON(results, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"key\"") {
		t.Error("expected JSON key field in output")
	}
}
