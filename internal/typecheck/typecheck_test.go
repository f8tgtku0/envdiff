package typecheck

import (
	"bytes"
	"strings"
	"testing"
)

func TestCheck_NoIssues(t *testing.T) {
	env := map[string]string{
		"SERVER_PORT":    "8080",
		"DATABASE_URL":   "http://localhost:5432",
		"FEATURE_ENABLED": "true",
	}
	issues, err := Check(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestCheck_InvalidPort(t *testing.T) {
	env := map[string]string{"APP_PORT": "not-a-number"}
	issues, err := Check(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Expected != TypeInt {
		t.Errorf("expected TypeInt, got %s", issues[0].Expected)
	}
}

func TestCheck_InvalidBool(t *testing.T) {
	env := map[string]string{"DARK_MODE_ENABLED": "maybe"}
	issues, err := Check(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DARK_MODE_ENABLED" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
}

func TestCheck_InvalidURL(t *testing.T) {
	env := map[string]string{"API_URL": "ftp://bad"}
	issues, err := Check(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestCheck_InvalidEmail(t *testing.T) {
	env := map[string]string{"ADMIN_EMAIL": "not-an-email"}
	opts := Options{Rules: []Rule{{KeyPattern: "_EMAIL$", Type: TypeEmail}}}
	issues, err := Check(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestCheck_InvalidPattern(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	opts := Options{Rules: []Rule{{KeyPattern: "[", Type: TypeString}}}
	_, err := Check(env, opts)
	if err == nil {
		t.Error("expected error for invalid regexp, got nil")
	}
}

func TestCheck_CustomRule_Float(t *testing.T) {
	env := map[string]string{"THRESHOLD_RATIO": "abc"}
	opts := Options{Rules: []Rule{{KeyPattern: "_RATIO$", Type: TypeFloat}}}
	issues, err := Check(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestWriteText_Clean(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteText(nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "all values match") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWriteText_WithIssues(t *testing.T) {
	issues := []Issue{
		{Key: "APP_PORT", Value: "abc", Expected: TypeInt, Reason: `"abc" is not a valid integer`},
	}
	var buf bytes.Buffer
	if err := WriteText(issues, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "APP_PORT") {
		t.Errorf("expected APP_PORT in output, got: %q", buf.String())
	}
}

func TestWriteJSON_WithIssues(t *testing.T) {
	issues := []Issue{
		{Key: "DB_URL", Value: "ftp://x", Expected: TypeURL, Reason: `"ftp://x" is not a valid URL`},
	}
	var buf bytes.Buffer
	if err := WriteJSON(issues, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB_URL") {
		t.Errorf("expected DB_URL in JSON output, got: %q", buf.String())
	}
}
