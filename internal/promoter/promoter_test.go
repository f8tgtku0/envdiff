package promoter_test

import (
	"testing"

	"github.com/user/envdiff/internal/promoter"
)

func base() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV":  "staging",
		"SECRET":  "abc123",
	}
}

func TestPromote_NoFilter(t *testing.T) {
	src := base()
	dst := map[string]string{}

	res, err := promoter.Promote(src, dst, promoter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 4 {
		t.Errorf("expected 4 promoted, got %d", len(res.Promoted))
	}
	if dst["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost in dst")
	}
}

func TestPromote_SkipsExistingByDefault(t *testing.T) {
	src := base()
	dst := map[string]string{"DB_HOST": "prod-host"}

	res, err := promoter.Promote(src, dst, promoter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["DB_HOST"] != "prod-host" {
		t.Errorf("expected DB_HOST to remain prod-host")
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in skipped, got %v", res.Skipped)
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := base()
	dst := map[string]string{"DB_HOST": "prod-host"}
	opts := promoter.Options{Overwrite: true}

	_, err := promoter.Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST overwritten to localhost, got %s", dst["DB_HOST"])
	}
}

func TestPromote_PrefixFilter(t *testing.T) {
	src := base()
	dst := map[string]string{}
	opts := promoter.Options{Prefix: "DB_"}

	res, err := promoter.Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d: %v", len(res.Promoted), res.Promoted)
	}
	if _, ok := dst["APP_ENV"]; ok {
		t.Error("APP_ENV should not have been promoted")
	}
}

func TestPromote_ExplicitKeys(t *testing.T) {
	src := base()
	dst := map[string]string{}
	opts := promoter.Options{Keys: []string{"SECRET", "APP_ENV"}}

	res, err := promoter.Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if dst["DB_HOST"] != "" {
		t.Error("DB_HOST should not have been promoted")
	}
}

func TestPromote_NilSrcReturnsError(t *testing.T) {
	_, err := promoter.Promote(nil, map[string]string{}, promoter.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil src")
	}
}

func TestPromote_NilDstReturnsError(t *testing.T) {
	_, err := promoter.Promote(map[string]string{}, nil, promoter.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil dst")
	}
}
