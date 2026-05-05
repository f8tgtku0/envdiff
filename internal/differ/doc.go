// Package differ provides functionality for comparing two sets of
// environment variables parsed from .env files.
//
// The primary entry point is the Compare function, which accepts two
// maps of key-value pairs and returns a Result describing any
// discrepancies found between them, including:
//
//   - Keys present in the left map but missing from the right
//   - Keys present in the right map but missing from the left
//   - Keys present in both maps but with differing values
//
// Example usage:
//
//	left  := map[string]string{"HOST": "localhost", "PORT": "8080"}
//	right := map[string]string{"HOST": "prod.example.com", "PORT": "8080"}
//	result := differ.Compare(left, right)
//	// result.Mismatched contains {"HOST": {Left: "localhost", Right: "prod.example.com"}}
package differ
