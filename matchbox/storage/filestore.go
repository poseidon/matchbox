package storage

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

// Config initializes a fileStore.
type Config struct {
	Root   string
	Logger *logrus.Logger
}

// fileStore implements ths Store interface. Queries to the file system
// are restricted to the specified directory tree.
type fileStore struct {
	root   string
	logger *logrus.Logger
}

// NewFileStore returns a new memory-backed Store.
func NewFileStore(config *Config) Store {
	return &fileStore{
		root:   config.Root,
		logger: config.Logger,
	}
}

// GroupPut writes the given Group.
func (s *fileStore) GroupPut(group *storagepb.Group) error {
	richGroup, err := group.ToRichGroup()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(richGroup, "", "\t")
	if err != nil {
		return err
	}
	return Dir(s.root).writeFile(filepath.Join("groups", group.Id+".json"), data)
}

// GroupGet returns a machine Group by id.
func (s *fileStore) GroupGet(id string) (*storagepb.Group, error) {
	data, err := Dir(s.root).readFile(filepath.Join("groups", id+".json"))
	if err != nil {
		return nil, err
	}
	group, err := storagepb.ParseGroup(data)
	if err != nil {
		return nil, err
	}
	return group, err
}

// GroupDelete deletes a machine Group by id.
func (s *fileStore) GroupDelete(id string) error {
	return Dir(s.root).deleteFile(filepath.Join("groups", id+".json"))
}

// GroupList lists all machine Groups.
func (s *fileStore) GroupList() ([]*storagepb.Group, error) {
	files, err := Dir(s.root).readDir("groups")
	if err != nil {
		return nil, err
	}
	groups := make([]*storagepb.Group, 0, len(files))
	for _, finfo := range files {
		name := strings.TrimSuffix(finfo.Name(), filepath.Ext(finfo.Name()))
		group, err := s.GroupGet(name)
		if err == nil {
			groups = append(groups, group)
		} else if s.logger != nil {
			s.logger.Infof("Group %q: %v", name, err)
		}
	}
	return groups, nil
}

// ProfilePut writes the given Profile.
func (s *fileStore) ProfilePut(profile *storagepb.Profile) error {
	data, err := json.MarshalIndent(profile, "", "\t")
	if err != nil {
		return err
	}
	return Dir(s.root).writeFile(filepath.Join("profiles", profile.Id+".json"), data)
}

// ProfileGet gets a profile by id.
func (s *fileStore) ProfileGet(id string) (*storagepb.Profile, error) {
	data, err := Dir(s.root).readFile(filepath.Join("profiles", id+".json"))
	if err != nil {
		return nil, err
	}
	profile := new(storagepb.Profile)
	err = json.Unmarshal(data, profile)
	if err != nil {
		return nil, err
	}
	if err := profile.AssertValid(); err != nil {
		return nil, err
	}
	return profile, err
}

// ProfileDelete deletes a profile by id.
func (s *fileStore) ProfileDelete(id string) error {
	return Dir(s.root).deleteFile(filepath.Join("profiles", id+".json"))
}

// ProfileList lists all profiles.
func (s *fileStore) ProfileList() ([]*storagepb.Profile, error) {
	files, err := Dir(s.root).readDir("profiles")
	if err != nil {
		return nil, err
	}
	profiles := make([]*storagepb.Profile, 0, len(files))
	for _, finfo := range files {
		name := strings.TrimSuffix(finfo.Name(), filepath.Ext(finfo.Name()))
		profile, err := s.ProfileGet(name)
		if err == nil {
			profiles = append(profiles, profile)
		} else if s.logger != nil {
			s.logger.Infof("Profile %q: %v", name, err)
		}
	}
	return profiles, nil
}

// IgnitionPut creates or updates an Ignition template.
func (s *fileStore) IgnitionPut(name string, config []byte) error {
	return Dir(s.root).writeFile(filepath.Join("ignition", name), config)
}

// IgnitionGet gets an Ignition template by name.
func (s *fileStore) IgnitionGet(name string) (string, error) {
	data, err := Dir(s.root).readFile(filepath.Join("ignition", name))
	return string(data), err
}

// IgnitionDelete deletes an Ignition template by name.
func (s *fileStore) IgnitionDelete(name string) error {
	return Dir(s.root).deleteFile(filepath.Join("ignition", name))
}

// CloudGet gets a Cloud-Config template by name.
func (s *fileStore) CloudGet(name string) (string, error) {
	data, err := Dir(s.root).readFile(filepath.Join("cloud", name))
	return string(data), err
}

// GenericGet gets a generic template by name.
func (s *fileStore) GenericGet(name string) (string, error) {
	data, err := Dir(s.root).readFile(filepath.Join("generic", name))
	return string(data), err
}
