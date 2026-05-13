package annotator

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls annotation behaviour.
type Options struct {
	// Prefix is prepended to each annotation comment, e.g. "#".
	CommentPrefix string
	// Overwrite replaces an existing annotation for a key.
	Overwrite bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		CommentPrefix: "#",
		Overwrite:     false,
	}
}

// Annotation holds the comment text attached to a key.
type Annotation struct {
	Key     string
	Comment string
}

// Result is the annotated environment map plus metadata.
type Result struct {
	// Env is the original env map (unmodified).
	Env map[string]string
	// Annotations maps each key to its annotation.
	Annotations map[string]Annotation
	// Skipped lists keys whose annotation was not applied (Overwrite=false and key already had one).
	Skipped []string
}

// Annotate attaches comments from notes to matching keys in env.
// notes maps key names to free-form comment text.
func Annotate(env map[string]string, notes map[string]string, opts Options) (*Result, error) {
	if env == nil {
		return nil, fmt.Errorf("annotator: env must not be nil")
	}

	prefix := opts.CommentPrefix
	if prefix == "" {
		prefix = "#"
	}

	existing := make(map[string]Annotation, len(env))
	var skipped []string

	for key := range env {
		note, ok := notes[key]
		if !ok {
			continue
		}
		note = strings.TrimSpace(note)
		if note == "" {
			continue
		}

		if _, already := existing[key]; already && !opts.Overwrite {
			skipped = append(skipped, key)
			continue
		}

		existing[key] = Annotation{
			Key:     key,
			Comment: fmt.Sprintf("%s %s", prefix, note),
		}
	}

	sort.Strings(skipped)

	envCopy := make(map[string]string, len(env))
	for k, v := range env {
		envCopy[k] = v
	}

	return &Result{
		Env:         envCopy,
		Annotations: existing,
		Skipped:     skipped,
	}, nil
}
