// Package scheduler provides a simple job runner that periodically diffs
// pairs of .env files and reports any changes.
//
// A Scheduler is created with one or more Job values, each specifying the
// paths of two env files to compare and the interval between checks.
// Output is written to a configurable io.Writer in either text or JSON
// format, matching the formats produced by the reporter package.
//
// Jobs can be defined programmatically or loaded from a JSON configuration
// file using LoadConfig and Config.ToJobs.
//
// Example:
//
//	cfg, _ := scheduler.LoadConfig("schedule.json")
//	jobs, _ := cfg.ToJobs()
//	s := scheduler.New(jobs, scheduler.DefaultOptions(os.Stdout))
//	_ = s.Run(context.Background())
package scheduler
