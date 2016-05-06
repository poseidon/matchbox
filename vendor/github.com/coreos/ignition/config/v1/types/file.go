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
	"os"
)

var (
	ErrFileIllegalMode = errors.New("illegal file mode")
)

type FileMode os.FileMode

type File struct {
	Path     Path     `json:"path,omitempty"     yaml:"path"`
	Contents string   `json:"contents,omitempty" yaml:"contents"`
	Mode     FileMode `json:"mode,omitempty"     yaml:"mode"`
	Uid      int      `json:"uid,omitempty"      yaml:"uid"`
	Gid      int      `json:"gid,omitempty"      yaml:"gid"`
}

func (m *FileMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return m.unmarshal(unmarshal)
}

func (m *FileMode) UnmarshalJSON(data []byte) error {
	return m.unmarshal(func(tm interface{}) error {
		return json.Unmarshal(data, tm)
	})
}

type fileMode FileMode

func (m *FileMode) unmarshal(unmarshal func(interface{}) error) error {
	tm := fileMode(*m)
	if err := unmarshal(&tm); err != nil {
		return err
	}
	*m = FileMode(tm)
	return m.assertValid()
}

func (m FileMode) assertValid() error {
	if (m &^ 07777) != 0 {
		return ErrFileIllegalMode
	}
	return nil
}
