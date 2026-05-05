package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/filter"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.local",
		"DB_PASSWORD": "secret",
		"LOG_LEVEL":   "info",
	}
}

func TestApply_NoFilter(t *testing.T) {
	env := baseEnv()
	result := filter.Apply(env, filter.Options{})
	if len(result) != len(env) {
		t.Fatalf("expected %d keys, got %d", len(env), len(result))
	}
}

func TestApply_IncludePrefix(t *testing.T) {
	result := filter.Apply(baseEnv(), filter.Options{
		IncludePrefixes: []string{"APP_"},
	})
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST to be present")
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be absent")
	}
}

func TestApply_ExcludePrefix(t *testing.T) {
	result := filter.Apply(baseEnv(), filter.Options{
		ExcludePrefixes: []string{"DB_"},
	})
	if _, ok := result["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be excluded")
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be excluded")
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 keys after exclusion, got %d", len(result))
	}
}

func TestApply_IncludeAndExclude(t *testing.T) {
	// Include APP_ and DB_, then exclude DB_PASSWORD
	result := filter.Apply(baseEnv(), filter.Options{
		IncludePrefixes: []string{"APP_", "DB_"},
		ExcludePrefixes: []string{"DB_PASSWORD"},
	})
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be excluded")
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be present")
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result))
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	result := filter.Apply(map[string]string{}, filter.Options{
		IncludePrefixes: []string{"APP_"},
	})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d keys", len(result))
	}
}
