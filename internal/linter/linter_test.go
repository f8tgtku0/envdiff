package linter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/linter"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	issues := linter.Lint(env, linter.DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"database_url": "postgres://localhost/db",
	}
	issues := linter.Lint(env, linter.DefaultOptions())
	if !containsMessage(issues, "not uppercase") {
		t.Fatalf("expected uppercase issue, got %v", issues)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
	}
	issues := linter.Lint(env, linter.DefaultOptions())
	if !containsMessage(issues, "empty") {
		t.Fatalf("expected empty-value issue, got %v", issues)
	}
}

func TestLint_KeyWithSpace(t *testing.T) {
	env := map[string]string{
		"MY KEY": "value",
	}
	opts := linter.Options{DisallowSpaces: true}
	issues := linter.Lint(env, opts)
	if !containsMessage(issues, "whitespace") {
		t.Fatalf("expected whitespace issue, got %v", issues)
	}
}

func TestLint_DisabledRules(t *testing.T) {
	env := map[string]string{
		"lower_key": "",
	}
	opts := linter.Options{
		RequireUppercase: false,
		DisallowEmpty:    false,
		DisallowSpaces:   false,
	}
	issues := linter.Lint(env, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with all rules disabled, got %v", issues)
	}
}

func TestIssue_String(t *testing.T) {
	i := linter.Issue{Key: "FOO", Message: "value is empty"}
	s := i.String()
	if !strings.Contains(s, "FOO") || !strings.Contains(s, "value is empty") {
		t.Fatalf("unexpected Issue.String() output: %q", s)
	}
}

// containsMessage returns true if any issue message contains the given substring.
func containsMessage(issues []linter.Issue, substr string) bool {
	for _, i := range issues {
		if strings.Contains(i.Message, substr) {
			return true
		}
	}
	return false
}
