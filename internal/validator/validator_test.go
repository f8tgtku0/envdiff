package validator_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/validator"
)

func TestValidate_NoRules(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": ""}
	violations := validator.Validate(env, nil)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "", "API_KEY": "secret"}
	rules := []validator.Rule{
		{Required: true},
	}
	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL violation, got %q", violations[0].Key)
	}
}

func TestValidate_RequiredWithKeyPattern(t *testing.T) {
	env := map[string]string{"DB_HOST": "", "DB_PORT": "5432", "APP_NAME": ""}
	rules := []validator.Rule{
		{KeyPattern: `^DB_`, Required: true},
	}
	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if violations[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST violation, got %q", violations[0].Key)
	}
}

func TestValidate_ValuePattern(t *testing.T) {
	env := map[string]string{"PORT": "abc", "TIMEOUT": "30"}
	rules := []validator.Rule{
		{KeyPattern: `^PORT$`, ValuePattern: `^\d+$`},
	}
	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected PORT violation, got %q", violations[0].Key)
	}
}

func TestValidate_InvalidKeyPattern(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	rules := []validator.Rule{
		{KeyPattern: `[invalid`},
	}
	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for bad pattern, got %d", len(violations))
	}
	if violations[0].Key != "<rule>" {
		t.Errorf("expected <rule> key, got %q", violations[0].Key)
	}
}

func TestValidate_MultipleRules(t *testing.T) {
	env := map[string]string{"PORT": "notanumber", "SECRET": ""}
	rules := []validator.Rule{
		{KeyPattern: `^PORT$`, ValuePattern: `^\d+$`},
		{KeyPattern: `^SECRET$`, Required: true},
	}
	violations := validator.Validate(env, rules)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(violations), violations)
	}
}
