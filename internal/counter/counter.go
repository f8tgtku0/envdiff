// Package counter provides utilities for counting and summarising
// key statistics across one or more parsed environment maps.
package counter

import "sort"

// Result holds the counts produced by Count.
type Result struct {
	// Total is the number of keys across all envs.
	Total int
	// PerEnv maps each environment label to its key count.
	PerEnv map[string]int
	// UniqueKeys is the number of distinct keys seen across all envs.
	UniqueKeys int
	// EmptyValues is the number of keys whose value is the empty string.
	EmptyValues int
}

// Count tallies key statistics for the supplied labelled environments.
// The map keys are human-readable labels (e.g. ".env.staging").
func Count(envs map[string]map[string]string) Result {
	if len(envs) == 0 {
		return Result{PerEnv: map[string]int{}}
	}

	allKeys := make(map[string]struct{})
	perEnv := make(map[string]int, len(envs))
	total := 0
	empty := 0

	for label, env := range envs {
		perEnv[label] = len(env)
		total += len(env)
		for k, v := range env {
			allKeys[k] = struct{}{}
			if v == "" {
				empty++
			}
		}
	}

	return Result{
		Total:       total,
		PerEnv:      perEnv,
		UniqueKeys:  len(allKeys),
		EmptyValues: empty,
	}
}

// SortedLabels returns the environment labels in alphabetical order.
func SortedLabels(r Result) []string {
	labels := make([]string, 0, len(r.PerEnv))
	for l := range r.PerEnv {
		labels = append(labels, l)
	}
	sort.Strings(labels)
	return labels
}
