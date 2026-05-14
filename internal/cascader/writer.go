package cascader

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable summary of the cascade result to w.
// Each line shows: KEY=value  (from: layerName)
func WriteText(r *Result, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("cascader: writer must not be nil")
	}
	keys := sortedKeys(r.Env)
	for _, k := range keys {
		origin := r.Origins[k]
		if _, err := fmt.Fprintf(w, "%s=%s  (from: %s)\n", k, r.Env[k], origin); err != nil {
			return err
		}
	}
	return nil
}

type jsonEntry struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Origin string `json:"origin"`
}

// WriteJSON writes the cascade result as a JSON array to w.
func WriteJSON(r *Result, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("cascader: writer must not be nil")
	}
	keys := sortedKeys(r.Env)
	entries := make([]jsonEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, jsonEntry{
			Key:    k,
			Value:  r.Env[k],
			Origin: r.Origins[k],
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
