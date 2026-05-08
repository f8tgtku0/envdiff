package sorter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/sorter"
)

func TestSort_Ascending(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"APPLE": "a",
		"MANGO": "m",
	}
	var buf bytes.Buffer
	if err := sorter.Sort(env, sorter.DefaultOptions(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := nonEmpty(strings.Split(buf.String(), "\n"))
	if lines[0] != "APPLE=a" || lines[1] != "MANGO=m" || lines[2] != "ZEBRA=z" {
		t.Errorf("unexpected order: %v", lines)
	}
}

func TestSort_Descending(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"APPLE": "a",
		"MANGO": "m",
	}
	var buf bytes.Buffer
	opts := sorter.Options{Reverse: true}
	if err := sorter.Sort(env, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := nonEmpty(strings.Split(buf.String(), "\n"))
	if lines[0] != "ZEBRA=z" || lines[1] != "MANGO=m" || lines[2] != "APPLE=a" {
		t.Errorf("unexpected order: %v", lines)
	}
}

func TestSort_GroupPrefixes(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PORT": "5432",
		"APP_NAME": "myapp",
	}
	var buf bytes.Buffer
	opts := sorter.Options{GroupPrefixes: true}
	if err := sorter.Sort(env, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := nonEmpty(strings.Split(buf.String(), "\n"))
	// APP group should appear before DB group
	if !strings.HasPrefix(lines[0], "APP_") || !strings.HasPrefix(lines[1], "APP_") {
		t.Errorf("expected APP_ group first, got: %v", lines)
	}
	if !strings.HasPrefix(lines[2], "DB_") || !strings.HasPrefix(lines[3], "DB_") {
		t.Errorf("expected DB_ group second, got: %v", lines)
	}
}

func TestSort_NilEnv(t *testing.T) {
	var buf bytes.Buffer
	err := sorter.Sort(nil, sorter.DefaultOptions(), &buf)
	if err == nil {
		t.Fatal("expected error for nil env, got nil")
	}
}

func TestSort_NilWriter(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	err := sorter.Sort(env, sorter.DefaultOptions(), nil)
	if err == nil {
		t.Fatal("expected error for nil writer, got nil")
	}
}

func TestSort_EmptyEnv(t *testing.T) {
	var buf bytes.Buffer
	if err := sorter.Sort(map[string]string{}, sorter.DefaultOptions(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty env, got %q", buf.String())
	}
}

func nonEmpty(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
