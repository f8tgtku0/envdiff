package aliaser

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DATABASE_HOST": "localhost",
		"DATABASE_PORT": "5432",
		"APP_SECRET":    "s3cr3t",
	}
}

func TestAlias_Basic(t *testing.T) {
	env := baseEnv()
	aliases := AliasMap{"DB_HOST": "DATABASE_HOST", "DB_PORT": "DATABASE_PORT"}
	out, res, err := Alias(env, aliases, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Error("source key DATABASE_HOST should be preserved")
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(res.Applied))
	}
}

func TestAlias_SkipsExistingWithoutOverwrite(t *testing.T) {
	env := baseEnv()
	env["DB_HOST"] = "existing"
	aliases := AliasMap{"DB_HOST": "DATABASE_HOST"}
	out, res, err := Alias(env, aliases, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "existing" {
		t.Errorf("expected existing value preserved, got %q", out["DB_HOST"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in skipped, got %v", res.Skipped)
	}
}

func TestAlias_OverwriteExisting(t *testing.T) {
	env := baseEnv()
	env["DB_HOST"] = "existing"
	aliases := AliasMap{"DB_HOST": "DATABASE_HOST"}
	opts := Options{Overwrite: true}
	out, res, err := Alias(env, aliases, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected overwritten value localhost, got %q", out["DB_HOST"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestAlias_MissingSourceKey(t *testing.T) {
	env := baseEnv()
	aliases := AliasMap{"NEW_KEY": "NONEXISTENT"}
	out, res, err := Alias(env, aliases, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["NEW_KEY"]; ok {
		t.Error("NEW_KEY should not be present when source is missing")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "NONEXISTENT" {
		t.Errorf("expected NONEXISTENT in missing, got %v", res.Missing)
	}
}

func TestAlias_NilEnvReturnsError(t *testing.T) {
	_, _, err := Alias(nil, AliasMap{"A": "B"}, DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestAlias_EmptyAliasMap(t *testing.T) {
	env := baseEnv()
	out, res, err := Alias(env, AliasMap{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(env) {
		t.Errorf("expected output same size as input, got %d", len(out))
	}
	if len(res.Applied) != 0 {
		t.Errorf("expected no applied entries, got %d", len(res.Applied))
	}
}
