package digester_test

import (
	"testing"

	"github.com/user/envdiff/internal/digester"
)

func base() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"LOG_LEVEL": "info",
	}
}

func TestDigest_DeterministicOutput(t *testing.T) {
	opts := digester.DefaultOptions()
	r1, err := digester.Digest(base(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := digester.Digest(base(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r1.Hex != r2.Hex {
		t.Errorf("expected identical digests, got %q vs %q", r1.Hex, r2.Hex)
	}
}

func TestDigest_KeyCountReflected(t *testing.T) {
	r, err := digester.Digest(base(), digester.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.KeyCount != 4 {
		t.Errorf("expected KeyCount=4, got %d", r.KeyCount)
	}
}

func TestDigest_ChangedValueProducesDifferentHash(t *testing.T) {
	opts := digester.DefaultOptions()
	r1, _ := digester.Digest(base(), opts)

	modified := base()
	modified["APP_ENV"] = "staging"
	r2, _ := digester.Digest(modified, opts)

	if digester.Equal(r1, r2) {
		t.Error("expected different digests after value change")
	}
}

func TestDigest_RenameDetectedWhenIncludeKeysTrue(t *testing.T) {
	opts := digester.DefaultOptions() // IncludeKeys: true
	r1, _ := digester.Digest(map[string]string{"FOO": "bar"}, opts)
	r2, _ := digester.Digest(map[string]string{"BAZ": "bar"}, opts)

	if digester.Equal(r1, r2) {
		t.Error("expected different digests when key name changes and IncludeKeys=true")
	}
}

func TestDigest_RenameNotDetectedWhenIncludeKeysFalse(t *testing.T) {
	opts := digester.Options{Algorithm: digester.SHA256, IncludeKeys: false}
	r1, _ := digester.Digest(map[string]string{"FOO": "bar"}, opts)
	r2, _ := digester.Digest(map[string]string{"BAZ": "bar"}, opts)

	if !digester.Equal(r1, r2) {
		t.Error("expected same digest when only key name changes and IncludeKeys=false")
	}
}

func TestDigest_NilEnvReturnsError(t *testing.T) {
	_, err := digester.Digest(nil, digester.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestEqual_SameResult(t *testing.T) {
	r, _ := digester.Digest(base(), digester.DefaultOptions())
	if !digester.Equal(r, r) {
		t.Error("expected Equal to return true for identical result")
	}
}
