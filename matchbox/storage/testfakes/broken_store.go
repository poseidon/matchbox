package testfakes

import (
	"errors"

	"github.com/coreos/coreos-baremetal/matchbox/storage/storagepb"
)

var (
	errIntentional = errors.New("store: error for testing purposes")
)

// BrokenStore returns errors for testing purposes.
type BrokenStore struct{}

// GroupPut returns an error.
func (s *BrokenStore) GroupPut(group *storagepb.Group) error {
	return errIntentional
}

// GroupGet returns an error.
func (s *BrokenStore) GroupGet(id string) (*storagepb.Group, error) {
	return nil, errIntentional
}

// GroupList returns an error.
func (s *BrokenStore) GroupList() (groups []*storagepb.Group, err error) {
	return groups, errIntentional
}

// ProfilePut returns an error.
func (s *BrokenStore) ProfilePut(profile *storagepb.Profile) error {
	return errIntentional
}

// ProfileGet returns an error.
func (s *BrokenStore) ProfileGet(id string) (*storagepb.Profile, error) {
	return nil, errIntentional
}

// ProfileList returns an error.
func (s *BrokenStore) ProfileList() (profiles []*storagepb.Profile, err error) {
	return profiles, errIntentional
}

// IgnitionPut returns an error.
func (s *BrokenStore) IgnitionPut(name string, config []byte) error {
	return errIntentional
}

// IgnitionGet returns an error.
func (s *BrokenStore) IgnitionGet(name string) (string, error) {
	return "", errIntentional
}

// CloudGet returns an error.
func (s *BrokenStore) CloudGet(name string) (string, error) {
	return "", errIntentional
}

// GenericGet returns an error.
func (s *BrokenStore) GenericGet(name string) (string, error) {
	return "", errIntentional
}
