// Package encoder provides utilities for encoding env maps into
// various string representations suitable for shell injection,
// Docker --env-file format, or inline export statements.
package encoder

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the target encoding format.
type Format string

const (
	FormatInline Format = "inline" // KEY=VALUE KEY2=VALUE2 (single line)
	FormatExport Format = "export" // export KEY=VALUE (one per line)
	FormatDocker Format = "docker" // --env KEY=VALUE flags
	FormatShell  Format = "shell"  // KEY=VALUE (one per line, no export)
)

// Options controls encoder behaviour.
type Options struct {
	Format    Format
	Quote     bool // wrap values in double quotes
	SortKeys  bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format:   FormatShell,
		Quote:    false,
		SortKeys: true,
	}
}

// Encode writes the env map to w using the configured format.
func Encode(env map[string]string, w io.Writer, opts Options) error {
	if env == nil {
		return fmt.Errorf("encoder: env map must not be nil")
	}
	if w == nil {
		return fmt.Errorf("encoder: writer must not be nil")
	}

	keys := sortedKeys(env, opts.SortKeys)

	switch opts.Format {
	case FormatInline:
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, k+"="+maybeQuote(env[k], opts.Quote))
		}
		_, err := fmt.Fprintln(w, strings.Join(parts, " "))
		return err

	case FormatExport:
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "export %s=%s\n", k, maybeQuote(env[k], opts.Quote)); err != nil {
				return err
			}
		}
		return nil

	case FormatDocker:
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "--env %s=%s\n", k, maybeQuote(env[k], opts.Quote)); err != nil {
				return err
			}
		}
		return nil

	default: // FormatShell
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "%s=%s\n", k, maybeQuote(env[k], opts.Quote)); err != nil {
				return err
			}
		}
		return nil
	}
}

func maybeQuote(v string, quote bool) string {
	if !quote {
		return v
	}
	return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
}

func sortedKeys(env map[string]string, doSort bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if doSort {
		sort.Strings(keys)
	}
	return keys
}
