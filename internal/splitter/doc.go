// Package splitter partitions a flat env map into named buckets based on
// key prefixes.
//
// Given an env map such as:
//
//	DB_HOST=localhost
//	DB_PORT=5432
//	APP_NAME=myapp
//	DEBUG=true
//
// and prefixes ["DB", "APP"], Split produces three buckets:
//
//	"DB"  -> {HOST: localhost, PORT: 5432}
//	"APP" -> {NAME: myapp}
//	""    -> {DEBUG: true}   // unmatched
//
// The separator between prefix and key name defaults to "_" and is
// configurable via Options.PrefixSep.
package splitter
