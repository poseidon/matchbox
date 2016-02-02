package api

import (
	"fmt"
)

// fixedStore provides fixed Group, Spec, Ignition, and Cloud config resources
// for testing.
type fixedStore struct {
	Groups          []Group
	Specs           map[string]*Spec
	CloudConfigs    map[string]string
	IgnitionConfigs map[string]string
}

func (s *fixedStore) Spec(id string) (*Spec, error) {
	if spec, present := s.Specs[id]; present {
		return spec, nil
	}
	return nil, fmt.Errorf("no spec %s", id)
}

func (s *fixedStore) BootstrapGroups(groups []Group) error {
	s.Groups = groups
	return nil
}

func (s *fixedStore) ListGroups() ([]Group, error) {
	return s.Groups, nil
}

func (s *fixedStore) CloudConfig(id string) (string, error) {
	if config, present := s.CloudConfigs[id]; present {
		return config, nil
	}
	return "", fmt.Errorf("no cloud config %s", id)
}

func (s *fixedStore) IgnitionConfig(id string) (string, error) {
	if config, present := s.IgnitionConfigs[id]; present {
		return config, nil
	}
	return "", fmt.Errorf("no ignition config %s", id)
}

type emptyStore struct{}

func (s *emptyStore) Spec(id string) (*Spec, error) {
	return nil, fmt.Errorf("no group config %s", id)
}

func (s *emptyStore) BootstrapGroups(groups []Group) error {
	return nil
}

func (s *emptyStore) ListGroups() (groups []Group, err error) {
	return groups, nil
}

func (s emptyStore) CloudConfig(id string) (string, error) {
	return "", fmt.Errorf("no cloud config %s", id)
}

func (s emptyStore) IgnitionConfig(id string) (string, error) {
	return "", fmt.Errorf("no ignition config %s", id)
}
