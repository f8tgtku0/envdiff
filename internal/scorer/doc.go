// Package scorer computes a numeric health score (0–100) for an env map
// relative to an expected set of keys.
//
// Penalties are applied for missing keys, empty values, and — when comparing
// two environments — mismatched values.  The final score is clamped to [0,
// 100] and translated to a letter grade (A–F).
//
// Basic usage:
//
//	result := scorer.Score(env, expectedKeys, scorer.DefaultOptions())
//	fmt.Println(result.Score, result.Grade)
//
// To compare two environments:
//
//	result := scorer.ScorePair(staging, production, expectedKeys, scorer.DefaultOptions())
package scorer
