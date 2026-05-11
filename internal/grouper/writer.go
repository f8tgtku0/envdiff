package grouper

import (
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable grouped summary to w.
func WriteText(res Result, w io.Writer) error {
	for _, g := range res.Groups {
		fmt.Fprintf(w, "[%s] (%d keys)\n", g.Prefix, len(g.Keys))
		for _, k := range g.Keys {
			fmt.Fprintf(w, "  %s=%s\n", k, g.Env[k])
		}
	}

	if len(res.Ungrouped) > 0 {
		fmt.Fprintf(w, "[ungrouped] (%d keys)\n", len(res.Ungrouped))
		for _, k := range res.Ungrouped {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
	return nil
}

// WriteMarkdown writes a Markdown table per group to w.
func WriteMarkdown(res Result, w io.Writer) error {
	writeTable := func(prefix string, keys []string, env map[string]string) {
		fmt.Fprintf(w, "## %s\n\n", prefix)
		fmt.Fprintln(w, "| Key | Value |")
		fmt.Fprintln(w, "|-----|-------|")
		for _, k := range keys {
			v := strings.ReplaceAll(env[k], "|", "\\|")
			fmt.Fprintf(w, "| `%s` | `%s` |\n", k, v)
		}
		fmt.Fprintln(w, "")
	}

	for _, g := range res.Groups {
		writeTable(g.Prefix, g.Keys, g.Env)
	}

	if len(res.Ungrouped) > 0 {
		fmt.Fprintf(w, "## ungrouped\n\n")
		fmt.Fprintln(w, "| Key |")
		fmt.Fprintln(w, "|-----|")
		for _, k := range res.Ungrouped {
			fmt.Fprintf(w, "| `%s` |\n", k)
		}
		fmt.Fprintln(w, "")
	}
	return nil
}
