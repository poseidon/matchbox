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
	"path/filepath"
)

var (
	ErrPathRelative = errors.New("path not absolute")
)

type Path string

func (d *Path) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return d.unmarshal(unmarshal)
}

func (d *Path) UnmarshalJSON(data []byte) error {
	return d.unmarshal(func(td interface{}) error {
		return json.Unmarshal(data, td)
	})
}

type path Path

func (d *Path) unmarshal(unmarshal func(interface{}) error) error {
	td := path(*d)
	if err := unmarshal(&td); err != nil {
		return err
	}
	*d = Path(td)
	return d.assertValid()
}

func (d Path) assertValid() error {
	if !filepath.IsAbs(string(d)) {
		return ErrPathRelative
	}
	return nil
}
