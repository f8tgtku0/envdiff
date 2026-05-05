// Package parser provides functionality for reading and parsing .env files
// into key-value maps for use by the envdiff comparison engine.
//
// Supported .env syntax:
//   - KEY=VALUE pairs (one per line)
//   - Lines beginning with '#' are treated as comments and skipped
//   - Empty lines are ignored
//   - Values may be wrapped in single or double quotes, which are stripped
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(env["DATABASE_URL"])
package parser
