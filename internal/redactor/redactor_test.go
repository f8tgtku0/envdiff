package redactor_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/redactor"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	r := redactor.New(nil, "")

	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY", "APP_SECRET"}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}

	safe := []string{"APP_NAME", "PORT", "LOG_LEVEL", "DATABASE_HOST"}
	for _, key := range safe {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	r := redactor.New([]string{"INTERNAL", "CERT"}, "")

	if !r.IsSensitive("INTERNAL_URL") {
		t.Error("expected INTERNAL_URL to be sensitive")
	}
	if !r.IsSensitive("TLS_CERT_PATH") {
		t.Error("expected TLS_CERT_PATH to be sensitive")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("expected API_KEY to NOT be sensitive with custom patterns")
	}
}

func TestRedact_MasksSensitiveValues(t *testing.T) {
	r := redactor.New(nil, "REDACTED")

	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}

	result := r.Redact(env)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT to be unchanged, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "REDACTED" {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "REDACTED" {
		t.Errorf("expected API_KEY to be masked, got %q", result["API_KEY"])
	}
}

func TestRedact_DefaultMask(t *testing.T) {
	r := redactor.New(nil, "")
	env := map[string]string{"DB_PASSWORD": "hunter2"}
	result := r.Redact(env)
	if result["DB_PASSWORD"] != "***" {
		t.Errorf("expected default mask ***, got %q", result["DB_PASSWORD"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	r := redactor.New(nil, "")
	env := map[string]string{"API_KEY": "original"}
	_ = r.Redact(env)
	if env["API_KEY"] != "original" {
		t.Error("Redact must not mutate the input map")
	}
}
