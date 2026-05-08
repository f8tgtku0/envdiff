package duplicates_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/duplicates"
)

func TestDetectCross_NoDuplicates(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"DB_HOST": "prod-db", "PORT": "5432"},
		"staging": {"CACHE_URL": "redis://staging"},
	}

	result := duplicates.DetectCross(envs)
	if len(result.CrossEnv) != 0 {
		t.Fatalf("expected no cross-env duplicates, got %d", len(result.CrossEnv))
	}
}

func TestDetectCross_SingleSharedKey(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"DB_HOST": "prod-db", "PORT": "5432"},
		"staging": {"DB_HOST": "staging-db", "CACHE_URL": "redis://"},
	}

	result := duplicates.DetectCross(envs)
	if len(result.CrossEnv) != 1 {
		t.Fatalf("expected 1 cross-env duplicate, got %d", len(result.CrossEnv))
	}
	if result.CrossEnv[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", result.CrossEnv[0].Key)
	}
	if len(result.CrossEnv[0].Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(result.CrossEnv[0].Sources))
	}
}

func TestDetectCross_MultipleSharedKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"DB_HOST": "prod-db", "PORT": "5432", "LOG_LEVEL": "warn"},
		"staging": {"DB_HOST": "staging-db", "PORT": "5433", "LOG_LEVEL": "debug"},
		"dev":     {"DB_HOST": "localhost", "DEBUG": "true"},
	}

	result := duplicates.DetectCross(envs)
	// DB_HOST appears in all three; PORT and LOG_LEVEL appear in two.
	if len(result.CrossEnv) != 3 {
		t.Fatalf("expected 3 cross-env duplicates, got %d", len(result.CrossEnv))
	}

	keys := make(map[string]int)
	for _, e := range result.CrossEnv {
		keys[e.Key] = len(e.Sources)
	}

	if keys["DB_HOST"] != 3 {
		t.Errorf("DB_HOST should appear in 3 sources, got %d", keys["DB_HOST"])
	}
	if keys["PORT"] != 2 {
		t.Errorf("PORT should appear in 2 sources, got %d", keys["PORT"])
	}
	if keys["LOG_LEVEL"] != 2 {
		t.Errorf("LOG_LEVEL should appear in 2 sources, got %d", keys["LOG_LEVEL"])
	}
}

func TestDetectCross_EmptyInput(t *testing.T) {
	result := duplicates.DetectCross(map[string]map[string]string{})
	if len(result.CrossEnv) != 0 {
		t.Fatalf("expected empty result for empty input")
	}
}

func TestDetectCross_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"DB_HOST": "prod-db", "PORT": "5432"},
	}

	result := duplicates.DetectCross(envs)
	if len(result.CrossEnv) != 0 {
		t.Fatalf("single env cannot have cross-env duplicates")
	}
}
