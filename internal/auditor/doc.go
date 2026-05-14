// Package auditor compares two env maps (a "before" and an "after" snapshot)
// and produces a structured audit report describing what changed.
//
// Each entry in the report carries:
//   - Key       – the environment variable name
//   - OldValue  – value before the change (empty for added keys)
//   - NewValue  – value after the change  (empty for removed keys)
//   - Action    – one of "added", "removed", "changed", or "unchanged"
//   - Timestamp – when the audit was run
//
// Usage:
//
//	report, err := auditor.Audit(before, after, auditor.DefaultOptions())
//
// Options:
//   - IncludeUnchanged – emit entries for keys whose value did not change
//   - RedactValues     – replace all values with "***" in the report
package auditor
