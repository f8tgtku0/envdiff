// Package differ compares two env maps and reports differences.
package differ

// Compare returns a Result describing the differences between left and right.
// Keys present in one but not the other are flagged as missing;
// keys present in both but with different values are flagged as mismatched.
func Compare(left, right map[string]string) *Result {
	r := newResult()

	for k, lv := range left {
		if rv, ok := right[k]; !ok {
			r.MissingInRight[k] = lv
		} else if lv != rv {
			r.Mismatched[k] = Pair{Left: lv, Right: rv}
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			r.MissingInLeft[k] = rv
		}
	}

	return r
}
