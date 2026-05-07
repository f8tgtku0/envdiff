package redactor

import "strings"

// DefaultSensitivePatterns holds common substrings that indicate a sensitive key.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

// Redactor masks values of sensitive environment variables.
type Redactor struct {
	patterns []string
	mask     string
}

// New creates a Redactor with the given sensitive key patterns and mask string.
// If patterns is nil, DefaultSensitivePatterns is used.
// If mask is empty, "***" is used.
func New(patterns []string, mask string) *Redactor {
	if patterns == nil {
		patterns = DefaultSensitivePatterns
	}
	if mask == "" {
		mask = "***"
	}
	return &Redactor{patterns: patterns, mask: mask}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Redact returns a copy of env where values of sensitive keys are replaced
// with the mask string.
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			out[k] = r.mask
		} else {
			out[k] = v
		}
	}
	return out
}
