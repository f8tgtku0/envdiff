// Package scheduler provides utilities for scheduling periodic env diff
// checks across multiple environment file pairs.
package scheduler

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/loader"
	"github.com/your-org/envdiff/internal/reporter"
)

// Job represents a single scheduled diff task between two env files.
type Job struct {
	Name     string
	Left     string
	Right    string
	Interval time.Duration
}

// Options configures the Scheduler behaviour.
type Options struct {
	// Writer receives report output for every tick. Defaults to io.Discard.
	Writer io.Writer
	// Format is "text" or "json". Defaults to "text".
	Format string
	// OnDiff is called after each tick when differences are found.
	OnDiff func(job Job, result *differ.Result)
}

// DefaultOptions returns a safe default Options value.
func DefaultOptions(w io.Writer) Options {
	return Options{
		Writer: w,
		Format: "text",
	}
}

// Scheduler runs one or more Jobs on a fixed interval.
type Scheduler struct {
	jobs []Job
	opts Options
}

// New creates a Scheduler with the provided jobs and options.
func New(jobs []Job, opts Options) *Scheduler {
	if opts.Writer == nil {
		opts.Writer = io.Discard
	}
	if opts.Format == "" {
		opts.Format = "text"
	}
	return &Scheduler{jobs: jobs, opts: opts}
}

// Run starts all jobs and blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	if len(s.jobs) == 0 {
		return fmt.Errorf("scheduler: no jobs configured")
	}
	for _, j := range s.jobs {
		go s.loop(ctx, j)
	}
	<-ctx.Done()
	return ctx.Err()
}

func (s *Scheduler) loop(ctx context.Context, j Job) {
	tick := time.NewTicker(j.Interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			s.run(j)
		}
	}
}

func (s *Scheduler) run(j Job) {
	pair, err := loader.LoadPair(j.Left, j.Right, loader.Options{})
	if err != nil {
		fmt.Fprintf(s.opts.Writer, "[%s] load error: %v\n", j.Name, err)
		return
	}
	result := differ.Compare(pair.Left, pair.Right)
	if s.opts.OnDiff != nil && !result.Clean() {
		s.opts.OnDiff(j, result)
	}
	rep := reporter.NewReport(j.Left, j.Right, result)
	switch s.opts.Format {
	case "json":
		_ = reporter.WriteJSON(s.opts.Writer, rep)
	default:
		_ = reporter.WriteText(s.opts.Writer, rep)
	}
}
