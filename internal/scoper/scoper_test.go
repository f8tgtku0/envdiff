package scoper_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/scoper"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
		"REDIS_URL": "redis://localhost",
		"APP_NAME":  "envdiff",
		"APP_ENV":   "production",
	}
}

func TestScope_NoRules(t *testing.T) {
	env := baseEnv()
	opts := scoper.DefaultOptions()
	res, err := scoper.Scope(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Scopes) != 0 {
		t.Errorf("expected 0 scopes, got %d", len(res.Scopes))
	}
	if len(res.Unscoped) != len(env) {
		t.Errorf("expected %d unscoped keys, got %d", len(env), len(res.Unscoped))
	}
}

func TestScope_PrefixRules(t *testing.T) {
	env := baseEnv()
	opts := scoper.DefaultOptions()
	opts.Rules = map[string][]string{
		"database": {"DB_"},
		"cache":    {"REDIS_"},
	}
	res, err := scoper.Scope(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Scopes) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(res.Scopes))
	}
	if len(res.Unscoped) != 2 {
		t.Errorf("expected 2 unscoped keys, got %d", len(res.Unscoped))
	}
}

func TestScope_ScopedKeysAreSorted(t *testing.T) {
	env := baseEnv()
	opts := scoper.DefaultOptions()
	opts.Rules = map[string][]string{"database": {"DB_"}}
	res, _ := scoper.Scope(env, opts)
	if res.Scopes[0].Keys[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST first, got %s", res.Scopes[0].Keys[0])
	}
}

func TestScope_IncludeUnscopedFalse(t *testing.T) {
	env := baseEnv()
	opts := scoper.DefaultOptions()
	opts.IncludeUnscoped = false
	opts.Rules = map[string][]string{"database": {"DB_"}}
	res, _ := scoper.Scope(env, opts)
	if len(res.Unscoped) != 0 {
		t.Errorf("expected no unscoped keys, got %d", len(res.Unscoped))
	}
}

func TestScope_NilEnv(t *testing.T) {
	_, err := scoper.Scope(nil, scoper.DefaultOptions())
	if err != nil {
		t.Errorf("nil env should not error, got %v", err)
	}
}

func TestWriteText_ContainsScopeName(t *testing.T) {
	env := baseEnv()
	opts := scoper.DefaultOptions()
	opts.Rules = map[string][]string{"database": {"DB_"}}
	res, _ := scoper.Scope(env, opts)

	var buf bytes.Buffer
	if err := scoper.WriteText(&buf, res, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[database]") {
		t.Errorf("expected [database] section in output")
	}
}

func TestWriteJSON_ValidOutput(t *testing.T) {
	env := baseEnv()
	opts := scoper.DefaultOptions()
	opts.Rules = map[string][]string{"cache": {"REDIS_"}}
	res, _ := scoper.Scope(env, opts)

	var buf bytes.Buffer
	if err := scoper.WriteJSON(&buf, res, env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"cache\"") {
		t.Errorf("expected cache scope in JSON output")
	}
}
