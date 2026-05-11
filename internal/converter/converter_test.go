package converter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/converter"
)

var sampleEnv = map[string]string{
	"APP_NAME": "myapp",
	"DB_URL":   "postgres://localhost/dev",
	"DEBUG":    "true",
	"GREETING": "hello world",
}

func TestConvert_Dotenv(t *testing.T) {
	var sb strings.Builder
	if err := converter.Convert(sampleEnv, converter.FormatDotenv, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got:\n%s", out)
	}
	// Value with a space should be quoted.
	if !strings.Contains(out, `GREETING=`) {
		t.Errorf("expected GREETING line in output")
	}
}

func TestConvert_Export(t *testing.T) {
	var sb strings.Builder
	if err := converter.Convert(sampleEnv, converter.FormatExport, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if !strings.HasPrefix(line, "export ") {
			t.Errorf("expected 'export ' prefix, got: %s", line)
		}
	}
}

func TestConvert_JSON(t *testing.T) {
	var sb strings.Builder
	if err := converter.Convert(sampleEnv, converter.FormatJSON, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, `"APP_NAME"`) {
		t.Errorf("expected JSON key APP_NAME, got:\n%s", out)
	}
	if !strings.Contains(out, `"myapp"`) {
		t.Errorf("expected JSON value myapp, got:\n%s", out)
	}
}

func TestConvert_YAML(t *testing.T) {
	var sb strings.Builder
	if err := converter.Convert(sampleEnv, converter.FormatYAML, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_NAME: myapp") {
		t.Errorf("expected YAML line, got:\n%s", out)
	}
}

func TestConvert_UnknownFormat(t *testing.T) {
	var sb strings.Builder
	err := converter.Convert(sampleEnv, converter.Format("xml"), &sb)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestParseFormat_Valid(t *testing.T) {
	for _, f := range converter.ValidFormats {
		got, err := converter.ParseFormat(string(f))
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", f, err)
		}
		if got != f {
			t.Errorf("ParseFormat(%q) = %q, want %q", f, got, f)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := converter.ParseFormat("toml")
	if err == nil {
		t.Fatal("expected error for unsupported format 'toml'")
	}
}
