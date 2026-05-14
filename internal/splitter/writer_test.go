package splitter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/splitter"
)

func makeResult() *splitter.Result {
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB", "APP"}
	r, _ := splitter.Split(baseEnv(), opts)
	return r
}

func TestWriteText_ContainsBucketHeaders(t *testing.T) {
	var buf bytes.Buffer
	err := splitter.WriteText(makeResult(), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, label := range []string{"[DB]", "[APP]", "[(unmatched)]"}  {
		if !strings.Contains(out, label) {
			t.Errorf("expected output to contain %q", label)
		}
	}
}

func TestWriteText_NilWriter(t *testing.T) {
	err := splitter.WriteText(makeResult(), nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestWriteText_NilResult(t *testing.T) {
	var buf bytes.Buffer
	err := splitter.WriteText(nil, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no result") {
		t.Errorf("expected '(no result)' message")
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	err := splitter.WriteJSON(makeResult(), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB") || !strings.Contains(out, "APP") {
		t.Errorf("expected JSON to contain DB and APP keys")
	}
}

func TestWriteJSON_NilWriter(t *testing.T) {
	err := splitter.WriteJSON(makeResult(), nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}
