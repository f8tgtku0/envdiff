package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestLoad_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	env, err := loader.Load(path, loader.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Path != path {
		t.Errorf("path: got %q, want %q", env.Path, path)
	}
	if env.Vars["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", env.Vars["FOO"], "bar")
	}
	if env.Vars["BAZ"] != "qux" {
		t.Errorf("BAZ: got %q, want %q", env.Vars["BAZ"], "qux")
	}
}

func TestLoad_WithIncludeFilter(t *testing.T) {
	path := writeTempEnv(t, "APP_HOST=localhost\nDB_HOST=db\nAPP_PORT=8080\n")

	env, err := loader.Load(path, loader.Options{
		IncludePrefixes: []string{"APP_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env.Vars["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be excluded")
	}
	if env.Vars["APP_HOST"] != "localhost" {
		t.Errorf("APP_HOST: got %q, want %q", env.Vars["APP_HOST"], "localhost")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := loader.Load("/nonexistent/.env", loader.Options{})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadPair_BothFiles(t *testing.T) {
	leftPath := writeTempEnv(t, "KEY=left\n")
	rightPath := writeTempEnv(t, "KEY=right\n")

	left, right, err := loader.LoadPair(leftPath, rightPath, loader.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if left.Vars["KEY"] != "left" {
		t.Errorf("left KEY: got %q, want %q", left.Vars["KEY"], "left")
	}
	if right.Vars["KEY"] != "right" {
		t.Errorf("right KEY: got %q, want %q", right.Vars["KEY"], "right")
	}
}

func TestLoadPair_RightFileMissing(t *testing.T) {
	leftPath := writeTempEnv(t, "KEY=val\n")

	_, _, err := loader.LoadPair(leftPath, "/no/such/file", loader.Options{})
	if err == nil {
		t.Fatal("expected error when right file is missing")
	}
}
