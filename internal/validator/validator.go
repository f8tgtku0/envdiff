package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule applied to env variable keys or values.
type Rule struct {
	// KeyPattern, if non-empty, restricts the rule to keys matching this regex.
	KeyPattern string
	// Required means the key must be present and non-empty.
	Required bool
	// ValuePattern, if non-empty, means the value must match this regex.
	ValuePattern string
}

// Violation describes a single validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks the provided env map against the given rules and returns
// any violations found. An empty slice means the env is valid.
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		var keyRe *regexp.Regexp
		if rule.KeyPattern != "" {
			var err error
			keyRe, err = regexp.Compile(rule.KeyPattern)
			if err != nil {
				violations = append(violations, Violation{
					Key:     "<rule>",
					Message: fmt.Sprintf("invalid key pattern %q: %v", rule.KeyPattern, err),
				})
				continue
			}
		}

		for key, val := range env {
			if keyRe != nil && !keyRe.MatchString(key) {
				continue
			}

			if rule.Required && strings.TrimSpace(val) == "" {
				violations = append(violations, Violation{
					Key:     key,
					Message: "required but empty or missing value",
				})
			}

			if rule.ValuePattern != "" {
				valRe, err := regexp.Compile(rule.ValuePattern)
				if err != nil {
					violations = append(violations, Violation{
						Key:     key,
						Message: fmt.Sprintf("invalid value pattern %q: %v", rule.ValuePattern, err),
					})
					continue
				}
				if !valRe.MatchString(val) {
					violations = append(violations, Violation{
						Key:     key,
						Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.ValuePattern),
					})
				}
			}
		}
	}

	return violations
}
