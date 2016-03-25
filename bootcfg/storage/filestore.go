package storage

import (
	"encoding/json"
	"path/filepath"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

// Config initializes a fileStore.
type Config struct {
	Dir    string
	Groups []*storagepb.Group
}

// fileStore implements ths Store interface. Queries to the file system
// are restricted to the specified directory tree.
type fileStore struct {
	root   string
	groups map[string]*storagepb.Group
}

// NewFileStore returns a new memory-backed Store.
func NewFileStore(config *Config) Store {
	groups := make(map[string]*storagepb.Group)
	for _, group := range config.Groups {
		groups[group.Id] = group
	}
	return &fileStore{
		root:   config.Dir,
		groups: groups,
	}
}

// GroupGet returns a machine Group by id.
func (s *fileStore) GroupGet(id string) (*storagepb.Group, error) {
	val, ok := s.groups[id]
	if !ok {
		return nil, ErrGroupNotFound
	}
	return val, nil
}

// GroupList lists all machine Groups.
func (s *fileStore) GroupList() ([]*storagepb.Group, error) {
	groups := make([]*storagepb.Group, len(s.groups))
	i := 0
	for _, g := range s.groups {
		groups[i] = g
		i++
	}
	return groups, nil
}

// ProfileGet gets a profile by id.
func (s *fileStore) ProfileGet(id string) (*storagepb.Profile, error) {
	data, err := Dir(s.root).readFile(filepath.Join("profiles", id, "profile.json"))
	if err != nil {
		return nil, err
	}
	profile := new(storagepb.Profile)
	err = json.Unmarshal(data, profile)
	if err != nil {
		return nil, err
	}
	return profile, err
}

// ProfileList lists all profiles.
func (s *fileStore) ProfileList() ([]*storagepb.Profile, error) {
	files, err := Dir(s.root).readDir("profiles")
	if err != nil {
		return nil, err
	}
	profiles := make([]*storagepb.Profile, 0, len(files))
	for _, finfo := range files {
		profile, err := s.ProfileGet(finfo.Name())
		if err == nil {
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}

// IgnitionGet gets an Ignition Config template by name.
func (s *fileStore) IgnitionGet(id string) (string, error) {
	data, err := Dir(s.root).readFile(filepath.Join("ignition", id))
	return string(data), err
}

// CloudGet gets a Cloud-Config template by name.
func (s *fileStore) CloudGet(id string) (string, error) {
	data, err := Dir(s.root).readFile(filepath.Join("cloud", id))
	return string(data), err
}
