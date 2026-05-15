// Package digester produces stable, deterministic content hashes for
// env maps (map[string]string).
//
// A digest uniquely fingerprints an environment snapshot. Two digests
// can be compared with Equal to decide whether anything has changed
// since the last snapshot without needing to inspect individual keys.
//
// Usage:
//
//	opts := digester.DefaultOptions()
//	before, _ := digester.Digest(envA, opts)
//	after,  _ := digester.Digest(envB, opts)
//	if !digester.Equal(before, after) {
//		// environment has changed
//	}
package digester
