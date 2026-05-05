// Package filter provides key-level filtering for environment variable maps.
//
// It supports two complementary filtering mechanisms:
//
//   - IncludePrefixes: when set, only keys whose names begin with one of the
//     specified prefixes are retained. An empty slice means "include all".
//
//   - ExcludePrefixes: keys whose names begin with one of the specified
//     prefixes are removed from the result. Exclusion is evaluated after
//     inclusion, so it can be used to carve out specific sub-sets.
//
// Typical usage:
//
//	opts := filter.Options{
//	    IncludePrefixes: []string{"APP_", "SERVICE_"},
//	    ExcludePrefixes: []string{"APP_SECRET_"},
//	}
//	filtered := filter.Apply(rawEnv, opts)
package filter
