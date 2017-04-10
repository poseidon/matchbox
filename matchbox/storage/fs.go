package storage

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	defaultDirectoryMode        os.FileMode = 0755
	defaultFileMode             os.FileMode = 0644
	errInvalidFilePathCharacter             = errors.New("invalid character in file path")
)

// Dir implements access to a collection of named files, restricted to a
// specific directory tree. It is very similar to net/http.Dir, but provides
// write access and some io/ioutil utilities.
// An empty directory is treated as ".".
type Dir string

// readFile reads data from a file at a given path, restricted to a specific
// directory tree.
func (d Dir) readFile(path string) ([]byte, error) {
	path, err := d.sanitize(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(path)
}

// readDir reads the directory named by the given path and returns a list of
// sorted directory entries. Restricted to a specified directory tree.
func (d Dir) readDir(dirname string) ([]os.FileInfo, error) {
	path, err := d.sanitize(dirname)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadDir(path)
}

// writeFile writes the data as a file at given path, restricted to a specific
// directory tree.
func (d Dir) writeFile(path string, data []byte) error {
	// make parent directories as needed
	path, err := d.sanitize(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), defaultDirectoryMode); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, defaultFileMode)
}

// deleteFile removes the file at the given path, restricted to a specific
// directory tree.
func (d Dir) deleteFile(path string) error {
	path, err := d.sanitize(path)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

// Borrowed directly from net/http Dir.Open and FileServer.
func (d Dir) sanitize(name string) (string, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
		strings.Contains(name, "\x00") {
		return "", errInvalidFilePathCharacter
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))), nil
}
