package annotator

import (
	"bytes"
	"strings"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"SECRET_KEY":   "s3cr3t",
		"PORT":         "8080",
	}
}

func TestAnnotate_Basic(t *testing.T) {
	notes := map[string]string{
		"DATABASE_URL": "Primary DB connection.",
		"SECRET_KEY":   "Keep this safe.",
	}
	res, err := Annotate(baseEnv(), notes, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Annotations) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(res.Annotations))
	}
	ann := res.Annotations["DATABASE_URL"]
	if !strings.HasPrefix(ann.Comment, "#") {
		t.Errorf("comment should start with '#', got %q", ann.Comment)
	}
	if !strings.Contains(ann.Comment, "Primary DB connection.") {
		t.Errorf("comment missing note text: %q", ann.Comment)
	}
}

func TestAnnotate_SkipsEmptyNotes(t *testing.T) {
	notes := map[string]string{"PORT": "   "}
	res, err := Annotate(baseEnv(), notes, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Annotations["PORT"]; ok {
		t.Error("empty note should not produce an annotation")
	}
}

func TestAnnotate_NilEnvReturnsError(t *testing.T) {
	_, err := Annotate(nil, map[string]string{}, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestAnnotate_CustomPrefix(t *testing.T) {
	opts := DefaultOptions()
	opts.CommentPrefix = "//"
	notes := map[string]string{"PORT": "HTTP port"}
	res, err := Annotate(baseEnv(), notes, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ann := res.Annotations["PORT"]
	if !strings.HasPrefix(ann.Comment, "//") {
		t.Errorf("expected '//' prefix, got %q", ann.Comment)
	}
}

func TestWriteText_AnnotationsAboveKeys(t *testing.T) {
	notes := map[string]string{"PORT": "HTTP port number"}
	res, _ := Annotate(baseEnv(), notes, DefaultOptions())

	var buf bytes.Buffer
	if err := WriteText(res, &buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	idxComment := strings.Index(out, "# HTTP port number")
	idxKey := strings.Index(out, "PORT=8080")
	if idxComment < 0 {
		t.Error("comment not found in output")
	}
	if idxKey < 0 {
		t.Error("key=value not found in output")
	}
	if idxComment > idxKey {
		t.Error("comment should appear before its key")
	}
}

func TestWriteSummary_Counts(t *testing.T) {
	notes := map[string]string{
		"DATABASE_URL": "DB",
		"SECRET_KEY":   "SK",
	}
	res, _ := Annotate(baseEnv(), notes, DefaultOptions())

	var buf bytes.Buffer
	if err := WriteSummary(res, &buf); err != nil {
		t.Fatalf("WriteSummary error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "annotated: 2") {
		t.Errorf("unexpected summary: %q", out)
	}
}
