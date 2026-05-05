package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestCompare_NoDiff(t *testing.T) {
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := differ.Compare(left, right)

	if result.HasDiff() {
		t.Errorf("expected no diff, got %+v", result)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"FOO": "bar", "ONLY_LEFT": "yes"}
	right := map[string]string{"FOO": "bar"}

	result := differ.Compare(left, right)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("expected ONLY_LEFT missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no keys missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "bar", "ONLY_RIGHT": "yes"}

	result := differ.Compare(left, right)

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("expected ONLY_RIGHT missing in left, got %v", result.MissingInLeft)
	}
	if len(result.MissingInRight) != 0 {
		t.Errorf("expected no keys missing in right, got %v", result.MissingInRight)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"FOO": "bar", "PORT": "8080"}
	right := map[string]string{"FOO": "bar", "PORT": "9090"}

	result := differ.Compare(left, right)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "PORT" || m.LeftValue != "8080" || m.RightValue != "9090" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := differ.Compare(map[string]string{}, map[string]string{})
	if result.HasDiff() {
		t.Errorf("expected no diff for empty maps, got %+v", result)
	}
}
