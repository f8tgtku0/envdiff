// Package converter transforms env maps between different serialisation
// formats such as dotenv, JSON, YAML, and TOML-style key=value.
package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents a supported output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
	FormatExport Format = "export"
)

// ValidFormats lists all accepted format strings.
var ValidFormats = []Format{FormatDotenv, FormatJSON, FormatYAML, FormatExport}

// ParseFormat returns a Format from a raw string, or an error if unrecognised.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatDotenv:
		return FormatDotenv, nil
	case FormatJSON:
		return FormatJSON, nil
	case FormatYAML:
		return FormatYAML, nil
	case FormatExport:
		return FormatExport, nil
	default:
		return "", fmt.Errorf("converter: unknown format %q; valid: %v", s, ValidFormats)
	}
}

// Convert writes the env map to w in the requested format.
func Convert(env map[string]string, format Format, w io.Writer) error {
	keys := sortedKeys(env)

	switch format {
	case FormatDotenv:
		return writeDotenv(keys, env, w, false)
	case FormatExport:
		return writeDotenv(keys, env, w, true)
	case FormatJSON:
		return writeJSON(keys, env, w)
	case FormatYAML:
		return writeYAML(keys, env, w)
	default:
		return fmt.Errorf("converter: unsupported format %q", format)
	}
}

func writeDotenv(keys []string, env map[string]string, w io.Writer, export bool) error {
	prefix := ""
	if export {
		prefix = "export "
	}
	for _, k := range keys {
		v := env[k]
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(w, "%s%s=%s\n", prefix, k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(keys []string, env map[string]string, w io.Writer) error {
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = env[k]
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}

func writeYAML(keys []string, env map[string]string, w io.Writer) error {
	for _, k := range keys {
		v := env[k]
		if strings.ContainsAny(v, ":\n#") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(w, "%s: %s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
