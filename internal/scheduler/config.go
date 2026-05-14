package scheduler

import (
	"fmt"
	"os"
	"time"

	"encoding/json"
)

// JobConfig is the serialisable representation of a Job loaded from a file.
type JobConfig struct {
	Name     string `json:"name"`
	Left     string `json:"left"`
	Right    string `json:"right"`
	Interval string `json:"interval"`
}

// Config holds the top-level scheduler configuration file structure.
type Config struct {
	Format string      `json:"format"`
	Jobs   []JobConfig `json:"jobs"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("scheduler: read config %q: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("scheduler: parse config: %w", err)
	}
	return &cfg, nil
}

// ToJobs converts the config's JobConfig entries into Job values.
func (c *Config) ToJobs() ([]Job, error) {
	if len(c.Jobs) == 0 {
		return nil, fmt.Errorf("scheduler: config contains no jobs")
	}
	jobs := make([]Job, 0, len(c.Jobs))
	for _, jc := range c.Jobs {
		if jc.Name == "" {
			return nil, fmt.Errorf("scheduler: job missing name")
		}
		if jc.Left == "" || jc.Right == "" {
			return nil, fmt.Errorf("scheduler: job %q missing left or right path", jc.Name)
		}
		d, err := time.ParseDuration(jc.Interval)
		if err != nil {
			return nil, fmt.Errorf("scheduler: job %q invalid interval %q: %w", jc.Name, jc.Interval, err)
		}
		if d <= 0 {
			return nil, fmt.Errorf("scheduler: job %q interval must be positive", jc.Name)
		}
		jobs = append(jobs, Job{Name: jc.Name, Left: jc.Left, Right: jc.Right, Interval: d})
	}
	return jobs, nil
}
