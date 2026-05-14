package normalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/normalizer"
)

func TestNormalize_TrimSpace(t *testing.T) {
	env := map[string]string{
		"KEY": "  hello  ",
		"OTHER": "\tworld\t",
	}
	opts := normalizer.DefaultOptions()
	opts.TrimSpace = true

	got, err := normalizer.Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "hello" {
		t.Errorf("KEY: got %q, want %q", got["KEY"], "hello")
	}
	if got["OTHER"] != "world" {
		t.Errorf("OTHER: got %q, want %q", got["OTHER"], "world")
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
		"App_Port": "8080",
	}
	opts := normalizer.DefaultOptions()
	opts.UppercaseKeys = true

	got, err := normalizer.Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to exist after uppercasing")
	}
	if _, ok := got["APP_PORT"]; !ok {
		t.Error("expected APP_PORT to exist after uppercasing")
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	env := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
	}
	opts := normalizer.DefaultOptions()
	opts.RemoveEmpty = true

	got, err := normalizer.Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if got["PRESENT"] != "value" {
		t.Errorf("PRESENT: got %q, want %q", got["PRESENT"], "value")
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": "  spaced  "}
	opts := normalizer.DefaultOptions()

	_, err := normalizer.Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "  spaced  " {
		t.Error("Normalize must not mutate the input map")
	}
}

func TestNormalize_NilEnvReturnsError(t *testing.T) {
	_, err := normalizer.Normalize(nil, normalizer.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env, got nil")
	}
}

func TestNormalize_NoOpsPreserveOriginal(t *testing.T) {
	env := map[string]string{"key": " val ", "EMPTY": ""}
	opts := normalizer.Options{TrimSpace: false, UppercaseKeys: false, RemoveEmpty: false}

	got, err := normalizer.Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != " val " {
		t.Errorf("key: got %q, want %q", got["key"], " val ")
	}
	if got["EMPTY"] != "" {
		t.Errorf("EMPTY: got %q, want empty string", got["EMPTY"])
	}
}
