package counter_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/counter"
)

func makeEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		".env.dev": {
			"HOST": "localhost",
			"PORT": "8080",
			"DEBUG": "",
		},
		".env.prod": {
			"HOST": "example.com",
			"PORT": "443",
			"SECRET": "abc123",
		},
	}
}

func TestCount_Total(t *testing.T) {
	r := counter.Count(makeEnvs())
	if r.Total != 6 {
		t.Errorf("expected Total=6, got %d", r.Total)
	}
}

func TestCount_PerEnv(t *testing.T) {
	r := counter.Count(makeEnvs())
	if r.PerEnv[".env.dev"] != 3 {
		t.Errorf("expected 3 keys for .env.dev, got %d", r.PerEnv[".env.dev"])
	}
	if r.PerEnv[".env.prod"] != 3 {
		t.Errorf("expected 3 keys for .env.prod, got %d", r.PerEnv[".env.prod"])
	}
}

func TestCount_UniqueKeys(t *testing.T) {
	r := counter.Count(makeEnvs())
	// HOST, PORT, DEBUG, SECRET = 4 unique keys
	if r.UniqueKeys != 4 {
		t.Errorf("expected UniqueKeys=4, got %d", r.UniqueKeys)
	}
}

func TestCount_EmptyValues(t *testing.T) {
	r := counter.Count(makeEnvs())
	if r.EmptyValues != 1 {
		t.Errorf("expected EmptyValues=1, got %d", r.EmptyValues)
	}
}

func TestCount_EmptyInput(t *testing.T) {
	r := counter.Count(map[string]map[string]string{})
	if r.Total != 0 || r.UniqueKeys != 0 || r.EmptyValues != 0 {
		t.Errorf("expected zero counts for empty input, got %+v", r)
	}
}

func TestSortedLabels(t *testing.T) {
	r := counter.Count(makeEnvs())
	labels := counter.SortedLabels(r)
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
	if labels[0] != ".env.dev" || labels[1] != ".env.prod" {
		t.Errorf("unexpected order: %v", labels)
	}
}
