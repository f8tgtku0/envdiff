package pinner_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/pinner"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV": "staging",
		"DEBUG":   "true",
		"PORT":    "8080",
	}
}

func TestPin_NoPins(t *testing.T) {
	env := baseEnv()
	r, err := pinner.Pin(env, nil, pinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestPin_AddsMissingKey(t *testing.T) {
	env := baseEnv()
	pins := map[string]string{"NEW_KEY": "hello"}
	r, err := pinner.Pin(env, pins, pinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", env["NEW_KEY"])
	}
	if len(r.Pinned()) != 1 {
		t.Errorf("expected 1 pinned, got %d", len(r.Pinned()))
	}
}

func TestPin_SkipsExistingByDefault(t *testing.T) {
	env := baseEnv()
	pins := map[string]string{"APP_ENV": "production"}
	r, err := pinner.Pin(env, pins, pinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV unchanged, got %q", env["APP_ENV"])
	}
	if len(r.Skipped()) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(r.Skipped()))
	}
}

func TestPin_OverwriteExisting(t *testing.T) {
	env := baseEnv()
	pins := map[string]string{"APP_ENV": "production"}
	opts := pinner.DefaultOptions()
	opts.Overwrite = true
	r, err := pinner.Pin(env, pins, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env["APP_ENV"])
	}
	if len(r.Pinned()) != 1 {
		t.Errorf("expected 1 pinned, got %d", len(r.Pinned()))
	}
}

func TestPin_SameValueSkipped(t *testing.T) {
	env := baseEnv()
	pins := map[string]string{"DEBUG": "true"}
	r, err := pinner.Pin(env, pins, pinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Skipped()) != 1 {
		t.Errorf("expected 1 skipped (already matches), got %d", len(r.Skipped()))
	}
	if r.Skipped()[0].Reason != "already matches" {
		t.Errorf("unexpected reason: %q", r.Skipped()[0].Reason)
	}
}

func TestPin_StrictMissingReturnsError(t *testing.T) {
	env := baseEnv()
	pins := map[string]string{"DOES_NOT_EXIST": "val"}
	opts := pinner.DefaultOptions()
	opts.StrictMissing = true
	_, err := pinner.Pin(env, pins, opts)
	if err == nil {
		t.Fatal("expected error for missing key under StrictMissing")
	}
}

func TestPin_NilEnvReturnsError(t *testing.T) {
	_, err := pinner.Pin(nil, map[string]string{"K": "v"}, pinner.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestWriteText_ShowsPinned(t *testing.T) {
	env := baseEnv()
	pins := map[string]string{"APP_ENV": "production"}
	opts := pinner.DefaultOptions()
	opts.Overwrite = true
	r, _ := pinner.Pin(env, pins, opts)

	var buf bytes.Buffer
	if err := pinner.WriteText(&buf, r); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PINNED") {
		t.Errorf("expected PINNED in output, got: %s", out)
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	env := baseEnv()
	pins := map[string]string{"PORT": "9090"}
	opts := pinner.DefaultOptions()
	opts.Overwrite = true
	r, _ := pinner.Pin(env, pins, opts)

	var buf bytes.Buffer
	if err := pinner.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(buf.String(), "PORT") {
		t.Errorf("expected PORT in JSON output")
	}
}
