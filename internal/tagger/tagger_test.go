package tagger

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":        "localhost",
		"DB_PASSWORD":    "secret",
		"API_TOKEN":      "tok_abc",
		"FEATURE_DARK":   "true",
		"APP_NAME":       "envdiff",
		"LOG_LEVEL":      "info",
	}
}

func TestTag_DefaultRules(t *testing.T) {
	res, err := Tag(baseEnv(), DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Tags["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD to be tagged")
	}
	if _, ok := res.Tags["API_TOKEN"]; !ok {
		t.Error("expected API_TOKEN to be tagged")
	}
	if _, ok := res.Tags["FEATURE_DARK"]; !ok {
		t.Error("expected FEATURE_DARK to be tagged")
	}
}

func TestTag_UntaggedKeys(t *testing.T) {
	res, err := Tag(baseEnv(), DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	untagged := map[string]bool{}
	for _, k := range res.Untagged {
		untagged[k] = true
	}
	if !untagged["APP_NAME"] {
		t.Error("expected APP_NAME to be untagged")
	}
	if !untagged["LOG_LEVEL"] {
		t.Error("expected LOG_LEVEL to be untagged")
	}
}

func TestTag_AllowMultiple(t *testing.T) {
	env := map[string]string{
		"DB_SECRET_TOKEN": "x",
	}
	opts := DefaultOptions()
	opts.AllowMultiple = true
	res, err := Tag(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := res.Tags["DB_SECRET_TOKEN"]
	if len(tags) < 2 {
		t.Errorf("expected multiple tags, got %v", tags)
	}
}

func TestTag_InvalidPattern(t *testing.T) {
	opts := Options{
		Rules: map[string]string{
			"bad": `[invalid`,
		},
	}
	_, err := Tag(baseEnv(), opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestTag_EmptyEnv(t *testing.T) {
	res, err := Tag(map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Tags) != 0 {
		t.Errorf("expected no tags, got %d", len(res.Tags))
	}
	if len(res.Untagged) != 0 {
		t.Errorf("expected no untagged keys, got %d", len(res.Untagged))
	}
}

func TestTag_CustomRules(t *testing.T) {
	env := map[string]string{
		"STRIPE_KEY": "sk_live_abc",
		"REGION":     "us-east-1",
	}
	opts := Options{
		Rules: map[string]string{
			"payment": `(?i)stripe`,
		},
	}
	res, err := Tag(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Tags["STRIPE_KEY"]; !ok {
		t.Error("expected STRIPE_KEY to be tagged as payment")
	}
	if len(res.Untagged) != 1 || res.Untagged[0] != "REGION" {
		t.Errorf("expected REGION to be untagged, got %v", res.Untagged)
	}
}
