package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env map.
type Snapshot struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// Options controls Snapshot behaviour.
type Options struct {
	// Label is a human-readable name for the snapshot.
	Label string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Label: "snapshot",
	}
}

// Take creates a Snapshot from the provided env map.
// The env map is shallow-copied so later mutations do not affect the snapshot.
func Take(env map[string]string, opts Options) (*Snapshot, error) {
	if env == nil {
		return nil, fmt.Errorf("snapshotter: env map must not be nil")
	}
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		Label:     opts.Label,
		Timestamp: time.Now().UTC(),
		Env:       copy,
	}, nil
}

// Save writes a Snapshot to a JSON file at the given path.
func Save(s *Snapshot, path string) error {
	if s == nil {
		return fmt.Errorf("snapshotter: snapshot must not be nil")
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshotter: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshotter: encode: %w", err)
	}
	return nil
}

// Load reads a Snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshotter: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshotter: decode: %w", err)
	}
	return &s, nil
}
