// Package aliaser provides functionality for creating key aliases within an
// environment map.
//
// An alias copies the value of an existing key into a new key name, leaving
// the original key intact. This is useful when migrating from one naming
// convention to another while maintaining backward compatibility.
//
// Example usage:
//
//	aliases := aliaser.AliasMap{
//		"DB_HOST": "DATABASE_HOST",
//		"DB_PORT": "DATABASE_PORT",
//	}
//	out, result, err := aliaser.Alias(env, aliases, aliaser.DefaultOptions())
package aliaser
