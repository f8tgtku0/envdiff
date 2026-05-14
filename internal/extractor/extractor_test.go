package extractor_test

import (
	"testing"

	"github.com/user/envdiff/internal/extractor"
)

var baseEnv = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"LOG_LEVEL":   "info",
	"SECRET_KEY":  "abc123",
}

func TestExtract_NoFilter(t *testing.T) {
	res, err := extractor.Extract(baseEnv, extractor.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != len(baseEnv) {
		t.Errorf("expected %d keys, got %d", len(baseEnv), len(res.Env))
	}
}

func TestExtract_PrefixFilter(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Prefixes = []string{"APP_"}

	res, err := extractor.Extract(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if _, ok := res.Env["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestExtract_SuffixFilter(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Suffixes = []string{"_HOST"}

	res, err := extractor.Extract(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestExtract_ExplicitKeys(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Keys = []string{"LOG_LEVEL", "SECRET_KEY"}

	res, err := extractor.Extract(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if res.Env["LOG_LEVEL"] != "info" {
		t.Errorf("unexpected value for LOG_LEVEL: %s", res.Env["LOG_LEVEL"])
	}
}

func TestExtract_StripPrefix(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Prefixes = []string{"APP_"}
	opts.StripPrefix = true

	res, err := extractor.Extract(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_ prefix")
	}
	if _, ok := res.Env["PORT"]; !ok {
		t.Error("expected PORT after stripping APP_ prefix")
	}
}

func TestExtract_NilEnvReturnsError(t *testing.T) {
	_, err := extractor.Extract(nil, extractor.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestExtract_ResultIsSorted(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Prefixes = []string{"DB_"}

	res, err := extractor.Extract(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(res.Extracted); i++ {
		if res.Extracted[i-1] > res.Extracted[i] {
			t.Errorf("extracted keys not sorted: %v", res.Extracted)
		}
	}
}
