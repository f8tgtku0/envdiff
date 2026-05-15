package differ

import (
	"testing"
)

var leftEnv = map[string]string{
	"DB_HOST":  "localhost",
	"DB_PORT":  "5432",
	"APP_ENV":  "production",
	"API_KEY":  "abc123",
	"LOG_LEVEL": "INFO",
}

var rightEnv = map[string]string{
	"DB_HOST":  "localhost",
	"DB_PORT":  "5433",
	"APP_ENV":  "PRODUCTION",
	"API_KEY":  "xyz789",
	"LOG_LEVEL": "info",
}

func TestMatchValues_ExactMatch(t *testing.T) {
	results := MatchValues(leftEnv, rightEnv, MatchOptions{})
	conformCount := 0
	for _, r := range results {
		if r.Key == "DB_HOST" && !r.Conformant {
			t.Errorf("DB_HOST should be conformant (same value)")
		}
		if r.Conformant {
			conformCount++
		}
	}
	if conformCount != 1 {
		t.Errorf("expected 1 conformant key (DB_HOST), got %d", conformCount)
	}
}

func TestMatchValues_IgnoreCase(t *testing.T) {
	results := MatchValues(leftEnv, rightEnv, MatchOptions{IgnoreCase: true})
	for _, r := range results {
		switch r.Key {
		case "DB_HOST", "APP_ENV", "LOG_LEVEL":
			if !r.Conformant {
				t.Errorf("key %s should be conformant with IgnoreCase", r.Key)
			}
		case "DB_PORT", "API_KEY":
			if r.Conformant {
				t.Errorf("key %s should NOT be conformant even with IgnoreCase", r.Key)
			}
		}
	}
}

func TestMatchValues_PatternOverride(t *testing.T) {
	opts := MatchOptions{
		ValuePatterns: map[string]string{
			"API_KEY": `^[a-z0-9]+$`,
		},
	}
	results := MatchValues(leftEnv, rightEnv, opts)
	for _, r := range results {
		if r.Key == "API_KEY" {
			if !r.Conformant {
				t.Errorf("API_KEY both match pattern, expected Conformant=true")
			}
			if r.PatternUsed == "" {
				t.Errorf("expected PatternUsed to be set")
			}
		}
	}
}

func TestMatchValues_InvalidPattern(t *testing.T) {
	opts := MatchOptions{
		ValuePatterns: map[string]string{
			"DB_PORT": `[invalid`,
		},
	}
	// Should not panic; invalid patterns are silently skipped.
	results := MatchValues(leftEnv, rightEnv, opts)
	for _, r := range results {
		if r.Key == "DB_PORT" && r.PatternUsed != "" {
			t.Errorf("invalid pattern should not be recorded in PatternUsed")
		}
	}
}

func TestMatchValues_OnlySharedKeys(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "2", "C": "3"}
	results := MatchValues(left, right, MatchOptions{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result for shared key B, got %d", len(results))
	}
	if results[0].Key != "B" {
		t.Errorf("expected key B, got %s", results[0].Key)
	}
}
