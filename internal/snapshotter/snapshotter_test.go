package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/snapshotter"
)

func TestTake_Basic(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "DB_HOST": "localhost"}
	opts := snapshotter.DefaultOptions()
	opts.Label = "test"

	s, err := snapshotter.Take(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Label != "test" {
		t.Errorf("expected label 'test', got %q", s.Label)
	}
	if len(s.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(s.Env))
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestTake_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"KEY": "original"}
	s, _ := snapshotter.Take(env, snapshotter.DefaultOptions())
	s.Env["KEY"] = "mutated"
	if env["KEY"] != "original" {
		t.Error("Take mutated the original env map")
	}
}

func TestTake_NilEnvReturnsError(t *testing.T) {
	_, err := snapshotter.Take(nil, snapshotter.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env, got nil")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"PORT": "8080", "LOG_LEVEL": "debug"}
	opts := snapshotter.DefaultOptions()
	opts.Label = "roundtrip"

	s, err := snapshotter.Take(env, opts)
	if err != nil {
		t.Fatalf("Take: %v", err)
	}
	s.Timestamp = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	path := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshotter.Save(s, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshotter.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Label != "roundtrip" {
		t.Errorf("label mismatch: got %q", loaded.Label)
	}
	if loaded.Env["PORT"] != "8080" {
		t.Errorf("PORT mismatch: got %q", loaded.Env["PORT"])
	}
	if !loaded.Timestamp.Equal(s.Timestamp) {
		t.Errorf("timestamp mismatch: got %v", loaded.Timestamp)
	}
}

func TestSave_NilSnapshotReturnsError(t *testing.T) {
	err := snapshotter.Save(nil, filepath.Join(t.TempDir(), "snap.json"))
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestLoad_MissingFileReturnsError(t *testing.T) {
	_, err := snapshotter.Load(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSONReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := snapshotter.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
