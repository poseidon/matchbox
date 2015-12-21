package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// Store maintains associations between machine attributes and configs.
type Store interface {
	// BootConfig returns the boot config (kernel, options) for the machine.
	BootConfig(attrs MachineAttrs) (*BootConfig, error)
	// CloudConfig returns the cloud config user data for the machine.
	CloudConfig(attrs MachineAttrs) (*CloudConfig, error)
}

// fileStore maps machine attributes to configs based on an http.Filesystem.
type fileStore struct {
	root http.FileSystem
}

// NewFileStore returns a Store backed by a filesystem directory.
func NewFileStore(root http.FileSystem) Store {
	return &fileStore{
		root: root,
	}
}

const (
	bootPrefix  = "boot"
	cloudPrefix = "cloud"
)

// BootConfig returns the boot config (kernel, options) for the machine.
func (s *fileStore) BootConfig(attrs MachineAttrs) (*BootConfig, error) {
	file, err := s.find(bootPrefix, attrs)
	if err != nil {
		log.Infof("no boot config for machine %+v", attrs)
		return nil, err
	}
	defer file.Close()

	config := new(BootConfig)
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		log.Errorf("error decoding boot config: %s", err)
	}
	return config, err
}

// CloudConfig returns the cloud config for the machine.
func (s *fileStore) CloudConfig(attrs MachineAttrs) (*CloudConfig, error) {
	file, err := s.find(cloudPrefix, attrs)
	if err != nil {
		log.Infof("no cloud config for machine %+v", attrs)
		return nil, err
	}
	defer file.Close()
	// cloudinit requires reading the entire file
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &CloudConfig{
		Content: string(b),
	}, nil
}

// find searches the prefix subdirectory of root for the first config file
// which matches the given machine attributes. If the error is non-nil, the
// caller must be sure to close the matched http.File. Matches are searched
// in priority order: uuid/<UUID>, mac/<MAC aaddress>, default.
func (s *fileStore) find(prefix string, attrs MachineAttrs) (http.File, error) {
	search := []string{
		filepath.Join("uuid", attrs.UUID),
		filepath.Join("mac", attrs.MAC.String()),
		"/default",
	}
	for _, path := range filter(search) {
		fullPath := filepath.Join(prefix, path)
		if file, err := openFile(s.root, fullPath); err == nil {
			return file, err
		}
	}
	return nil, fmt.Errorf("no %s config for machine %+v", prefix, attrs)
}

// filter returns only paths which have non-empty directory paths. For example,
// "uuid/123" has a directory path "uuid", while path "uuid" does not.
func filter(inputs []string) (paths []string) {
	for _, path := range inputs {
		if filepath.Dir(path) != "." {
			paths = append(paths, path)
		}
	}
	return paths
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
