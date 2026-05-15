package cloner_test

import (
	"testing"

	"github.com/user/envdiff/internal/cloner"
)

var base = map[string]string{
	"APP_HOST": "localhost",
	"APP_PORT": "8080",
	"DB_URL":   "postgres://localhost/dev",
}

func TestClone_BasicCopy(t *testing.T) {
	out, err := cloner.Clone(base, nil, cloner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(base) {
		t.Fatalf("expected %d keys, got %d", len(base), len(out))
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", out["APP_HOST"])
	}
}

func TestClone_AddPrefix(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.AddPrefix = "STAGING_"
	out, err := cloner.Clone(base, nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["STAGING_APP_HOST"]; !ok {
		t.Error("expected STAGING_APP_HOST to exist")
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("original key APP_HOST should not exist after prefix addition")
	}
}

func TestClone_StripPrefix(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.StripPrefix = "APP_"
	out, err := cloner.Clone(base, nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost after strip, got %q", out["HOST"])
	}
	// DB_URL has no APP_ prefix — should pass through unchanged
	if out["DB_URL"] != "postgres://localhost/dev" {
		t.Errorf("expected DB_URL to be preserved, got %q", out["DB_URL"])
	}
}

func TestClone_StrictPrefix_Error(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.StripPrefix = "APP_"
	opts.StrictPrefix = true
	_, err := cloner.Clone(base, nil, opts)
	if err == nil {
		t.Fatal("expected error for key without required prefix, got nil")
	}
}

func TestClone_NoOverwrite(t *testing.T) {
	dst := map[string]string{"APP_HOST": "original"}
	opts := cloner.DefaultOptions()
	opts.Overwrite = false
	out, err := cloner.Clone(base, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_HOST"] != "original" {
		t.Errorf("expected APP_HOST to remain 'original', got %q", out["APP_HOST"])
	}
}

func TestClone_NilSrcReturnsError(t *testing.T) {
	_, err := cloner.Clone(nil, nil, cloner.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil src, got nil")
	}
}
