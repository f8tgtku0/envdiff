// Package annotator attaches inline comments to env keys.
//
// Given an environment map and a separate map of key→comment strings,
// Annotate produces a Result that pairs each key with its annotation.
// The accompanying WriteText helper renders the annotated env as a dotenv
// file, placing each comment on the line immediately above its key.
//
// Example usage:
//
//	notes := map[string]string{
//		"DATABASE_URL": "Primary Postgres connection string.",
//		"SECRET_KEY":   "Must be at least 32 characters.",
//	}
//	res, err := annotator.Annotate(env, notes, annotator.DefaultOptions())
//	if err != nil { ... }
//	annotator.WriteText(res, os.Stdout)
package annotator
