package flattener_test

import (
	"testing"

	"github.com/user/envdiff/internal/flattener"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP__DB__HOST":  "localhost",
		"APP__DB__PORT":  "5432",
		"APP__LOG_LEVEL": "info",
		"PLAIN":          "value",
	}
}

func TestFlatten_NoStrip(t *testing.T) {
	env := baseEnv()
	opts := flattener.DefaultOptions()
	out, err := flattener.Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP__DB__HOST"] != "localhost" {
		t.Errorf("expected APP__DB__HOST=localhost, got %q", out["APP__DB__HOST"])
	}
	if out["PLAIN"] != "value" {
		t.Errorf("expected PLAIN=value, got %q", out["PLAIN"])
	}
}

func TestFlatten_StripOneSegment(t *testing.T) {
	env := baseEnv()
	opts := flattener.DefaultOptions()
	opts.StripDepth = 1
	out, err := flattener.Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB__HOST"]; !ok {
		t.Error("expected key DB__HOST after stripping one segment")
	}
	if _, ok := out["DB__PORT"]; !ok {
		t.Error("expected key DB__PORT after stripping one segment")
	}
}

func TestFlatten_Lowercase(t *testing.T) {
	env := map[string]string{"APP__NAME": "envdiff"}
	opts := flattener.DefaultOptions()
	opts.Lowercase = true
	out, err := flattener.Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["app__name"] != "envdiff" {
		t.Errorf("expected app__name=envdiff, got %v", out)
	}
}

func TestFlatten_CollisionLastWins(t *testing.T) {
	env := map[string]string{
		"A__KEY": "first",
		"B__KEY": "second",
	}
	opts := flattener.DefaultOptions()
	opts.StripDepth = 1
	out, err := flattener.Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// sorted order: A__KEY < B__KEY, so B__KEY's value wins
	if out["KEY"] != "second" {
		t.Errorf("expected KEY=second (last sorted wins), got %q", out["KEY"])
	}
}

func TestFlatten_NilEnvReturnsError(t *testing.T) {
	_, err := flattener.Flatten(nil, flattener.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestFlatten_EmptySeparatorReturnsError(t *testing.T) {
	opts := flattener.DefaultOptions()
	opts.Separator = ""
	_, err := flattener.Flatten(map[string]string{"K": "v"}, opts)
	if err == nil {
		t.Error("expected error for empty separator")
	}
}
