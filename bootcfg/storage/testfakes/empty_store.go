package testfakes

import (
	"fmt"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

// EmptyStore is used for testing purposes.
type EmptyStore struct{}

// GroupGet returns a group not found error.
func (s *EmptyStore) GroupGet(id string) (*storagepb.Group, error) {
	return nil, fmt.Errorf("Group not found")
}

// GroupList returns an empty list of groups.
func (s *EmptyStore) GroupList() (groups []*storagepb.Group, err error) {
	return groups, nil
}

// ProfilePut returns an error writing any Profile.
func (s *EmptyStore) ProfilePut(profile *storagepb.Profile) error {
	return fmt.Errorf("emptyStore does not accept Profiles")
}

// ProfileGet returns a profile not found error.
func (s *EmptyStore) ProfileGet(id string) (*storagepb.Profile, error) {
	return nil, fmt.Errorf("Profile not found")
}

// ProfileList returns an empty list of profiles.
func (s *EmptyStore) ProfileList() (profiles []*storagepb.Profile, err error) {
	return profiles, nil
}

// IgnitionGet get returns an Ignition config not found error.
func (s *EmptyStore) IgnitionGet(id string) (string, error) {
	return "", fmt.Errorf("no Ignition Config %s", id)
}

// CloudGet returns a Cloud config not found error.
func (s *EmptyStore) CloudGet(id string) (string, error) {
	return "", fmt.Errorf("no Cloud Config %s", id)
}
