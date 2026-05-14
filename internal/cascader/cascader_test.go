package cascader

import (
	"testing"
)

func layers(pairs ...any) []Layer {
	var out []Layer
	for i := 0; i+1 < len(pairs); i += 2 {
		name := pairs[i].(string)
		env := pairs[i+1].(map[string]string)
		out = append(out, Layer{Name: name, Env: env})
	}
	return out
}

func TestCascade_FirstWins(t *testing.T) {
	ls := layers(
		"base", map[string]string{"KEY": "base", "ONLY_BASE": "yes"},
		"override", map[string]string{"KEY": "override", "ONLY_OVERRIDE": "yes"},
	)
	res, err := Cascade(ls, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "base" {
		t.Errorf("expected base to win, got %q", res.Env["KEY"])
	}
	if res.Origins["KEY"] != "base" {
		t.Errorf("expected origin 'base', got %q", res.Origins["KEY"])
	}
	if res.Env["ONLY_OVERRIDE"] != "yes" {
		t.Errorf("expected ONLY_OVERRIDE to be present")
	}
}

func TestCascade_OverwriteTrue(t *testing.T) {
	ls := layers(
		"base", map[string]string{"KEY": "base"},
		"prod", map[string]string{"KEY": "prod"},
	)
	opts := Options{Overwrite: true}
	res, err := Cascade(ls, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "prod" {
		t.Errorf("expected prod to win, got %q", res.Env["KEY"])
	}
	if res.Origins["KEY"] != "prod" {
		t.Errorf("expected origin 'prod', got %q", res.Origins["KEY"])
	}
}

func TestCascade_NilLayerSkipped(t *testing.T) {
	ls := []Layer{
		{Name: "base", Env: map[string]string{"A": "1"}},
		{Name: "empty", Env: nil},
		{Name: "extra", Env: map[string]string{"B": "2"}},
	}
	res, err := Cascade(ls, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["A"] != "1" || res.Env["B"] != "2" {
		t.Errorf("unexpected env: %v", res.Env)
	}
}

func TestCascade_StrictNilReturnsError(t *testing.T) {
	ls := []Layer{
		{Name: "base", Env: map[string]string{"A": "1"}},
		{Name: "bad", Env: nil},
	}
	_, err := Cascade(ls, Options{Strict: true})
	if err == nil {
		t.Fatal("expected error for nil layer in strict mode")
	}
}

func TestCascade_EmptyLayers(t *testing.T) {
	res, err := Cascade(nil, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}
