package scheduler_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/envdiff/internal/scheduler"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestNew_DefaultOptions(t *testing.T) {
	s := scheduler.New(nil, scheduler.Options{})
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRun_NoJobs_ReturnsError(t *testing.T) {
	s := scheduler.New(nil, scheduler.DefaultOptions(nil))
	ctx := context.Background()
	if err := s.Run(ctx); err == nil {
		t.Fatal("expected error for empty job list")
	}
}

func TestRun_TicksAndWrites(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	right := writeTempEnv(t, "FOO=bar\n")

	var buf bytes.Buffer
	jobs := []scheduler.Job{
		{Name: "test", Left: left, Right: right, Interval: 20 * time.Millisecond},
	}
	s := scheduler.New(jobs, scheduler.DefaultOptions(&buf))

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	_ = s.Run(ctx)

	if buf.Len() == 0 {
		t.Error("expected output written to writer")
	}
}

func TestRun_OnDiff_CalledWhenDiffsExist(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nMISSING=val\n")
	right := writeTempEnv(t, "FOO=bar\n")

	called := make(chan struct{}, 1)
	opts := scheduler.DefaultOptions(nil)
	opts.OnDiff = func(j scheduler.Job, _ interface{ Clean() bool }) {
		called <- struct{}{}
	}

	jobs := []scheduler.Job{
		{Name: "diff-job", Left: left, Right: right, Interval: 20 * time.Millisecond},
	}
	s := scheduler.New(jobs, opts)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go func() { _ = s.Run(ctx) }()

	select {
	case <-called:
		// success
	case <-ctx.Done():
		t.Error("OnDiff was never called")
	}
}

func TestRun_JSONFormat(t *testing.T) {
	left := writeTempEnv(t, "A=1\n")
	right := writeTempEnv(t, "A=1\n")

	var buf bytes.Buffer
	opts := scheduler.DefaultOptions(&buf)
	opts.Format = "json"

	jobs := []scheduler.Job{
		{Name: "json-job", Left: left, Right: right, Interval: 20 * time.Millisecond},
	}
	s := scheduler.New(jobs, opts)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	_ = s.Run(ctx)

	out := buf.String()
	if len(out) == 0 {
		t.Error("expected JSON output")
	}
	_ = filepath.Join // suppress unused import
}
