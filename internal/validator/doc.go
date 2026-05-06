// Package validator provides rule-based validation for environment variable
// maps produced by the loader package.
//
// A Rule can target a subset of keys via a regular-expression pattern and
// enforce one or more constraints:
//
//   - Required — the key must be present with a non-empty value.
//   - ValuePattern — the value must match the supplied regular expression.
//
// Usage:
//
//	rules := []validator.Rule{
//		{KeyPattern: `^DB_`, Required: true},
//		{KeyPattern: `^PORT$`, ValuePattern: `^\d+$`},
//	}
//	violations := validator.Validate(env, rules)
//	for _, v := range violations {
//		fmt.Println(v)
//	}
package validator
