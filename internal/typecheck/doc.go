// Package typecheck validates env var values against expected types.
//
// Rules map key patterns (regular expressions) to a Type such as int,
// bool, url, or email. Check iterates over an env map and returns an
// Issue for every value that does not satisfy its matched rule.
//
// DefaultOptions provides a ready-made rule set covering common
// conventions (e.g. keys ending in _PORT must be integers, _URL must
// be a valid HTTP/S URL).
//
// Results can be rendered as plain text or JSON via WriteText and
// WriteJSON respectively.
package typecheck
