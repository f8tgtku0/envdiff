package auditor

import (
	"fmt"
	"sort"
	"time"
)

// Entry records a single audit event for a key-value change.
type Entry struct {
	Timestamp time.Time
	Key       string
	OldValue  string
	NewValue  string
	Action    string // "added", "removed", "changed", "unchanged"
}

// Report holds all audit entries produced by a diff.
type Report struct {
	Entries   []Entry
	CreatedAt time.Time
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		IncludeUnchanged: false,
		RedactValues:     false,
	}
}

// Options controls Audit behaviour.
type Options struct {
	IncludeUnchanged bool
	RedactValues     bool
}

const redacted = "***"

// Audit compares two env maps (before, after) and produces an audit Report.
func Audit(before, after map[string]string, opts Options) (*Report, error) {
	if before == nil {
		return nil, fmt.Errorf("auditor: before env must not be nil")
	}
	if after == nil {
		return nil, fmt.Errorf("auditor: after env must not be nil")
	}

	keys := unionKeys(before, after)
	report := &Report{CreatedAt: time.Now()}

	for _, k := range keys {
		oldVal, inBefore := before[k]
		newVal, inAfter := after[k]

		var action string
		switch {
		case inBefore && !inAfter:
			action = "removed"
		case !inBefore && inAfter:
			action = "added"
			oldVal = ""
		case oldVal != newVal:
			action = "changed"
		default:
			action = "unchanged"
		}

		if action == "unchanged" && !opts.IncludeUnchanged {
			continue
		}

		if opts.RedactValues {
			if oldVal != "" {
				oldVal = redacted
			}
			if newVal != "" {
				newVal = redacted
			}
		}

		report.Entries = append(report.Entries, Entry{
			Timestamp: report.CreatedAt,
			Key:       k,
			OldValue:  oldVal,
			NewValue:  newVal,
			Action:    action,
		})
	}
	return report, nil
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
