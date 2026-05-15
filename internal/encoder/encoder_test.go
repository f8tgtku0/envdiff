package encoder_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/encoder"
)

var baseEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"SECRET":   "p@ssw0rd",
}

func TestEncode_ShellFormat(t *testing.T) {
	var sb strings.Builder
	opts := encoder.DefaultOptions()
	if err := encoder.Encode(baseEnv, &sb, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, want := range []string{"APP_ENV=production", "DB_HOST=localhost", "SECRET=p@ssw0rd"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestEncode_ExportFormat(t *testing.T) {
	var sb strings.Builder
	opts := encoder.DefaultOptions()
	opts.Format = encoder.FormatExport
	if err := encoder.Encode(baseEnv, &sb, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestEncode_DockerFormat(t *testing.T) {
	var sb strings.Builder
	opts := encoder.DefaultOptions()
	opts.Format = encoder.FormatDocker
	if err := encoder.Encode(map[string]string{"FOO": "bar"}, &sb, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := sb.String(); !strings.Contains(got, "--env FOO=bar") {
		t.Errorf("expected docker flag format, got: %s", got)
	}
}

func TestEncode_InlineFormat(t *testing.T) {
	var sb strings.Builder
	opts := encoder.DefaultOptions()
	opts.Format = encoder.FormatInline
	if err := encoder.Encode(map[string]string{"A": "1", "B": "2"}, &sb, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(sb.String())
	if strings.Contains(out, "\n") {
		t.Errorf("inline format should be single line, got:\n%s", out)
	}
	if !strings.Contains(out, "A=1") || !strings.Contains(out, "B=2") {
		t.Errorf("missing key-value pairs in inline output: %s", out)
	}
}

func TestEncode_QuotedValues(t *testing.T) {
	var sb strings.Builder
	opts := encoder.DefaultOptions()
	opts.Quote = true
	if err := encoder.Encode(map[string]string{"KEY": "hello world"}, &sb, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := sb.String(); !strings.Contains(got, `KEY="hello world"`) {
		t.Errorf("expected quoted value, got: %s", got)
	}
}

func TestEncode_NilEnvReturnsError(t *testing.T) {
	var sb strings.Builder
	if err := encoder.Encode(nil, &sb, encoder.DefaultOptions()); err == nil {
		t.Error("expected error for nil env, got nil")
	}
}

func TestEncode_NilWriterReturnsError(t *testing.T) {
	if err := encoder.Encode(baseEnv, nil, encoder.DefaultOptions()); err == nil {
		t.Error("expected error for nil writer, got nil")
	}
}
