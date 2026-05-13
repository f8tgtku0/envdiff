// Package tagger assigns user-defined labels to environment variable keys
// based on regular-expression rules.
//
// Labels (tags) are matched against key names, making it easy to categorise
// variables — for example flagging secrets, database credentials, or feature
// flags — without inspecting their values.
//
// Basic usage:
//
//	env := map[string]string{
//		"DB_PASSWORD": "hunter2",
//		"APP_NAME":    "envdiff",
//	}
//
//	res, err := tagger.Tag(env, tagger.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	// res.Tags["DB_PASSWORD"] => ["secret"]
//	// res.Untagged           => ["APP_NAME"]
package tagger
