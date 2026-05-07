package resolver_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/resolver"
)

func TestResolve_OverwriteTrue(t *testing.T) {
	sources := []resolver.Source{
		{Label: "base", Env: map[string]string{"APP_ENV": "development", "DB_HOST": "localhost"}},
		{Label: "override", Env: map[string]string{"APP_ENV": "production"}},
	}
	env, prov, err := resolver.Resolve(sources, resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected production, got %q", env["APP_ENV"])
	}
	if prov["APP_ENV"] != "override" {
		t.Errorf("expected provenance override, got %q", prov["APP_ENV"])
	}
	if env["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", env["DB_HOST"])
	}
}

func TestResolve_OverwriteFalse(t *testing.T) {
	sources := []resolver.Source{
		{Label: "first", Env: map[string]string{"KEY": "first-value"}},
		{Label: "second", Env: map[string]string{"KEY": "second-value"}},
	}
	opts := resolver.Options{Overwrite: false, Strict: false}
	env, prov, err := resolver.Resolve(sources, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "first-value" {
		t.Errorf("expected first-value, got %q", env["KEY"])
	}
	if prov["KEY"] != "first" {
		t.Errorf("expected provenance first, got %q", prov["KEY"])
	}
}

func TestResolve_StrictConflict(t *testing.T) {
	sources := []resolver.Source{
		{Label: "a", Env: map[string]string{"SECRET": "abc"}},
		{Label: "b", Env: map[string]string{"SECRET": "xyz"}},
	}
	opts := resolver.Options{Overwrite: true, Strict: true}
	_, _, err := resolver.Resolve(sources, opts)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestResolve_StrictSameValue(t *testing.T) {
	sources := []resolver.Source{
		{Label: "a", Env: map[string]string{"PORT": "8080"}},
		{Label: "b", Env: map[string]string{"PORT": "8080"}},
	}
	opts := resolver.Options{Overwrite: true, Strict: true}
	env, _, err := resolver.Resolve(sources, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", env["PORT"])
	}
}

func TestResolve_EmptySources(t *testing.T) {
	env, prov, err := resolver.Resolve(nil, resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 0 || len(prov) != 0 {
		t.Errorf("expected empty maps, got env=%v prov=%v", env, prov)
	}
}
