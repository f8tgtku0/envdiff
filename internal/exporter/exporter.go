package exporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents an export output format.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
)

// Options configures the export behaviour.
type Options struct {
	Format  Format
	OutFile string // empty means stdout
}

// Write serialises env vars in the requested format to the configured output.
func Write(vars map[string]string, opts Options) error {
	w := io.Writer(os.Stdout)
	if opts.OutFile != "" {
		f, err := os.Create(opts.OutFile)
		if err != nil {
			return fmt.Errorf("exporter: create file: %w", err)
		}
		defer f.Close()
		w = f
	}

	keys := sortedKeys(vars)

	for _, k := range keys {
		v := vars[k]
		var line string
		switch opts.Format {
		case FormatExport:
			line = fmt.Sprintf("export %s=%s", k, quoteValue(v))
		case FormatShell:
			line = fmt.Sprintf("%s=%s", k, quoteValue(v))
		default: // dotenv
			line = fmt.Sprintf("%s=%s", k, quoteValue(v))
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("exporter: write: %w", err)
		}
	}
	return nil
}

func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n#$") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
