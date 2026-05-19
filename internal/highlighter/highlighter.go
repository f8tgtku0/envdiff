package highlighter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls how diffs are highlighted in output.
type Options struct {
	// UseColor enables ANSI color codes in output.
	UseColor bool
	// MissingPrefix is the prefix printed before missing keys.
	MissingPrefix string
	// MismatchPrefix is the prefix printed before mismatched keys.
	MismatchPrefix string
}

// DefaultOptions returns sensible defaults for highlighting.
func DefaultOptions() Options {
	return Options{
		UseColor:       true,
		MissingPrefix:  "- ",
		MismatchPrefix: "~ ",
	}
}

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
)

// Line represents a single highlighted output line.
type Line struct {
	Key    string
	Left   string
	Right  string
	Status string // "missing_left", "missing_right", "mismatch", "ok"
}

// Highlight produces a slice of annotated Lines from two env maps.
func Highlight(left, right map[string]string, opts Options) []Line {
	keys := unionKeys(left, right)
	sort.Strings(keys)

	lines := make([]Line, 0, len(keys))
	for _, k := range keys {
		lv, lok := left[k]
		rv, rok := right[k]
		switch {
		case lok && !rok:
			lines = append(lines, Line{Key: k, Left: lv, Right: "", Status: "missing_right"})
		case !lok && rok:
			lines = append(lines, Line{Key: k, Left: "", Right: rv, Status: "missing_left"})
		case lv != rv:
			lines = append(lines, Line{Key: k, Left: lv, Right: rv, Status: "mismatch"})
		default:
			lines = append(lines, Line{Key: k, Left: lv, Right: rv, Status: "ok"})
		}
	}
	return lines
}

// WriteText writes highlighted diff output to w.
func WriteText(w io.Writer, lines []Line, opts Options) error {
	if w == nil {
		return fmt.Errorf("highlighter: writer must not be nil")
	}
	for _, l := range lines {
		switch l.Status {
		case "missing_left", "missing_right":
			prefix := opts.MissingPrefix
			line := fmt.Sprintf("%s%s (left=%q right=%q)\n", prefix, l.Key, l.Left, l.Right)
			if opts.UseColor {
				line = colorRed + line + colorReset
			}
			fmt.Fprint(w, line)
		case "mismatch":
			prefix := opts.MismatchPrefix
			line := fmt.Sprintf("%s%s (left=%q right=%q)\n", prefix, l.Key, l.Left, l.Right)
			if opts.UseColor {
				line = colorYellow + line + colorReset
			}
			fmt.Fprint(w, line)
		default:
			line := fmt.Sprintf("  %s=%q\n", l.Key, l.Left)
			if opts.UseColor {
				line = colorGreen + line + colorReset
			}
			fmt.Fprint(w, line)
		}
	}
	return nil
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

// Summary returns a one-line human-readable summary of the highlighted lines.
func Summary(lines []Line) string {
	var missing, mismatch, ok int
	for _, l := range lines {
		switch l.Status {
		case "missing_left", "missing_right":
			missing++
		case "mismatch":
			mismatch++
		default:
			ok++
		}
	}
	parts := []string{fmt.Sprintf("%d ok", ok)}
	if missing > 0 {
		parts = append(parts, fmt.Sprintf("%d missing", missing))
	}
	if mismatch > 0 {
		parts = append(parts, fmt.Sprintf("%d mismatched", mismatch))
	}
	return strings.Join(parts, ", ")
}
