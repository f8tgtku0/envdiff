package validator

import (
	"encoding/json"
	"fmt"
	"os"
)

// RulesFile is the on-disk representation of a rules configuration.
type RulesFile struct {
	Rules []Rule `json:"rules"`
}

// LoadRules reads a JSON file at path and returns the parsed rules.
// The expected format is:
//
//	{
//	  "rules": [
//	    {"key_pattern": "^DB_", "required": true},
//	    {"key_pattern": "^PORT$", "value_pattern": "^\\d+$"}
//	  ]
//	}
func LoadRules(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("validator: open rules file: %w", err)
	}
	defer f.Close()

	var rf rulesFileJSON
	if err := json.NewDecoder(f).Decode(&rf); err != nil {
		return nil, fmt.Errorf("validator: decode rules file: %w", err)
	}

	rules := make([]Rule, len(rf.Rules))
	for i, r := range rf.Rules {
		rules[i] = Rule{
			KeyPattern:   r.KeyPattern,
			Required:     r.Required,
			ValuePattern: r.ValuePattern,
		}
	}
	return rules, nil
}

// rulesFileJSON mirrors Rule with json struct tags for unmarshalling.
type rulesFileJSON struct {
	Rules []struct {
		KeyPattern   string `json:"key_pattern"`
		Required     bool   `json:"required"`
		ValuePattern string `json:"value_pattern"`
	} `json:"rules"`
}
