package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

// Errors querying a Store.
var (
	ErrGroupNotFound   = errors.New("storage: No Group found")
	ErrProfileNotFound = errors.New("storage: No Profile found")
)

// A Store stores machine Groups and Profiles.
type Store interface {
	// GroupGet returns a machine Group by id.
	GroupGet(id string) (*storagepb.Group, error)
	// GroupList lists all machine Groups.
	GroupList() ([]*storagepb.Group, error)
	// ProfileGet gets a profile by id.
	ProfileGet(id string) (*storagepb.Profile, error)
	// ProfileList lists all profiles.
	ProfileList() ([]*storagepb.Profile, error)
	// IgnitionGet gets an Ignition Config template by name.
	IgnitionGet(name string) (string, error)
	// CloudGet gets a Cloud-Config template by name.
	CloudGet(name string) (string, error)
}

// Config initializes a fileStore.
type Config struct {
	Dir    string
	Groups []*storagepb.Group
}

// fileStore implements ths Store interface. Queries to the file system
// are restricted to the specified directory tree.
type fileStore struct {
	dir    string
	groups map[string]*storagepb.Group
}

// NewFileStore returns a new memory-backed Store.
func NewFileStore(config *Config) Store {
	groups := make(map[string]*storagepb.Group)
	for _, group := range config.Groups {
		groups[group.Id] = group
	}
	return &fileStore{
		dir:    config.Dir,
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
	file, err := openFile(http.Dir(s.dir), filepath.Join("profiles", id, "profile.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	profile := new(storagepb.Profile)
	err = json.NewDecoder(file).Decode(profile)
	if err != nil {
		return nil, err
	}
	return profile, err
}

// ProfileList lists all profiles.
func (s *fileStore) ProfileList() ([]*storagepb.Profile, error) {
	finfos, err := ioutil.ReadDir(filepath.Join(s.dir, "profiles"))
	if err != nil {
		return nil, err
	}
	profiles := make([]*storagepb.Profile, 0, len(finfos))
	for _, finfo := range finfos {
		profile, err := s.ProfileGet(finfo.Name())
		if err == nil {
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}

// IgnitionGet gets an Ignition Config template by name.
func (s *fileStore) IgnitionGet(id string) (string, error) {
	file, err := openFile(http.Dir(s.dir), filepath.Join("ignition", id))
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// CloudGet gets a Cloud-Config template by name.
func (s *fileStore) CloudGet(id string) (string, error) {
	file, err := openFile(http.Dir(s.dir), filepath.Join("cloud", id))
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// openFile attempts to open the file within the specified Filesystem. If
// successful, the http.File is returned and must be closed by the caller.
// Otherwise, the path was not a regular file that could be opened and an
// error is returned.
func openFile(fs http.FileSystem, path string) (http.File, error) {
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	if info.Mode().IsRegular() {
		return file, nil
	}
	file.Close()
	return nil, fmt.Errorf("%s is not a file on the given filesystem", path)
}
