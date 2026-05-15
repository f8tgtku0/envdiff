// Package differ compares two env maps and returns a structured result.
package differ

// Compare takes two env maps (left and right) and returns a Result describing
// keys that are missing in either side or have mismatched values.
func Compare(left, right map[string]string) Result {
	r := newResult()

	for k, lv := range left {
		if rv, ok := right[k]; !ok {
			r.MissingInRight[k] = lv
		} else if lv != rv {
			r.Mismatched[k] = [2]string{lv, rv}
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			r.MissingInLeft[k] = rv
		}
	}

	return r
}
