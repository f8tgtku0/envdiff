// Package deduplicator merges multiple env maps into a single canonical map
// and reports any keys whose values differ across sources.
//
// # Usage
//
//	sources := []map[string]string{
//		{"DB_HOST": "localhost", "PORT": "5432"},
//		{"DB_HOST": "prod.db",   "LOG_LEVEL": "info"},
//	}
//
//	merged, report := deduplicator.Deduplicate(sources, deduplicator.DefaultOptions())
//	if report.HasConflicts() {
//		for _, c := range report.Conflicts {
//			fmt.Printf("conflict: %s => %v\n", c.Key, c.Values)
//		}
//	}
//
// By default the first value seen for a key wins. Set Options.PreferLast to
// true to keep the last value instead. Set Options.ReportOnly to skip
// mutation and only collect the conflict report.
package deduplicator
