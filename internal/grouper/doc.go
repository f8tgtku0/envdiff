// Package grouper partitions an env map into named groups based on key
// prefixes.
//
// Keys are split on a configurable delimiter (default "_") and the leading
// segment(s) are used as the group name. For example, "DB_HOST" and
// "DB_PORT" both belong to the "DB" group.
//
// Usage:
//
//	env := map[string]string{
//		"DB_HOST": "localhost",
//		"APP_PORT": "8080",
//	}
//	res := grouper.GroupEnv(env, grouper.DefaultOptions())
//	grouper.WriteText(res, os.Stdout)
//
// The Result contains a sorted slice of Group values and a separate
// Ungrouped slice for keys that have no delimiter.
package grouper
