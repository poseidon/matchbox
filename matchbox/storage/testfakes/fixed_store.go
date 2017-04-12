package testfakes

import (
	"fmt"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// FixedStore is used for testing purposes.
type FixedStore struct {
	Groups          map[string]*storagepb.Group
	Profiles        map[string]*storagepb.Profile
	IgnitionConfigs map[string]string
	CloudConfigs    map[string]string
	GenericConfigs  map[string]string
}

// NewFixedStore returns a new FixedStore.
func NewFixedStore() *FixedStore {
	return &FixedStore{
		Groups:          make(map[string]*storagepb.Group),
		Profiles:        make(map[string]*storagepb.Profile),
		IgnitionConfigs: make(map[string]string),
		CloudConfigs:    make(map[string]string),
		GenericConfigs:  make(map[string]string),
	}
}

// GroupPut write the given Group the Groups map.
func (s *FixedStore) GroupPut(group *storagepb.Group) error {
	s.Groups[group.Id] = group
	return nil
}

// GroupGet returns the Group from the Groups map with the given id.
func (s *FixedStore) GroupGet(id string) (*storagepb.Group, error) {
	if group, present := s.Groups[id]; present {
		return group, nil
	}
	return nil, fmt.Errorf("Group not found")
}

// GroupDelete deletes the Group from the Groups map with the given id.
func (s *FixedStore) GroupDelete(id string) error {
	delete(s.Groups, id)
	return nil
}

// GroupList returns the groups in the Groups map.
func (s *FixedStore) GroupList() ([]*storagepb.Group, error) {
	groups := make([]*storagepb.Group, len(s.Groups))
	i := 0
	for _, g := range s.Groups {
		groups[i] = g
		i++
	}
	return groups, nil
}

// ProfilePut writes the given Profile to the Profiles map.
func (s *FixedStore) ProfilePut(profile *storagepb.Profile) error {
	s.Profiles[profile.Id] = profile
	return nil
}

// ProfileGet returns the Profile from the Profiles map with the given id.
func (s *FixedStore) ProfileGet(id string) (*storagepb.Profile, error) {
	if profile, present := s.Profiles[id]; present {
		return profile, nil
	}
	return nil, fmt.Errorf("Profile not found")
}

// ProfileDelete deletes the Profile from the Profiles map with the given id.
func (s *FixedStore) ProfileDelete(id string) error {
	delete(s.Profiles, id)
	return nil
}

// ProfileList returns the profiles in the Profiles map.
func (s *FixedStore) ProfileList() ([]*storagepb.Profile, error) {
	profiles := make([]*storagepb.Profile, len(s.Profiles))
	i := 0
	for _, p := range s.Profiles {
		profiles[i] = p
		i++
	}
	return profiles, nil
}

// IgnitionPut create or updates an Ignition template.
func (s *FixedStore) IgnitionPut(name string, config []byte) error {
	s.IgnitionConfigs[name] = string(config)
	return nil
}

// IgnitionGet returns an Ignition template by name.
func (s *FixedStore) IgnitionGet(name string) (string, error) {
	if config, present := s.IgnitionConfigs[name]; present {
		return config, nil
	}
	return "", fmt.Errorf("no Ignition template %s", name)
}

// IgnitionDelete deletes an Ignition template by name.
func (s *FixedStore) IgnitionDelete(name string) error {
	delete(s.IgnitionConfigs, name)
	return nil
}

// CloudGet returns a Cloud-config template by name.
func (s *FixedStore) CloudGet(name string) (string, error) {
	if config, present := s.CloudConfigs[name]; present {
		return config, nil
	}
	return "", fmt.Errorf("no Cloud-Config template %s", name)
}

// GenericGet returns a generic template by name.
func (s *FixedStore) GenericGet(name string) (string, error) {
	if config, present := s.GenericConfigs[name]; present {
		return config, nil
	}
	return "", fmt.Errorf("no generic template %s", name)
}
