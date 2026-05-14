// Package scoper partitions an environment map into named scopes based on
// key-prefix rules. Each scope groups keys that share a common prefix,
// making it easier to reason about large .env files that contain variables
// for multiple services or subsystems.
//
// Example usage:
//
//	env := map[string]string{
//		"DB_HOST": "localhost",
//		"DB_PORT": "5432",
//		"REDIS_URL": "redis://localhost",
//		"APP_NAME": "envdiff",
//	}
//
//	opts := scoper.DefaultOptions()
//	opts.Rules = map[string][]string{
//		"database": {"DB_"},
//		"cache":    {"REDIS_"},
//	}
//
//	res, _ := scoper.Scope(env, opts)
package scoper
