package sanitizer

import (
	"sort"
	"testing"
)

func TestSanitize_NoChanges(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	res, err := Sanitize(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("unexpected value: %q", res.Env["HOST"])
	}
}

func TestSanitize_TrimSpace(t *testing.T) {
	env := map[string]string{"KEY": "  hello world  "}
	res, err := Sanitize(env, Options{TrimSpace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "hello world" {
		t.Errorf("got %q, want %q", res.Env["KEY"], "hello world")
	}
	if len(res.Changed) != 1 || res.Changed[0] != "KEY" {
		t.Errorf("expected KEY in changed, got %v", res.Changed)
	}
}

func TestSanitize_StripControl(t *testing.T) {
	env := map[string]string{"VAL": "hello\x00world\x1b"}
	res, err := Sanitize(env, Options{StripControl: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["VAL"] != "helloworld" {
		t.Errorf("got %q, want %q", res.Env["VAL"], "helloworld")
	}
}

func TestSanitize_MaxValueLen(t *testing.T) {
	env := map[string]string{"SECRET": "supersecretvalue"}
	res, err := Sanitize(env, Options{MaxValueLen: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["SECRET"] != "super" {
		t.Errorf("got %q, want %q", res.Env["SECRET"], "super")
	}
}

func TestSanitize_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"K": "  value  "}
	_, err := Sanitize(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["K"] != "  value  " {
		t.Error("original map was mutated")
	}
}

func TestSanitize_NilEnvReturnsError(t *testing.T) {
	_, err := Sanitize(nil, DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestSanitize_MultipleChanges(t *testing.T) {
	env := map[string]string{
		"A": "  trimme  ",
		"B": "clean",
		"C": "\x01ctrl",
	}
	res, err := Sanitize(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(res.Changed)
	if len(res.Changed) != 2 {
		t.Errorf("expected 2 changed keys, got %v", res.Changed)
	}
	if res.Env["B"] != "clean" {
		t.Errorf("clean value should be unchanged, got %q", res.Env["B"])
	}
}
