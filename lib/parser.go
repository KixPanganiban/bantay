package lib

import (
	"gopkg.in/yaml.v2"
)

// ParsedChecks represents the unmarshalled YAML file
type ParsedChecks struct {
	Checks *[]Check `yaml:"checks"`
}

// ParseYAML parses the given YAML File and outputs a set of Checks
func ParseYAML(b []byte) (*[]Check, error) {
	checks := ParsedChecks{}
	err := yaml.Unmarshal(b, &checks)
	if err != nil {
		return nil, err
	}
	return checks.Checks, nil
}
