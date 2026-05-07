package templater_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envdiff/internal/templater"
)

func TestGenerate_BasicOutput(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
	}
	var buf strings.Builder
	err := templater.Generate(&buf, env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "APP_ENV=") {
		t.Errorf("expected APP_ENV= in output, got:\n%s", got)
	}
	if !strings.Contains(got, "DB_HOST=") {
		t.Errorf("expected DB_HOST= in output, got:\n%s", got)
	}
	if strings.Contains(got, "production") {
		t.Errorf("original value should not appear in template output")
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc123"}
	var buf strings.Builder
	opts := templater.Options{Placeholder: "CHANGEME", CommentValues: false}
	if err := templater.Generate(&buf, env, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SECRET_KEY=CHANGEME") {
		t.Errorf("expected placeholder in output, got: %s", buf.String())
	}
}

func TestGenerate_CommentValues(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	var buf strings.Builder
	opts := templater.Options{Placeholder: "", CommentValues: true}
	if err := templater.Generate(&buf, env, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "# was: 8080") {
		t.Errorf("expected inline comment with original value, got: %s", buf.String())
	}
}

func TestGenerate_SortedKeys(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}
	var buf strings.Builder
	if err := templater.Generate(&buf, env, templater.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first line ALPHA, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line ZEBRA, got %s", lines[2])
	}
}

func TestGenerate_NilWriter(t *testing.T) {
	err := templater.Generate(nil, map[string]string{"K": "v"}, templater.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil writer")
	}
}

func TestMerge_DeduplicatesKeys(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "99", "BAZ": "3"}
	merged := templater.Merge(a, b)
	if len(merged) != 3 {
		t.Errorf("expected 3 keys, got %d", len(merged))
	}
	if merged["FOO"] != "1" {
		t.Errorf("expected first-wins value '1', got %s", merged["FOO"])
	}
}
