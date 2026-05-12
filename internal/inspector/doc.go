// Package inspector provides cross-environment key inspection.
//
// It accepts a map of named environments (each a map[string]string) and
// produces a structured report describing every key: which environments
// define it, what value each environment holds, and whether the values are
// consistent across all environments that carry the key.
//
// Basic usage:
//
//	envs := map[string]map[string]string{
//		"dev":  {"APP_NAME": "myapp", "DEBUG": "true"},
//		"prod": {"APP_NAME": "myapp"},
//	}
//	result := inspector.Inspect(envs, inspector.DefaultOptions())
//	for _, info := range result.Keys {
//		fmt.Println(info.Key, info.Consistent)
//	}
package inspector
