// Package promoter provides functionality for promoting environment variables
// from one environment (e.g. staging) into another (e.g. production).
//
// Use Promote to copy keys from a source map into a destination map.
// Promotion can be scoped by key prefix or an explicit allow-list, and
// existing keys in the destination can be preserved or overwritten.
//
// Example:
//
//	src := map[string]string{"DB_HOST": "staging-db", "APP_ENV": "staging"}
//	dst := map[string]string{"APP_ENV": "production"}
//
//	opts := promoter.Options{Prefix: "DB_", Overwrite: false}
//	res, err := promoter.Promote(src, dst, opts)
//	// dst now contains DB_HOST="staging-db"; APP_ENV is unchanged.
package promoter
