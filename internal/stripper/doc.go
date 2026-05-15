// Package stripper provides a way to remove keys from an env map based on
// configurable rules: prefix matching, suffix matching, regular-expression
// patterns, or an explicit list of key names.
//
// Basic usage:
//
//	opts := stripper.DefaultOptions()
//	opts.Prefixes = []string{"INTERNAL_", "DEBUG_"}
//	opts.Keys    = []string{"LEGACY_TOKEN"}
//
//	result, err := stripper.Strip(env, opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("stripped:", result.Stripped)
//
// Set DryRun to true to inspect which keys would be removed without
// modifying the returned Env map.
package stripper
