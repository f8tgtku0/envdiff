// Package comparator provides value-level comparison of environment variables
// across multiple named environments.
//
// Unlike the differ package — which works on exactly two environments — the
// comparator supports an arbitrary number of environments identified by a
// string label (e.g. "dev", "staging", "prod").
//
// Basic usage:
//
//	envs := map[string]map[string]string{
//		"dev":     {"PORT": "3000", "DEBUG": "true"},
//		"staging": {"PORT": "3000", "DEBUG": "false"},
//		"prod":    {"PORT": "8080"},
//	}
//
//	results, err := comparator.Compare(envs, comparator.DefaultOptions())
//
// Results can be rendered as plain text or JSON via WriteText / WriteJSON.
package comparator
