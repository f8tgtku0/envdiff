package profiler_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/profiler"
)

func makeEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"APP_ENV":  "development",
			"DB_HOST":  "localhost",
			"API_KEY":  "dev-key",
		},
		"staging": {
			"APP_ENV":  "staging",
			"DB_HOST":  "localhost",
			"API_KEY":  "staging-key",
			"LOG_LEVEL": "info",
		},
		"prod": {
			"APP_ENV": "production",
			"DB_HOST": "db.prod.internal",
			"API_KEY": "prod-key",
		},
	}
}

func TestProfile_KeyCount(t *testing.T) {
	profiles := profiler.Profile(makeEnvs())
	// Keys: APP_ENV, DB_HOST, API_KEY, LOG_LEVEL = 4
	if len(profiles) != 4 {
		t.Fatalf("expected 4 profiles, got %d", len(profiles))
	}
}

func TestProfile_Sorted(t *testing.T) {
	profiles := profiler.Profile(makeEnvs())
	for i := 1; i < len(profiles); i++ {
		if profiles[i].Key < profiles[i-1].Key {
			t.Errorf("profiles not sorted: %s before %s", profiles[i-1].Key, profiles[i].Key)
		}
	}
}

func TestProfile_MissingIn(t *testing.T) {
	profiles := profiler.Profile(makeEnvs())
	var logProfile *profiler.KeyProfile
	for i := range profiles {
		if profiles[i].Key == "LOG_LEVEL" {
			logProfile = &profiles[i]
			break
		}
	}
	if logProfile == nil {
		t.Fatal("LOG_LEVEL profile not found")
	}
	if len(logProfile.MissingIn) != 2 {
		t.Errorf("expected LOG_LEVEL missing in 2 envs, got %d: %v", len(logProfile.MissingIn), logProfile.MissingIn)
	}
}

func TestProfile_UniqueFlag(t *testing.T) {
	profiles := profiler.Profile(makeEnvs())
	for _, p := range profiles {
		if p.Key == "DB_HOST" {
			// dev and staging share "localhost", prod differs — not unique
			if p.Unique {
				t.Error("DB_HOST should not be unique across all envs")
			}
			return
		}
	}
	t.Fatal("DB_HOST profile not found")
}

func TestSummary_Consistent(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080"},
		"prod": {"PORT": "8080"},
	}
	profiles := profiler.Profile(envs)
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	summary := profiler.Summary(profiles[0])
	if summary != "PORT: consistent across all environments" {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestSummary_WithDiffs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"SECRET": "abc"},
		"prod": {"SECRET": "xyz"},
	}
	profiles := profiler.Profile(envs)
	summary := profiler.Summary(profiles[0])
	if summary == "SECRET: consistent across all environments" {
		t.Error("expected diff summary but got consistent")
	}
}
