package interpolator

import (
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "5432" {
		t.Errorf("values should be unchanged, got %v", out)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE": "/app",
		"DATA": "${BASE}/data",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATA"] != "/app/data" {
		t.Errorf("expected /app/data, got %q", out["DATA"])
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	env := map[string]string{
		"SCHEME": "https",
		"HOST":   "example.com",
		"URL":    "$SCHEME://$HOST",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "https://example.com" {
		t.Errorf("expected https://example.com, got %q", out["URL"])
	}
}

func TestInterpolate_ChainedReferences(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["C"] != "hello_world!" {
		t.Errorf("expected hello_world!, got %q", out["C"])
	}
}

func TestInterpolate_UndefinedLenient(t *testing.T) {
	env := map[string]string{
		"VAL": "${MISSING}_suffix",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error in lenient mode: %v", err)
	}
	// unresolved reference should remain as-is
	if out["VAL"] != "${MISSING}_suffix" {
		t.Errorf("expected reference to remain, got %q", out["VAL"])
	}
}

func TestInterpolate_UndefinedStrict(t *testing.T) {
	env := map[string]string{
		"VAL": "${MISSING}_suffix",
	}
	opts := Options{Strict: true}
	_, err := Interpolate(env, opts)
	if err == nil {
		t.Fatal("expected error in strict mode for undefined variable")
	}
}

func TestInterpolate_SelfReferenceStrict(t *testing.T) {
	env := map[string]string{
		"A": "${A}_loop",
	}
	opts := Options{Strict: true}
	_, err := Interpolate(env, opts)
	if err == nil {
		t.Fatal("expected error for self-referential key in strict mode")
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"BASE": "/opt",
		"DIR":  "${BASE}/bin",
	}
	original := env["DIR"]
	_, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DIR"] != original {
		t.Errorf("input map was mutated")
	}
}
