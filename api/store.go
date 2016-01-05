package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	ignition "github.com/coreos/ignition/src/config"
)

// Store provides Machine, Spec, and config resources.
type Store interface {
	Machine(id string) (*Machine, error)
	Spec(id string) (*Spec, error)
	CloudConfig(id string) (*CloudConfig, error)
	IgnitionConfig(id string) (*ignition.Config, error)
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

// Machine returns the configuration for the machine with the given id.
func (s *fileStore) Machine(id string) (*Machine, error) {
	file, err := openFile(s.root, filepath.Join("machines", id, "machine.json"))
	if err != nil {
		log.Debugf("no machine config %s", id)
		return nil, err
	}
	defer file.Close()

	machine := new(Machine)
	err = json.NewDecoder(file).Decode(machine)
	if err != nil {
		log.Errorf("error decoding machine config: %s", err)
		return nil, err
	}

	if machine.Spec == nil && machine.SpecID != "" {
		// machine references a Spec, attempt to add Spec properties
		spec, err := s.Spec(machine.SpecID)
		if err == nil {
			machine.Spec = spec
		}
	}
	return machine, err
}

// Spec returns the Spec with the given id.
func (s *fileStore) Spec(id string) (*Spec, error) {
	file, err := openFile(s.root, filepath.Join("specs", id, "spec.json"))
	if err != nil {
		log.Debugf("no spec %s", id)
		return nil, err
	}
	defer file.Close()

	spec := new(Spec)
	err = json.NewDecoder(file).Decode(spec)
	if err != nil {
		log.Errorf("error decoding spec: %s", err)
		return nil, err
	}
	return spec, err
}

// CloudConfig returns the cloud config with the given id.
func (s *fileStore) CloudConfig(id string) (*CloudConfig, error) {
	file, err := openFile(s.root, filepath.Join("cloud", id))
	if err != nil {
		log.Debugf("no cloud config %s", id)
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("error reading cloud config: %s", err)
		return nil, err
	}
	return &CloudConfig{
		Content: string(b),
	}, nil
}

// IgnitionConfig returns the ignition config with the given id.
func (s *fileStore) IgnitionConfig(id string) (*ignition.Config, error) {
	file, err := openFile(s.root, filepath.Join("ignition", id))
	if err != nil {
		log.Debugf("no ignition config %s", id)
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("error reading ignition config: %s", err)
		return nil, err
	}
	config, err := ignition.Parse(b)
	if err != nil {
		log.Errorf("error parsing ignition config: %s", err)
		return nil, err
	}
	return &config, err
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
