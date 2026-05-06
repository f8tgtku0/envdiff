package exporter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envdiff/internal/exporter"
)

func TestWrite_DotenvStdout(t *testing.T) {
	// Redirect stdout via OutFile to a temp file for assertion.
	tmp := filepath.Join(t.TempDir(), "out.env")
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	err := exporter.Write(vars, exporter.Options{Format: exporter.FormatDotenv, OutFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	content := string(data)
	if !strings.Contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", content)
	}
	if !strings.Contains(content, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", content)
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.sh")
	vars := map[string]string{"HOME": "/home/user"}
	err := exporter.Write(vars, exporter.Options{Format: exporter.FormatExport, OutFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "export HOME=/home/user") {
		t.Errorf("expected export prefix, got: %s", string(data))
	}
}

func TestWrite_QuotedValues(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.env")
	vars := map[string]string{"MSG": "hello world"}
	err := exporter.Write(vars, exporter.Options{Format: exporter.FormatDotenv, OutFile: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.env")
	vars := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	_ = exporter.Write(vars, exporter.Options{Format: exporter.FormatDotenv, OutFile: tmp})
	data, _ := os.ReadFile(tmp)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "AAA=2" || lines[1] != "MMM=3" || lines[2] != "ZZZ=1" {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}

func TestWrite_InvalidOutFile(t *testing.T) {
	vars := map[string]string{"K": "v"}
	err := exporter.Write(vars, exporter.Options{Format: exporter.FormatDotenv, OutFile: "/no/such/dir/out.env"})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
