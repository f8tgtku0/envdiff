package scheduler_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/your-org/envdiff/internal/scheduler"
)

func writeTempConfig(t *testing.T, cfg scheduler.Config) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write(data)
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	cfg := scheduler.Config{
		Format: "text",
		Jobs: []scheduler.JobConfig{
			{Name: "prod", Left: "prod.env", Right: "staging.env", Interval: "1m"},
		},
	}
	path := writeTempConfig(t, cfg)
	got, err := scheduler.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Format != "text" {
		t.Errorf("format: got %q want %q", got.Format, "text")
	}
	if len(got.Jobs) != 1 {
		t.Fatalf("jobs: got %d want 1", len(got.Jobs))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := scheduler.LoadConfig("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestToJobs_Valid(t *testing.T) {
	cfg := &scheduler.Config{
		Jobs: []scheduler.JobConfig{
			{Name: "dev", Left: "dev.env", Right: "prod.env", Interval: "30s"},
		},
	}
	jobs, err := cfg.ToJobs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Interval != 30*time.Second {
		t.Errorf("interval: got %v want 30s", jobs[0].Interval)
	}
}

func TestToJobs_MissingName(t *testing.T) {
	cfg := &scheduler.Config{
		Jobs: []scheduler.JobConfig{
			{Left: "a.env", Right: "b.env", Interval: "1m"},
		},
	}
	if _, err := cfg.ToJobs(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestToJobs_InvalidInterval(t *testing.T) {
	cfg := &scheduler.Config{
		Jobs: []scheduler.JobConfig{
			{Name: "x", Left: "a.env", Right: "b.env", Interval: "bad"},
		},
	}
	if _, err := cfg.ToJobs(); err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestToJobs_EmptyJobs(t *testing.T) {
	cfg := &scheduler.Config{}
	if _, err := cfg.ToJobs(); err == nil {
		t.Fatal("expected error for empty jobs")
	}
}
