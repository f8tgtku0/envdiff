// Package redactor provides utilities for masking sensitive values in
// environment variable maps before they are displayed or exported.
//
// A Redactor is configured with a list of key-name patterns (e.g. "SECRET",
// "TOKEN", "PASSWORD") and a mask string. Any environment variable whose key
// contains one of the patterns (case-insensitive) is considered sensitive and
// its value is replaced with the mask when Redact is called.
//
// Example usage:
//
//	r := redactor.New(nil, "")        // use defaults
//	safe := r.Redact(loadedEnvMap)    // values like DB_PASSWORD become "***"
package redactor
