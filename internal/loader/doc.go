// Package loader handles reading one or more .env files from disk, delegating
// parsing to the parser package and key filtering to the filter package.
//
// Typical usage:
//
//	left, right, err := loader.LoadPair(".env.staging", ".env.production", loader.Options{
//		IncludePrefixes: []string{"APP_"},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	// left.Vars and right.Vars are ready for diffing.
//
// The Options struct mirrors the filter.Options fields so callers do not need
// to import the filter package directly when using the loader.
package loader
