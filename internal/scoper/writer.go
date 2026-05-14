package scoper

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes the scoping result as human-readable text to w.
func WriteText(w io.Writer, res Result, env map[string]string) error {
	if w == nil {
		return fmt.Errorf("scoper: writer must not be nil")
	}

	for _, scope := range res.Scopes {
		fmt.Fprintf(w, "[%s]\n", scope.Name)
		for _, key := range scope.Keys {
			fmt.Fprintf(w, "  %s=%s\n", key, env[key])
		}
	}

	if len(res.Unscoped) > 0 {
		fmt.Fprintln(w, "[unscoped]")
		for _, key := range res.Unscoped {
			fmt.Fprintf(w, "  %s=%s\n", key, env[key])
		}
	}

	return nil
}

// jsonScope is the JSON representation of a single scope.
type jsonScope struct {
	Name   string            `json:"name"`
	Keys   []string          `json:"keys"`
	Values map[string]string `json:"values"`
}

// jsonResult is the top-level JSON output structure.
type jsonResult struct {
	Scopes   []jsonScope `json:"scopes"`
	Unscoped []string    `json:"unscoped"`
}

// WriteJSON writes the scoping result as JSON to w.
func WriteJSON(w io.Writer, res Result, env map[string]string) error {
	if w == nil {
		return fmt.Errorf("scoper: writer must not be nil")
	}

	out := jsonResult{Unscoped: res.Unscoped}
	for _, scope := range res.Scopes {
		vals := make(map[string]string, len(scope.Keys))
		for _, k := range scope.Keys {
			vals[k] = env[k]
		}
		out.Scopes = append(out.Scopes, jsonScope{
			Name:   scope.Name,
			Keys:   scope.Keys,
			Values: vals,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
