// Package highlighter annotates pairs of env maps with per-key diff status
// and renders the result as colour-coded text output.
//
// Each key is classified as one of:
//
//	"ok"            – present and identical in both envs
//	"mismatch"      – present in both envs but with different values
//	"missing_left"  – absent from the left env, present in the right
//	"missing_right" – present in the left env, absent from the right
//
// Basic usage:
//
//	lines := highlighter.Highlight(left, right, highlighter.DefaultOptions())
//	highlighter.WriteText(os.Stdout, lines, highlighter.DefaultOptions())
//	fmt.Println(highlighter.Summary(lines))
//
// Color output can be disabled by setting Options.UseColor = false, which is
// useful when writing to files or in CI pipelines that do not support ANSI
// escape codes.
package highlighter
