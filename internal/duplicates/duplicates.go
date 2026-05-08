// Package duplicates detects duplicate keys within a single .env map
// or across multiple environment maps.
package duplicates

// Entry records a key that appears more than once and in which sources.
type Entry struct {
	Key     string
	Sources []string
}

// Result holds the outcome of a duplicate-key detection run.
type Result struct {
	// CrossEnv contains keys that appear in more than one named environment.
	CrossEnv []Entry
}

// DetectCross scans a map of environment-name → key/value pairs and returns
// every key that is present in more than one environment.
//
// The envs parameter is keyed by an arbitrary environment label (e.g. "prod",
// "staging") whose value is the parsed key/value map for that file.
func DetectCross(envs map[string]map[string]string) Result {
	// Build an inverted index: key → list of env names that contain it.
	index := make(map[string][]string)

	// Iterate in a stable order so results are deterministic.
	names := sortedKeys(envs)
	for _, name := range names {
		for key := range envs[name] {
			index[key] = append(index[key], name)
		}
	}

	var cross []Entry
	for _, key := range sortedStringKeys(index) {
		sources := index[key]
		if len(sources) > 1 {
			cross = append(cross, Entry{Key: key, Sources: sources})
		}
	}

	return Result{CrossEnv: cross}
}

// sortedKeys returns the keys of a map[string]map[string]string in sorted order.
func sortedKeys(m map[string]map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

// sortedStringKeys returns the keys of a map[string][]string in sorted order.
func sortedStringKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
