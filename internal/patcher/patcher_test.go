package patcher_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/patcher"
)

func TestApply_AddNewKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	patch := map[string]string{"B": "2"}

	out, changes, err := patcher.Apply(base, patch, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %q", out["B"])
	}
	if len(changes) != 1 || changes[0].Action != "add" {
		t.Errorf("expected one 'add' change, got %+v", changes)
	}
}

func TestApply_SkipExistingByDefault(t *testing.T) {
	base := map[string]string{"A": "original"}
	patch := map[string]string{"A": "new"}

	out, changes, err := patcher.Apply(base, patch, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "original" {
		t.Errorf("expected A to remain 'original', got %q", out["A"])
	}
	if len(changes) != 1 || changes[0].Action != "skip" {
		t.Errorf("expected one 'skip' change, got %+v", changes)
	}
}

func TestApply_OverwriteExisting(t *testing.T) {
	base := map[string]string{"A": "original"}
	patch := map[string]string{"A": "new"}
	opts := patcher.Options{Overwrite: true}

	out, changes, err := patcher.Apply(base, patch, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %q", out["A"])
	}
	if len(changes) != 1 || changes[0].Action != "update" {
		t.Errorf("expected one 'update' change, got %+v", changes)
	}
}

func TestApply_DryRunDoesNotMutate(t *testing.T) {
	base := map[string]string{"A": "1"}
	patch := map[string]string{"B": "2"}
	opts := patcher.Options{DryRun: true}

	out, changes, err := patcher.Apply(base, patch, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["B"]; ok {
		t.Error("dry-run should not add B to output map")
	}
	if len(changes) != 1 || changes[0].Action != "add" {
		t.Errorf("expected change recorded even in dry-run, got %+v", changes)
	}
}

func TestApply_NilBaseReturnsError(t *testing.T) {
	_, _, err := patcher.Apply(nil, map[string]string{}, patcher.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil base")
	}
}

func TestWritePatched(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "hello world"}
	tmp := filepath.Join(t.TempDir(), ".env")

	if err := patcher.WritePatched(tmp, env); err != nil {
		t.Fatalf("WritePatched error: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("expected non-empty file")
	}
}
