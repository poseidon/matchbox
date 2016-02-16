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

package config

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/coreos/ignition/third_party/github.com/camlistore/camlistore/pkg/errorutil"
)

type Config struct {
	Version  int      `json:"ignitionVersion"    yaml:"ignition_version"`
	Storage  Storage  `json:"storage,omitempty"  yaml:"storage"`
	Systemd  Systemd  `json:"systemd,omitempty"  yaml:"systemd"`
	Networkd Networkd `json:"networkd,omitempty" yaml:"networkd"`
	Passwd   Passwd   `json:"passwd,omitempty"   yaml:"passwd"`
}

const (
	Version = 1
)

var (
	ErrVersion     = errors.New("incorrect config version")
	ErrCloudConfig = errors.New("not a config (found coreos-cloudconfig)")
	ErrScript      = errors.New("not a config (found coreos-cloudinit script)")
	ErrEmpty       = errors.New("not a config (empty)")
)

func Parse(config []byte) (cfg Config, err error) {
	if err = json.Unmarshal(config, &cfg); err == nil {
		if cfg.Version != Version {
			err = ErrVersion
		}
	} else if isCloudConfig(decompressIfGzipped(config)) {
		err = ErrCloudConfig
	} else if isScript(decompressIfGzipped(config)) {
		err = ErrScript
	} else if isEmpty(config) {
		err = ErrEmpty
	}
	if serr, ok := err.(*json.SyntaxError); ok {
		line, col, highlight := errorutil.HighlightBytePosition(bytes.NewReader(config), serr.Offset)
		err = fmt.Errorf("error at line %d, column %d\n%s%v", line, col, highlight, err)
	}

	return
}

func isEmpty(userdata []byte) bool {
	return len(userdata) == 0
}

func decompressIfGzipped(data []byte) []byte {
	if reader, err := gzip.NewReader(bytes.NewReader(data)); err == nil {
		uncompressedData, err := ioutil.ReadAll(reader)
		reader.Close()
		if err == nil {
			return uncompressedData
		} else {
			return data
		}
	} else {
		return data
	}
}
