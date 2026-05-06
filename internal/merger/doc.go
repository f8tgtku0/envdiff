// Package merger provides functionality for merging multiple .env file
// maps into a single unified map.
//
// When keys appear in more than one source, the conflict resolution
// strategy determines which value wins:
//
//   - StrategyFirst – keep the value from the earliest source that defines the key.
//   - StrategyLast  – keep the value from the latest source that defines the key.
//   - StrategyError – return an error if the same key has different values across sources.
//
// Example usage:
//
//	result, err := merger.Merge(
//		[]map[string]string{base, override},
//		merger.StrategyLast,
//	)
package merger
