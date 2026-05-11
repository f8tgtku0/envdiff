package grouper_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/grouper"
)

func TestWriteText_ContainsPrefix(t *testing.T) {
	res := grouper.GroupEnv(baseEnv(), grouper.DefaultOptions())
	var buf strings.Builder
	if err := grouper.WriteText(res, &buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected [DB] in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[APP]") {
		t.Errorf("expected [APP] in output, got:\n%s", out)
	}
}

func TestWriteText_UngroupedSection(t *testing.T) {
	res := grouper.GroupEnv(baseEnv(), grouper.DefaultOptions())
	var buf strings.Builder
	_ = grouper.WriteText(res, &buf)
	out := buf.String()
	if !strings.Contains(out, "[ungrouped]") {
		t.Errorf("expected [ungrouped] section, got:\n%s", out)
	}
	if !strings.Contains(out, "STANDALONE") {
		t.Errorf("expected STANDALONE in ungrouped, got:\n%s", out)
	}
}

func TestWriteText_EmptyResult(t *testing.T) {
	res := grouper.GroupEnv(map[string]string{}, grouper.DefaultOptions())
	var buf strings.Builder
	_ = grouper.WriteText(res, &buf)
	if buf.String() != "" {
		t.Errorf("expected empty output, got: %q", buf.String())
	}
}

func TestWriteMarkdown_Headers(t *testing.T) {
	res := grouper.GroupEnv(baseEnv(), grouper.DefaultOptions())
	var buf strings.Builder
	if err := grouper.WriteMarkdown(res, &buf); err != nil {
		t.Fatalf("WriteMarkdown error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "## DB") {
		t.Errorf("expected ## DB header, got:\n%s", out)
	}
	if !strings.Contains(out, "| Key | Value |") {
		t.Errorf("expected table header, got:\n%s", out)
	}
}

func TestWriteMarkdown_KeyValues(t *testing.T) {
	res := grouper.GroupEnv(map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}, grouper.DefaultOptions())
	var buf strings.Builder
	_ = grouper.WriteMarkdown(res, &buf)
	out := buf.String()
	if !strings.Contains(out, "`DB_HOST`") {
		t.Errorf("expected DB_HOST in markdown, got:\n%s", out)
	}
	if !strings.Contains(out, "`localhost`") {
		t.Errorf("expected localhost value in markdown, got:\n%s", out)
	}
}
