package renamer_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/renamer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV":  "production",
	}
}

func TestRename_Basic(t *testing.T) {
	out, res, err := renamer.Rename(baseEnv(), renamer.RenameMap{
		"DB_HOST": "DATABASE_HOST",
		"DB_PORT": "DATABASE_PORT",
	}, renamer.DefaultOptions())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 2 {
		t.Errorf("expected 2 renamed, got %d", len(res.Renamed))
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("old key DB_HOST should have been removed")
	}
}

func TestRename_SkipsMissingKey(t *testing.T) {
	_, res, err := renamer.Rename(baseEnv(), renamer.RenameMap{
		"MISSING_KEY": "NEW_KEY",
	}, renamer.DefaultOptions())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in skipped, got %v", res.Skipped)
	}
}

func TestRename_ConflictNoOverwrite(t *testing.T) {
	env := map[string]string{
		"OLD_KEY": "value1",
		"NEW_KEY": "value2",
	}
	_, _, err := renamer.Rename(env, renamer.RenameMap{
		"OLD_KEY": "NEW_KEY",
	}, renamer.DefaultOptions())

	if err == nil {
		t.Error("expected conflict error, got nil")
	}
}

func TestRename_ConflictWithOverwrite(t *testing.T) {
	env := map[string]string{
		"OLD_KEY": "value1",
		"NEW_KEY": "value2",
	}
	out, res, err := renamer.Rename(env, renamer.RenameMap{
		"OLD_KEY": "NEW_KEY",
	}, renamer.Options{OverwriteConflicts: true})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1 after overwrite, got %q", out["NEW_KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict recorded, got %d", len(res.Conflicts))
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	original := map[string]string{}
	for k, v := range env {
		original[k] = v
	}

	renamer.Rename(env, renamer.RenameMap{"DB_HOST": "DATABASE_HOST"}, renamer.DefaultOptions()) //nolint

	for k, v := range original {
		if env[k] != v {
			t.Errorf("input env was mutated: key %q changed", k)
		}
	}
}
