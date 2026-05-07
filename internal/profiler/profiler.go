package profiler

import (
	"fmt"
	"strings"
)

// KeyProfile holds analysis results for a single key across multiple environments.
type KeyProfile struct {
	Key         string
	Environments []string
	Values      map[string]string // env name -> value
	Unique      bool             // true if all present values are identical
	MissingIn   []string         // environments where the key is absent
}

// Profile analyses a set of named environments and returns per-key profiles.
// envs maps an environment label (e.g. "staging") to its parsed key/value map.
func Profile(envs map[string]map[string]string) []KeyProfile {
	envNames := sortedKeys(envs)

	// Collect all unique keys across every environment.
	keySet := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	profiles := make([]KeyProfile, 0, len(keySet))
	for key := range keySet {
		profile := KeyProfile{
			Key:          key,
			Environments: envNames,
			Values:       make(map[string]string),
		}

		seenValues := map[string]struct{}{}
		for _, env := range envNames {
			val, ok := envs[env][key]
			if ok {
				profile.Values[env] = val
				seenValues[val] = struct{}{}
			} else {
				profile.MissingIn = append(profile.MissingIn, env)
			}
		}

		profile.Unique = len(seenValues) == 1
		profiles = append(profiles, profile)
	}

	// Sort profiles by key name for deterministic output.
	sortProfiles(profiles)
	return profiles
}

// Summary returns a human-readable summary line for a KeyProfile.
func Summary(p KeyProfile) string {
	var parts []string
	if len(p.MissingIn) > 0 {
		parts = append(parts, fmt.Sprintf("missing in [%s]", strings.Join(p.MissingIn, ", ")))
	}
	if !p.Unique && len(p.MissingIn) < len(p.Environments) {
		parts = append(parts, "values differ across environments")
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s: consistent across all environments", p.Key)
	}
	return fmt.Sprintf("%s: %s", p.Key, strings.Join(parts, "; "))
}

func sortedKeys(m map[string]map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}
