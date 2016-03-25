package storage

import (
	"fmt"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

// fixedStore is used for testing purposes.
type fixedStore struct {
	Groups          map[string]*storagepb.Group
	Profiles        map[string]*storagepb.Profile
	IgnitionConfigs map[string]string
	CloudConfigs    map[string]string
}

func (s *fixedStore) GroupGet(id string) (*storagepb.Group, error) {
	if group, present := s.Groups[id]; present {
		return group, nil
	}
	return nil, ErrGroupNotFound
}

func (s *fixedStore) GroupList() ([]*storagepb.Group, error) {
	groups := make([]*storagepb.Group, len(s.Groups))
	i := 0
	for _, g := range s.Groups {
		groups[i] = g
		i++
	}
	return groups, nil
}

func (s *fixedStore) ProfileGet(id string) (*storagepb.Profile, error) {
	if profile, present := s.Profiles[id]; present {
		return profile, nil
	}
	return nil, ErrProfileNotFound
}

func (s *fixedStore) ProfileList() ([]*storagepb.Profile, error) {
	profiles := make([]*storagepb.Profile, len(s.Profiles))
	i := 0
	for _, p := range s.Profiles {
		profiles[i] = p
		i++
	}
	return profiles, nil
}

func (s *fixedStore) IgnitionGet(id string) (string, error) {
	if config, present := s.IgnitionConfigs[id]; present {
		return config, nil
	}
	return "", fmt.Errorf("no Ignition Config %s", id)
}

func (s *fixedStore) CloudGet(id string) (string, error) {
	if config, present := s.CloudConfigs[id]; present {
		return config, nil
	}
	return "", fmt.Errorf("no Cloud Config %s", id)
}
