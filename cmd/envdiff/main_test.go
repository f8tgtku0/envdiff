package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary compiles the envdiff binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	bin := filepath.Join(tmpDir, "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestCLI_NoDiff_ExitZero(t *testing.T) {
	bin := buildBinary(t)
	left := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	right := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	cmd := exec.Command(bin, left, right)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
}

func TestCLI_WithDiff_ExitTwo(t *testing.T) {
	bin := buildBinary(t)
	left := writeTempEnv(t, "FOO=bar\nONLY_LEFT=1\n")
	right := writeTempEnv(t, "FOO=different\nONLY_RIGHT=1\n")

	cmd := exec.Command(bin, left, right)
	out, err := cmd.CombinedOutput()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 2 {
		t.Fatalf("expected exit 2, got: %v\noutput: %s", err, out)
	}
}

func TestCLI_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	left := writeTempEnv(t, "FOO=bar\n")
	right := writeTempEnv(t, "FOO=bar\n")

	cmd := exec.Command(bin, "-format", "json", left, right)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\noutput: %s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected JSON output, got empty")
	}
}

func TestCLI_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 1 {
		t.Fatalf("expected exit 1 for missing args, got: %v", err)
	}
}
