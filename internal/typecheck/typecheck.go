package typecheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Type represents an expected value type for an env var.
type Type string

const (
	TypeString  Type = "string"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeBool    Type = "bool"
	TypeURL     Type = "url"
	TypeEmail   Type = "email"
)

var urlRe = regexp.MustCompile(`^https?://[^\s]+$`)
var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// Rule maps a key pattern to an expected type.
type Rule struct {
	KeyPattern string
	Type       Type
}

// Issue describes a type violation found during checking.
type Issue struct {
	Key      string
	Value    string
	Expected Type
	Reason   string
}

func (i Issue) String() string {
	return fmt.Sprintf("%-30s expected %s: %s", i.Key, i.Expected, i.Reason)
}

// Options controls typecheck behaviour.
type Options struct {
	Rules []Rule
}

// DefaultOptions returns a sensible default set of type rules.
func DefaultOptions() Options {
	return Options{
		Rules: []Rule{
			{KeyPattern: "_PORT$", Type: TypeInt},
			{KeyPattern: "_URL$", Type: TypeURL},
			{KeyPattern: "_EMAIL$", Type: TypeEmail},
			{KeyPattern: "_ENABLED$", Type: TypeBool},
			{KeyPattern: "_FLAG$", Type: TypeBool},
		},
	}
}

// Check validates env var values against the configured type rules.
// Keys are matched against each rule's KeyPattern (regexp). The first
// matching rule is applied.
func Check(env map[string]string, opts Options) ([]Issue, error) {
	var issues []Issue

	for _, rule := range opts.Rules {
		re, err := regexp.Compile(rule.KeyPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid key pattern %q: %w", rule.KeyPattern, err)
		}
		for key, val := range env {
			if !re.MatchString(key) {
				continue
			}
			if reason, ok := validate(val, rule.Type); !ok {
				issues = append(issues, Issue{
					Key:      key,
					Value:    val,
					Expected: rule.Type,
					Reason:   reason,
				})
			}
		}
	}
	return issues, nil
}

func validate(val string, t Type) (string, bool) {
	switch t {
	case TypeInt:
		if _, err := strconv.Atoi(strings.TrimSpace(val)); err != nil {
			return fmt.Sprintf("%q is not a valid integer", val), false
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err != nil {
			return fmt.Sprintf("%q is not a valid float", val), false
		}
	case TypeBool:
		v := strings.ToLower(strings.TrimSpace(val))
		if v != "true" && v != "false" && v != "1" && v != "0" && v != "yes" && v != "no" {
			return fmt.Sprintf("%q is not a valid boolean", val), false
		}
	case TypeURL:
		if !urlRe.MatchString(strings.TrimSpace(val)) {
			return fmt.Sprintf("%q is not a valid URL", val), false
		}
	case TypeEmail:
		if !emailRe.MatchString(strings.TrimSpace(val)) {
			return fmt.Sprintf("%q is not a valid email", val), false
		}
	}
	return "", true
}
