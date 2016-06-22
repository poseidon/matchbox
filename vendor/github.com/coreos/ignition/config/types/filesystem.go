// Copyright 2016 CoreOS, Inc.
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

package types

import (
	"encoding/json"
	"errors"
)

var (
	ErrFilesystemInvalidFormat = errors.New("invalid filesystem format")
	ErrFilesystemNoMountPath   = errors.New("filesystem is missing mount or path")
	ErrFilesystemMountAndPath  = errors.New("filesystem has both mount and path defined")
)

type Filesystem struct {
	Name  string           `json:"name,omitempty"`
	Mount *FilesystemMount `json:"mount,omitempty"`
	Path  *Path            `json:"path,omitempty"`
}
type filesystem Filesystem

type FilesystemMount struct {
	Device Path              `json:"device,omitempty"`
	Format FilesystemFormat  `json:"format,omitempty"`
	Create *FilesystemCreate `json:"create,omitempty"`
}

type FilesystemCreate struct {
	Force   bool        `json:"force,omitempty"`
	Options MkfsOptions `json:"options,omitempty"`
}

func (f *Filesystem) UnmarshalJSON(data []byte) error {
	tf := filesystem(*f)
	if err := json.Unmarshal(data, &tf); err != nil {
		return err
	}
	*f = Filesystem(tf)
	return f.AssertValid()
}

func (f Filesystem) AssertValid() error {
	hasMount := false
	hasPath := false

	if f.Mount != nil {
		hasMount = true
		if err := f.Mount.AssertValid(); err != nil {
			return err
		}
	}

	if f.Path != nil {
		hasPath = true
		if err := f.Path.AssertValid(); err != nil {
			return err
		}
	}

	if !hasMount && !hasPath {
		return ErrFilesystemNoMountPath
	} else if hasMount && hasPath {
		return ErrFilesystemMountAndPath
	}

	return nil
}

type filesystemMount FilesystemMount

func (f *FilesystemMount) UnmarshalJSON(data []byte) error {
	tf := filesystemMount(*f)
	if err := json.Unmarshal(data, &tf); err != nil {
		return err
	}
	*f = FilesystemMount(tf)
	return f.AssertValid()
}

func (f FilesystemMount) AssertValid() error {
	if err := f.Device.AssertValid(); err != nil {
		return err
	}
	if err := f.Format.AssertValid(); err != nil {
		return err
	}
	return nil
}

type FilesystemFormat string
type filesystemFormat FilesystemFormat

func (f *FilesystemFormat) UnmarshalJSON(data []byte) error {
	tf := filesystemFormat(*f)
	if err := json.Unmarshal(data, &tf); err != nil {
		return err
	}
	*f = FilesystemFormat(tf)
	return f.AssertValid()
}

func (f FilesystemFormat) AssertValid() error {
	switch f {
	case "ext4", "btrfs", "xfs":
		return nil
	default:
		return ErrFilesystemInvalidFormat
	}
}

type MkfsOptions []string
type mkfsOptions MkfsOptions

func (o *MkfsOptions) UnmarshalJSON(data []byte) error {
	to := mkfsOptions(*o)
	if err := json.Unmarshal(data, &to); err != nil {
		return err
	}
	*o = MkfsOptions(to)
	return o.AssertValid()
}

func (o MkfsOptions) AssertValid() error {
	return nil
}
