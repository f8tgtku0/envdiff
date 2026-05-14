// Package freezer provides enforcement of a "frozen" .env snapshot.
//
// A frozen snapshot captures the exact set of keys and values that are
// considered stable for a given environment. Freeze then compares a live
// env map against that snapshot and reports any violations:
//
//   - Keys present in the frozen snapshot but absent from the live env.
//   - Values that have changed relative to the frozen snapshot.
//   - Keys added to the live env that were not in the frozen snapshot
//     (controlled by Options.AllowExpand).
//
// Typical usage:
//
//	result, err := freezer.Freeze(snapshot, live, freezer.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !result.Clean() {
//		freezer.WriteText(result, os.Stderr)
//	}
package freezer
