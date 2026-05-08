package trimmer_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/trimmer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":  "myapp",
		"APP_PORT":  "8080",
		"DB_HOST":   "localhost",
		"DB_PASS":   "secret",
		"LEGACY_V1": "old",
	}
}

func TestTrim_NoUnused(t *testing.T) {
	env := baseEnv()
	ref := baseEnv()
	res := trimmer.Trim(env, ref, trimmer.DefaultOptions())
	if len(res.Removed) != 0 {
		t.Fatalf("expected no removed keys, got %v", res.Removed)
	}
	if len(res.Kept) != len(env) {
		t.Fatalf("expected kept length %d, got %d", len(env), len(res.Kept))
	}
}

func TestTrim_RemovesUnused(t *testing.T) {
	env := baseEnv()
	ref := map[string]string{
		"APP_NAME": "myapp",
		"APP_PORT": "8080",
	}
	res := trimmer.Trim(env, ref, trimmer.DefaultOptions())
	if len(res.Removed) != 3 {
		t.Fatalf("expected 3 removed, got %d: %v", len(res.Removed), res.Removed)
	}
	if _, ok := res.Kept["DB_HOST"]; ok {
		t.Error("DB_HOST should have been removed from Kept")
	}
}

func TestTrim_DryRunDoesNotMutate(t *testing.T) {
	env := baseEnv()
	ref := map[string]string{"APP_NAME": "myapp"}
	opts := trimmer.DefaultOptions()
	opts.DryRun = true
	res := trimmer.Trim(env, ref, opts)
	if res.Kept != nil {
		t.Error("expected Kept to be nil in dry-run mode")
	}
	if len(res.Removed) == 0 {
		t.Error("expected Removed to be populated in dry-run mode")
	}
	// original env must be untouched
	if len(env) != 5 {
		t.Errorf("original env mutated; got length %d", len(env))
	}
}

func TestTrim_IgnorePrefix(t *testing.T) {
	env := baseEnv()
	ref := map[string]string{
		"APP_NAME": "myapp",
		"APP_PORT": "8080",
	}
	opts := trimmer.DefaultOptions()
	opts.IgnorePrefix = []string{"DB_", "LEGACY_"}
	res := trimmer.Trim(env, ref, opts)
	if len(res.Removed) != 0 {
		t.Fatalf("expected 0 removed (DB_ and LEGACY_ ignored), got %v", res.Removed)
	}
}

func TestTrim_EmptyEnv(t *testing.T) {
	res := trimmer.Trim(map[string]string{}, baseEnv(), trimmer.DefaultOptions())
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals from empty env, got %v", res.Removed)
	}
}

func TestTrim_RemovedIsSorted(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	res := trimmer.Trim(env, map[string]string{}, trimmer.DefaultOptions())
	for i := 1; i < len(res.Removed); i++ {
		if res.Removed[i-1] > res.Removed[i] {
			t.Errorf("Removed not sorted: %v", res.Removed)
		}
	}
}
