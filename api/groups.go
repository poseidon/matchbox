package api

import (
	"errors"

	"gopkg.in/yaml.v2"
)

// GroupConfig parser errors
var (
	ErrInvalidVersion = errors.New("api: mismatched API version")
)

// Group associates matcher conditions with a Specification identifier.
type Group struct {
	// Human readable name (optional)
	Name string `yaml:"name"`
	// Spec identifier
	Specification string `yaml:"spec"`
	// matcher conditions
	Matcher RequirementSet `yaml:"require"`
}

// GroupConfig define an group import structure.
type GroupConfig struct {
	APIVersion string  `yaml:"api_version"`
	Groups     []Group `yaml:"groups"`
}

// ParseGroupConfig parses a YAML group config and returns a GroupConfig.
func ParseGroupConfig(data []byte) (*GroupConfig, error) {
	config := new(GroupConfig)
	err := yaml.Unmarshal(data, config)
	if err == nil && config.APIVersion != APIVersion {
		return nil, ErrInvalidVersion
	}
	return config, err
}
