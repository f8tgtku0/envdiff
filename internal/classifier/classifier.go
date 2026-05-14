// Package classifier categorises env variables by their likely purpose
// based on key naming patterns (e.g. database, auth, network, feature flags).
package classifier

import (
	"regexp"
	"sort"
)

// Category represents a named grouping for env variables.
type Category string

const (
	CategoryDatabase Category = "database"
	CategoryAuth     Category = "auth"
	CategoryNetwork  Category = "network"
	CategoryFeature  Category = "feature_flag"
	CategoryLogging  Category = "logging"
	CategoryUnknown  Category = "unknown"
)

// Rule maps a compiled pattern to a category.
type Rule struct {
	Pattern  *regexp.Regexp
	Category Category
}

// Result holds the classification output for a single env map.
type Result struct {
	// Categories maps each Category to the list of keys assigned to it.
	Categories map[Category][]string
}

// DefaultRules returns the built-in classification rules.
func DefaultRules() []Rule {
	return []Rule{
		{regexp.MustCompile(`(?i)(DB_|DATABASE_|POSTGRES_|MYSQL_|MONGO_|REDIS_)`), CategoryDatabase},
		{regexp.MustCompile(`(?i)(SECRET|TOKEN|PASSWORD|API_KEY|AUTH_|JWT_|OAUTH_)`), CategoryAuth},
		{regexp.MustCompile(`(?i)(HOST|PORT|URL|ADDR|ENDPOINT|DOMAIN|PROXY_)`), CategoryNetwork},
		{regexp.MustCompile(`(?i)(FEATURE_|FLAG_|ENABLE_|DISABLE_)`), CategoryFeature},
		{regexp.MustCompile(`(?i)(LOG_|LOGGING_|LOG$)`), CategoryLogging},
	}
}

// Classify assigns each key in env to a category using the provided rules.
// Keys that match no rule are placed under CategoryUnknown.
// If rules is nil, DefaultRules are used.
func Classify(env map[string]string, rules []Rule) Result {
	if rules == nil {
		rules = DefaultRules()
	}

	out := Result{Categories: make(map[Category][]string)}

	for key := range env {
		cat := categorise(key, rules)
		out.Categories[cat] = append(out.Categories[cat], key)
	}

	// Sort keys within each category for deterministic output.
	for cat := range out.Categories {
		sort.Strings(out.Categories[cat])
	}

	return out
}

func categorise(key string, rules []Rule) Category {
	for _, r := range rules {
		if r.Pattern.MatchString(key) {
			return r.Category
		}
	}
	return CategoryUnknown
}
