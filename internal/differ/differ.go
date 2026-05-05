// Package differ compares parsed .env maps and reports differences.
package differ

// Result holds the outcome of comparing two .env environments.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched contains keys present in both environments whose values differ.
	Mismatched []MismatchedVar
}

// MismatchedVar describes a single variable whose value differs between environments.
type MismatchedVar struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare takes two env maps (key -> value) and returns a Result describing
// all missing and mismatched variables between them.
func Compare(left, right map[string]string) Result {
	var result Result

	for key, leftVal := range left {
		rightVal, ok := right[key]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, key)
			continue
		}
		if leftVal != rightVal {
			result.Mismatched = append(result.Mismatched, MismatchedVar{
				Key:        key,
				LeftValue:  leftVal,
				RightValue: rightVal,
			})
		}
	}

	for key := range right {
		if _, ok := left[key]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, key)
		}
	}

	return result
}

// HasDiff returns true when the Result contains any differences.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 ||
		len(r.MissingInLeft) > 0 ||
		len(r.Mismatched) > 0
}
