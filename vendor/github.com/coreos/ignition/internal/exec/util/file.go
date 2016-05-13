// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/util"

	"github.com/vincent-petithory/dataurl"
)

const (
	DefaultDirectoryPermissions os.FileMode = 0755
	DefaultFilePermissions      os.FileMode = 0644
)

var (
	ErrSchemeUnsupported = errors.New("unsupported source scheme")
	ErrStatusBad         = errors.New("bad HTTP response status")
)

type File struct {
	Path     types.Path
	Contents []byte
	Mode     os.FileMode
	Uid      int
	Gid      int
}

func RenderFile(l *log.Logger, f types.File) *File {
	var contents []byte
	var err error

	fetch := func() error {
		contents, err = fetchFile(l, f)
		return err
	}

	validate := func() error {
		return util.AssertValid(f.Contents.Verification, contents)
	}

	decompress := func() error {
		contents, err = decompressFile(l, f, contents)
		return err
	}

	if l.LogOp(fetch, "fetching file %q", f.Path) != nil {
		return nil
	}
	if l.LogOp(validate, "validating file contents") != nil {
		return nil
	}
	if l.LogOp(decompress, "decompressing file contents") != nil {
		return nil
	}

	return &File{
		Path:     f.Path,
		Contents: []byte(contents),
		Mode:     os.FileMode(f.Mode),
		Uid:      f.User.Id,
		Gid:      f.Group.Id,
	}
}

func fetchFile(l *log.Logger, f types.File) ([]byte, error) {
	switch f.Contents.Source.Scheme {
	case "http":
		client := util.NewHttpClient(l)
		data, status, err := client.Get(f.Contents.Source.String())
		if err != nil {
			return nil, err
		}

		l.Debug("GET result: %s", http.StatusText(status))
		if status != http.StatusOK {
			return nil, ErrStatusBad
		}

		return data, nil
	case "data":
		url, err := dataurl.DecodeString(f.Contents.Source.String())
		if err != nil {
			return nil, err
		}

		return url.Data, nil
	default:
		return nil, ErrSchemeUnsupported
	}
}

func decompressFile(l *log.Logger, f types.File, contents []byte) ([]byte, error) {
	switch f.Contents.Compression {
	case "":
		return contents, nil
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewReader(contents))
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		return ioutil.ReadAll(reader)
	default:
		return nil, types.ErrCompressionInvalid
	}
}

// WriteFile creates and writes the file described by f using the provided context
func (u Util) WriteFile(f *File) error {
	var err error

	path := u.JoinPath(string(f.Path))

	if err := mkdirForFile(path); err != nil {
		return err
	}

	// Create a temporary file in the same directory to ensure it's on the same filesystem
	var tmp *os.File
	if tmp, err = ioutil.TempFile(filepath.Dir(path), "tmp"); err != nil {
		return err
	}
	tmp.Close()
	defer func() {
		if err != nil {
			os.Remove(tmp.Name())
		}
	}()

	if err := ioutil.WriteFile(tmp.Name(), f.Contents, f.Mode); err != nil {
		return err
	}

	// XXX(vc): Note that we assume to be operating on the file we just wrote, this is only guaranteed
	// by using syscall.Fchown() and syscall.Fchmod()

	// Ensure the ownership and mode are as requested (since WriteFile can be affected by sticky bit)
	if err := os.Chown(tmp.Name(), f.Uid, f.Gid); err != nil {
		return err
	}

	if err := os.Chmod(tmp.Name(), f.Mode); err != nil {
		return err
	}

	if err := os.Rename(tmp.Name(), path); err != nil {
		return err
	}

	return nil
}

// mkdirForFile helper creates the directory components of path
func mkdirForFile(path string) error {
	return os.MkdirAll(filepath.Dir(path), DefaultDirectoryPermissions)
}
