package merger

import "fmt"

// Strategy controls how conflicting keys are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first source that defines a key.
	StrategyFirst Strategy = iota

	// StrategyLast keeps the value from the last source that defines a key.
	StrategyLast

	// StrategyError returns an error when the same key has differing values
	// across sources.
	StrategyError
)

// ParseStrategy converts a string name into a Strategy constant.
// Accepted values (case-insensitive): "first", "last", "error".
func ParseStrategy(s string) (Strategy, error) {
	switch s {
	case "first":
		return StrategyFirst, nil
	case "last":
		return StrategyLast, nil
	case "error":
		return StrategyError, nil
	default:
		return 0, fmt.Errorf("unknown merge strategy %q: must be one of first, last, error", s)
	}
}

// ValidStrategies lists the accepted strategy name strings.
var ValidStrategies = []string{"first", "last", "error"}
