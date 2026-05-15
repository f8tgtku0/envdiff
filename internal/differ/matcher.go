package differ

import "regexp"

// MatchOptions controls fuzzy / pattern-based matching behaviour when
// comparing env values.
type MatchOptions struct {
	// IgnoreCase treats values as equal regardless of casing.
	IgnoreCase bool
	// ValuePatterns maps key names (or glob prefixes) to regexp patterns that
	// both sides must satisfy for the key to be considered "matching". When a
	// pattern is provided the exact values are not compared — only conformance
	// to the pattern is checked.
	ValuePatterns map[string]string
}

// MatchResult holds the outcome of a pattern-aware comparison for a single key.
type MatchResult struct {
	Key           string
	LeftValue     string
	RightValue    string
	// Conformant is true when both values satisfy the declared pattern.
	Conformant    bool
	// PatternUsed is the regexp pattern that was applied, if any.
	PatternUsed   string
}

// MatchValues compares two env maps using the supplied MatchOptions and returns
// a slice of MatchResult for every key that appears in both maps.
func MatchValues(left, right map[string]string, opts MatchOptions) []MatchResult {
	var results []MatchResult

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			continue
		}

		mr := MatchResult{Key: k, LeftValue: lv, RightValue: rv}

		if pat, hasPat := opts.ValuePatterns[k]; hasPat {
			re, err := regexp.Compile(pat)
			if err == nil {
				mr.PatternUsed = pat
				mr.Conformant = re.MatchString(lv) && re.MatchString(rv)
			}
		} else if opts.IgnoreCase {
			mr.Conformant = equalFold(lv, rv)
		} else {
			mr.Conformant = lv == rv
		}

		results = append(results, mr)
	}

	return results
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 32
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}
