package masker_test

import (
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/masker"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":       "myapp",
		"APP_PORT":       "8080",
		"SECRET_KEY":     "supersecret",
		"TOKEN_API":      "tok-abc123",
		"PASSWORD_DB":    "hunter2",
		"PLAIN_VALUE":    "visible",
	}
}

func TestMask_DefaultPrefixes(t *testing.T) {
	env := baseEnv()
	opts := masker.DefaultOptions()

	out, err := masker.Mask(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, k := range []string{"SECRET_KEY", "TOKEN_API", "PASSWORD_DB"} {
		if out[k] != "***" {
			t.Errorf("expected %q to be masked, got %q", k, out[k])
		}
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unmasked")
	}
}

func TestMask_ExplicitKeys(t *testing.T) {
	env := baseEnv()
	opts := masker.Options{
		Keys: []string{"APP_PORT", "PLAIN_VALUE"},
		Mask: "REDACTED",
	}

	out, err := masker.Mask(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out["APP_PORT"] != "REDACTED" {
		t.Errorf("expected APP_PORT to be redacted")
	}
	if out["PLAIN_VALUE"] != "REDACTED" {
		t.Errorf("expected PLAIN_VALUE to be redacted")
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME untouched")
	}
}

func TestMask_ShowLength(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "supersecret"}
	opts := masker.Options{
		Prefixes:   []string{"SECRET_"},
		Mask:       "***",
		ShowLength: true,
	}

	out, err := masker.Mask(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "***[11]"
	if out["SECRET_KEY"] != want {
		t.Errorf("expected %q, got %q", want, out["SECRET_KEY"])
	}
}

func TestMask_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	opts := masker.DefaultOptions()

	_, err := masker.Mask(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env["SECRET_KEY"] != "supersecret" {
		t.Error("original map was mutated")
	}
}

func TestMask_NilEnvReturnsError(t *testing.T) {
	_, err := masker.Mask(nil, masker.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMask_DefaultMaskFallback(t *testing.T) {
	env := map[string]string{"SECRET_X": "value"}
	opts := masker.Options{
		Prefixes: []string{"SECRET_"},
		Mask:     "", // should fall back to "***"
	}

	out, err := masker.Mask(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET_X"] != "***" {
		t.Errorf("expected default mask, got %q", out["SECRET_X"])
	}
}
