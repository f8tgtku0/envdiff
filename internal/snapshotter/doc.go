// Package snapshotter provides utilities for capturing, saving, and loading
// point-in-time snapshots of environment variable maps.
//
// A Snapshot records the full key-value contents of an env map together with
// a label and a UTC timestamp. Snapshots are serialised to JSON files and can
// be reloaded later for comparison with the current state of an environment,
// enabling change-tracking and audit workflows.
//
// Basic usage:
//
//	env := map[string]string{"APP_ENV": "staging", "PORT": "3000"}
//	opts := snapshotter.DefaultOptions()
//	opts.Label = "pre-deploy"
//
//	snap, err := snapshotter.Take(env, opts)
//	if err != nil { ... }
//
//	if err := snapshotter.Save(snap, "snap.json"); err != nil { ... }
//
//	loaded, err := snapshotter.Load("snap.json")
package snapshotter
