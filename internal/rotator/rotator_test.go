package rotator_test

import (
	"testing"

	"github.com/user/envdiff/internal/rotator"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "old-pass",
		"API_KEY":     "old-key",
		"APP_ENV":     "production",
	}
}

func TestRotate_Basic(t *testing.T) {
	env := baseEnv()
	replacements := map[string]string{
		"DB_PASSWORD": "new-pass",
		"API_KEY":     "new-key",
	}
	res, err := rotator.Rotate(env, replacements, rotator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated keys, got %d", len(res.Rotated))
	}
	if res.Rotated["DB_PASSWORD"] != "new-pass" {
		t.Errorf("expected DB_PASSWORD=new-pass, got %s", res.Rotated["DB_PASSWORD"])
	}
	if res.Archived["DB_PASSWORD_OLD"] != "old-pass" {
		t.Errorf("expected DB_PASSWORD_OLD=old-pass, got %s", res.Archived["DB_PASSWORD_OLD"])
	}
}

func TestRotate_SkipsMissingKey(t *testing.T) {
	env := baseEnv()
	replacements := map[string]string{
		"NONEXISTENT": "value",
	}
	res, err := rotator.Rotate(env, replacements, rotator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 0 {
		t.Errorf("expected 0 rotated keys, got %d", len(res.Rotated))
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "NONEXISTENT" {
		t.Errorf("expected NONEXISTENT in skipped, got %v", res.Skipped)
	}
}

func TestRotate_ConflictNoOverwrite(t *testing.T) {
	env := baseEnv()
	env["DB_PASSWORD_OLD"] = "ancient-pass" // archive key already exists
	replacements := map[string]string{
		"DB_PASSWORD": "new-pass",
	}
	opts := rotator.DefaultOptions()
	opts.Overwrite = false
	res, err := rotator.Rotate(env, replacements, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Rotated["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be skipped due to archive conflict")
	}
	if len(res.Skipped) == 0 {
		t.Error("expected at least one skipped key")
	}
}

func TestRotate_ConflictWithOverwrite(t *testing.T) {
	env := baseEnv()
	env["DB_PASSWORD_OLD"] = "ancient-pass"
	replacements := map[string]string{
		"DB_PASSWORD": "new-pass",
	}
	opts := rotator.DefaultOptions()
	opts.Overwrite = true
	res, err := rotator.Rotate(env, replacements, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rotated["DB_PASSWORD"] != "new-pass" {
		t.Errorf("expected rotation with overwrite, got %v", res.Rotated)
	}
}

func TestRotate_CustomSuffix(t *testing.T) {
	env := baseEnv()
	replacements := map[string]string{"API_KEY": "new-key"}
	opts := rotator.Options{Suffix: "_PREV", Overwrite: false}
	res, err := rotator.Rotate(env, replacements, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Archived["API_KEY_PREV"]; !ok {
		t.Errorf("expected archive key API_KEY_PREV, got %v", res.Archived)
	}
}

func TestRotate_NilEnvReturnsError(t *testing.T) {
	_, err := rotator.Rotate(nil, map[string]string{}, rotator.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestRotate_NilReplacementsReturnsError(t *testing.T) {
	_, err := rotator.Rotate(baseEnv(), nil, rotator.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil replacements")
	}
}

func TestApply_MergesResult(t *testing.T) {
	env := baseEnv()
	replacements := map[string]string{"DB_PASSWORD": "new-pass"}
	res, _ := rotator.Rotate(env, replacements, rotator.DefaultOptions())
	out := rotator.Apply(env, res, "")
	if out["DB_PASSWORD"] != "new-pass" {
		t.Errorf("expected new-pass, got %s", out["DB_PASSWORD"])
	}
	if out["DB_PASSWORD_OLD"] != "old-pass" {
		t.Errorf("expected old-pass in archive, got %s", out["DB_PASSWORD_OLD"])
	}
	if out["APP_ENV"] != "production" {
		t.Error("expected unrelated key to be preserved")
	}
}
