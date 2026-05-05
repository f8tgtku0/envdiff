package differ

import "fmt"

// ValuePair holds the differing values for a single key found in both
// environments but with non-matching values.
type ValuePair struct {
	Left  string
	Right string
}

// String returns a human-readable representation of the ValuePair.
func (vp ValuePair) String() string {
	return fmt.Sprintf("left=%q right=%q", vp.Left, vp.Right)
}

// Result holds the full diff between two .env maps.
type Result struct {
	// MissingInRight contains keys present in the left map but absent from the right.
	MissingInRight []string

	// MissingInLeft contains keys present in the right map but absent from the left.
	MissingInLeft []string

	// Mismatched contains keys present in both maps but with different values.
	Mismatched map[string]ValuePair
}

// HasDiff reports whether any differences were found.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 ||
		len(r.MissingInLeft) > 0 ||
		len(r.Mismatched) > 0
}

// Summary returns a brief human-readable summary of the diff result.
func (r Result) Summary() string {
	if !r.HasDiff() {
		return "no differences found"
	}
	return fmt.Sprintf(
		"%d missing in right, %d missing in left, %d mismatched",
		len(r.MissingInRight),
		len(r.MissingInLeft),
		len(r.Mismatched),
	)
}

// newResult returns an initialised Result ready for population.
func newResult() Result {
	return Result{
		MissingInRight: []string{},
		MissingInLeft:  []string{},
		Mismatched:     make(map[string]ValuePair),
	}
}
