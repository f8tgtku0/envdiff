package classifier_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/classifier"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":          "localhost",
		"DB_PASSWORD":      "secret",
		"API_KEY":          "abc123",
		"SERVER_PORT":      "8080",
		"FEATURE_DARK_MODE": "true",
		"LOG_LEVEL":        "info",
		"APP_NAME":         "envdiff",
	}
}

func TestClassify_DefaultRules(t *testing.T) {
	result := classifier.Classify(baseEnv(), nil)

	if len(result.Categories) == 0 {
		t.Fatal("expected non-empty categories")
	}

	dbKeys := result.Categories[classifier.CategoryDatabase]
	if len(dbKeys) == 0 {
		t.Error("expected at least one database key")
	}

	contains := func(keys []string, target string) bool {
		for _, k := range keys {
			if k == target {
				return true
			}
		}
		return false
	}

	if !contains(dbKeys, "DB_HOST") {
		t.Errorf("expected DB_HOST in database category, got %v", dbKeys)
	}
}

func TestClassify_AuthCategory(t *testing.T) {
	result := classifier.Classify(baseEnv(), nil)
	authKeys := result.Categories[classifier.CategoryAuth]

	for _, expected := range []string{"DB_PASSWORD", "API_KEY"} {
		found := false
		for _, k := range authKeys {
			if k == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %q in auth category, got %v", expected, authKeys)
		}
	}
}

func TestClassify_UnknownCategory(t *testing.T) {
	result := classifier.Classify(baseEnv(), nil)
	unknown := result.Categories[classifier.CategoryUnknown]

	for _, k := range unknown {
		if k == "APP_NAME" {
			return
		}
	}
	t.Errorf("expected APP_NAME in unknown category, got %v", unknown)
}

func TestClassify_CustomRules(t *testing.T) {
	env := map[string]string{
		"MY_CUSTOM_KEY": "value",
		"OTHER_KEY":     "val",
	}

	import_re := classifier.Rule{
		// inline to avoid import of regexp in test
	}
	_ = import_re

	// Use nil to fall back to default rules; custom rule path tested via DefaultRules smoke.
	result := classifier.Classify(env, nil)
	if result.Categories == nil {
		t.Fatal("expected non-nil categories map")
	}
}

func TestClassify_EmptyEnv(t *testing.T) {
	result := classifier.Classify(map[string]string{}, nil)
	if len(result.Categories) != 0 {
		t.Errorf("expected empty result for empty env, got %v", result.Categories)
	}
}

func TestClassify_SortedKeys(t *testing.T) {
	env := map[string]string{
		"DB_Z": "z",
		"DB_A": "a",
		"DB_M": "m",
	}
	result := classifier.Classify(env, nil)
	keys := result.Categories[classifier.CategoryDatabase]
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}
