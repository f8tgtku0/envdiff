package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/scanner"
)

// scaffold creates a temporary directory tree and returns its root path.
func scaffold(t *testing.T, files []string) string {
	t.Helper()
	root := t.TempDir()
	for _, f := range files {
		full := filepath.Join(root, f)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(full, []byte("KEY=val\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", f, err)
		}
	}
	return root
}

func TestScan_DefaultPatterns(t *testing.T) {
	root := scaffold(t, []string{".env", ".env.local", "config.yaml", "README.md"})
	opts := scanner.DefaultOptions()
	got, err := scanner.Scan(root, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 results, got %d: %v", len(got), got)
	}
}

func TestScan_CustomPattern(t *testing.T) {
	root := scaffold(t, []string{"prod.env", "dev.env", "notes.txt"})
	opts := scanner.Options{Patterns: []string{"*.env"}}
	got, err := scanner.Scan(root, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 results, got %d", len(got))
	}
}

func TestScan_RecursiveFalse(t *testing.T) {
	root := scaffold(t, []string{".env", "sub/.env.staging"})
	opts := scanner.Options{Patterns: []string{".env*"}, Recursive: false}
	got, err := scanner.Scan(root, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 result (non-recursive), got %d: %v", len(got), got)
	}
}

func TestScan_RecursiveTrue(t *testing.T) {
	root := scaffold(t, []string{".env", "sub/.env.staging", "sub/deep/.env.prod"})
	opts := scanner.Options{Patterns: []string{".env*"}, Recursive: true}
	got, err := scanner.Scan(root, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 results, got %d: %v", len(got), got)
	}
}

func TestScan_IgnoreDirs(t *testing.T) {
	root := scaffold(t, []string{".env", "vendor/.env", ".git/.env"})
	opts := scanner.Options{
		Patterns:   []string{".env*"},
		Recursive:  true,
		IgnoreDirs: []string{"vendor", ".git"},
	}
	got, err := scanner.Scan(root, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 result (ignored dirs), got %d: %v", len(got), got)
	}
}

func TestScan_EmptyDir(t *testing.T) {
	root := t.TempDir()
	opts := scanner.DefaultOptions()
	got, err := scanner.Scan(root, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want 0 results for empty dir, got %d", len(got))
	}
}
