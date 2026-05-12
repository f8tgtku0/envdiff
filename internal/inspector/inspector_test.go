package inspector_test

import (
	"testing"

	"github.com/user/envdiff/internal/inspector"
)

func makeEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"APP_NAME": "myapp",
			"DB_HOST":  "localhost",
			"DEBUG":    "true",
		},
		"staging": {
			"APP_NAME": "myapp",
			"DB_HOST":  "staging.db.example.com",
			"LOG_LEVEL": "info",
		},
		"prod": {
			"APP_NAME": "myapp",
			"DB_HOST":  "prod.db.example.com",
			"LOG_LEVEL": "warn",
		},
	}
}

func TestInspect_AllKeys(t *testing.T) {
	result := inspector.Inspect(makeEnvs(), inspector.DefaultOptions())
	if len(result.Keys) != 4 {
		t.Fatalf("expected 4 keys, got %d", len(result.Keys))
	}
}

func TestInspect_ConsistentFlag(t *testing.T) {
	result := inspector.Inspect(makeEnvs(), inspector.DefaultOptions())
	for _, info := range result.Keys {
		switch info.Key {
		case "APP_NAME":
			if !info.Consistent {
				t.Errorf("APP_NAME should be consistent")
			}
		case "DB_HOST":
			if info.Consistent {
				t.Errorf("DB_HOST should not be consistent")
			}
		}
	}
}

func TestInspect_OnlyInconsistent(t *testing.T) {
	opts := inspector.Options{OnlyInconsistent: true}
	result := inspector.Inspect(makeEnvs(), opts)
	for _, info := range result.Keys {
		if info.Consistent {
			t.Errorf("key %q should not appear in inconsistent-only result", info.Key)
		}
	}
}

func TestInspect_EnvsList(t *testing.T) {
	result := inspector.Inspect(makeEnvs(), inspector.DefaultOptions())
	for _, info := range result.Keys {
		if info.Key == "DEBUG" {
			if len(info.Envs) != 1 || info.Envs[0] != "dev" {
				t.Errorf("DEBUG should only appear in dev, got %v", info.Envs)
			}
		}
	}
}

func TestInspect_EmptyInput(t *testing.T) {
	result := inspector.Inspect(map[string]map[string]string{}, inspector.DefaultOptions())
	if len(result.Keys) != 0 {
		t.Errorf("expected 0 keys for empty input")
	}
}

func TestInspect_SortedKeys(t *testing.T) {
	result := inspector.Inspect(makeEnvs(), inspector.DefaultOptions())
	for i := 1; i < len(result.Keys); i++ {
		if result.Keys[i].Key < result.Keys[i-1].Key {
			t.Errorf("keys not sorted: %q before %q", result.Keys[i-1].Key, result.Keys[i].Key)
		}
	}
}
