package stripper_test

import (
	"testing"

	"github.com/user/envdiff/internal/stripper"
)

var baseEnv = map[string]string{
	"APP_NAME":    "myapp",
	"APP_VERSION": "1.0",
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"INTERNAL_ID": "42",
	"DEBUG":       "true",
}

func TestStrip_NoRules(t *testing.T) {
	res, err := stripper.Strip(baseEnv, stripper.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Stripped) != 0 {
		t.Errorf("expected no stripped keys, got %v", res.Stripped)
	}
	if len(res.Env) != len(baseEnv) {
		t.Errorf("expected %d keys, got %d", len(baseEnv), len(res.Env))
	}
}

func TestStrip_ByPrefix(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	res, err := stripper.Strip(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("DB_HOST should have been stripped")
	}
	if _, ok := res.Env["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been stripped")
	}
	if len(res.Stripped) != 2 {
		t.Errorf("expected 2 stripped, got %d", len(res.Stripped))
	}
}

func TestStrip_BySuffix(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Suffixes = []string{"_ID"}
	res, err := stripper.Strip(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["INTERNAL_ID"]; ok {
		t.Error("INTERNAL_ID should have been stripped")
	}
}

func TestStrip_ByPattern(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Patterns = []string{"^APP_"}
	res, err := stripper.Strip(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["APP_NAME"]; ok {
		t.Error("APP_NAME should have been stripped")
	}
	if _, ok := res.Env["APP_VERSION"]; ok {
		t.Error("APP_VERSION should have been stripped")
	}
}

func TestStrip_ByExplicitKeys(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Keys = []string{"DEBUG", "APP_NAME"}
	res, err := stripper.Strip(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DEBUG"]; ok {
		t.Error("DEBUG should have been stripped")
	}
	if _, ok := res.Env["APP_NAME"]; ok {
		t.Error("APP_NAME should have been stripped")
	}
}

func TestStrip_DryRunDoesNotMutateResult(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	opts.DryRun = true
	res, err := stripper.Strip(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("dry run: DB_HOST should still be present in Env")
	}
	if len(res.Stripped) == 0 {
		t.Error("dry run: Stripped list should still be populated")
	}
}

func TestStrip_InvalidPattern(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Patterns = []string{"[invalid"}
	_, err := stripper.Strip(baseEnv, opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestStrip_NilEnvReturnsError(t *testing.T) {
	_, err := stripper.Strip(nil, stripper.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}
