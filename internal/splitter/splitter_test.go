package splitter_test

import (
	"testing"

	"github.com/user/envdiff/internal/splitter"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
		"APP_NAME":  "myapp",
		"APP_DEBUG": "false",
		"LOG_LEVEL": "info",
		"STANDALONE": "yes",
	}
}

func TestSplit_NoPrefixes(t *testing.T) {
	env := baseEnv()
	opts := splitter.DefaultOptions()
	r, err := splitter.Split(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Buckets[""]) != len(env) {
		t.Errorf("expected all %d keys in unmatched bucket, got %d", len(env), len(r.Buckets[""]))
	}
}

func TestSplit_BasicPrefixes(t *testing.T) {
	env := baseEnv()
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB", "APP"}
	r, err := splitter.Split(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Buckets["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(r.Buckets["DB"]))
	}
	if len(r.Buckets["APP"]) != 2 {
		t.Errorf("expected 2 APP keys, got %d", len(r.Buckets["APP"]))
	}
	if len(r.Buckets[""]) != 2 { // LOG_LEVEL + STANDALONE
		t.Errorf("expected 2 unmatched keys, got %d", len(r.Buckets[""]))
	}
}

func TestSplit_KeepPrefix(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB"}
	opts.KeepPrefix = true
	r, err := splitter.Split(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Buckets["DB"]["DB_HOST"]; !ok {
		t.Errorf("expected full key DB_HOST to be retained")
	}
}

func TestSplit_StripPrefix(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB"}
	opts.KeepPrefix = false
	r, err := splitter.Split(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Buckets["DB"]["HOST"]; !ok {
		t.Errorf("expected stripped key HOST in DB bucket")
	}
}

func TestSplit_NilEnvReturnsError(t *testing.T) {
	_, err := splitter.Split(nil, splitter.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestSplit_EmptyBucketCreated(t *testing.T) {
	env := map[string]string{"STANDALONE": "yes"}
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB"}
	r, err := splitter.Split(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Buckets["DB"]; !ok {
		t.Error("expected empty DB bucket to be pre-created")
	}
}
