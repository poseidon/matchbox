package testfakes

import (
	"fmt"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// EmptyStore is used for testing purposes.
type EmptyStore struct{}

// GroupPut returns an error writing any Group.
func (s *EmptyStore) GroupPut(group *storagepb.Group) error {
	return fmt.Errorf("emptyStore does not accept Groups")
}

// GroupGet returns a group not found error.
func (s *EmptyStore) GroupGet(id string) (*storagepb.Group, error) {
	return nil, fmt.Errorf("Group not found")
}

// GroupDelete returns a nil error (successful deletion).
func (s *EmptyStore) GroupDelete(id string) error {
	return nil
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

// ProfileDelete returns a nil error (successful deletion).
func (s *EmptyStore) ProfileDelete(id string) error {
	return nil
}

// ProfileList returns an empty list of profiles.
func (s *EmptyStore) ProfileList() (profiles []*storagepb.Profile, err error) {
	return profiles, nil
}

// IgnitionPut returns an error writing any Ignition template.
func (s *EmptyStore) IgnitionPut(name string, config []byte) error {
	return fmt.Errorf("emptyStore does not accept Ignition templates")
}

// IgnitionGet get returns an Ignition template not found error.
func (s *EmptyStore) IgnitionGet(name string) (string, error) {
	return "", fmt.Errorf("no Ignition template %s", name)
}

// IgnitionDelete returns a nil error (successful deletion).
func (s *EmptyStore) IgnitionDelete(name string) error {
	return nil
}

// GenericPut returns an error writing any Generic template.
func (s *EmptyStore) GenericPut(name string, config []byte) error {
	return fmt.Errorf("emptyStore does not accept Generic templates")
}

// GenericGet get returns an Generic template not found error.
func (s *EmptyStore) GenericGet(name string) (string, error) {
	return "", fmt.Errorf("no Generic template %s", name)
}

// GenericDelete returns a nil error (successful deletion).
func (s *EmptyStore) GenericDelete(name string) error {
	return nil
}

// CloudGet returns a Cloud-config template not found error.
func (s *EmptyStore) CloudGet(name string) (string, error) {
	return "", fmt.Errorf("no Cloud-Config template %s", name)
}
