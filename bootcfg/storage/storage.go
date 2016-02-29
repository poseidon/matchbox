package storage

import (
	"errors"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

// Errors querying a Store.
var (
	ErrGroupNotFound = errors.New("storage: No Group found")
)

// A Store provides machine Groups.
type Store interface {
	// Get a machine Group by id.
	GetGroup(id string) (*storagepb.Group, error)
	// List all machine Groups.
	ListGroups() ([]*storagepb.Group, error)
}

// Config initializes a memStore.
type Config struct {
	Groups []*storagepb.Group
}

// memStore implements ths Store interface.
type memStore struct {
	groups map[string] *storagepb.Group
}

// NewMemStore returns a new memory-backed Store.
func NewMemStore(config *Config) Store {
	groups := make(map[string]*storagepb.Group)
	for _, group := range config.Groups {
		groups[group.Id] = group
	}
	return &memStore{
		groups: groups,
	}
}

// GetGroup returns a machine Group by id.
func (s *memStore) GetGroup(id string) (*storagepb.Group, error) {
	val, ok := s.groups[id]
	if !ok {
		return nil, ErrGroupNotFound
	}
	return val, nil
}

// ListGroups lists all machine Groups.
func (s *memStore) ListGroups() ([]*storagepb.Group, error) {
	groups := make([]*storagepb.Group, len(s.groups))
	i := 0
	for _, g := range s.groups {
		groups[i] = g
		i++
	}
	return groups, nil
}
