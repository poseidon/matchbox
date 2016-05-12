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

package files

import (
	"reflect"
	"testing"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/exec/util"
	"github.com/coreos/ignition/internal/log"
)

func TestMapFilesToFilesystems(t *testing.T) {
	type in struct {
		config types.Config
	}
	type out struct {
		files map[types.Filesystem][]types.File
		err   error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{config: types.Config{}},
			out: out{files: map[types.Filesystem][]types.File{}},
		},
		{
			in:  in{config: types.Config{Storage: types.Storage{Files: []types.File{{Filesystem: "foo"}}}}},
			out: out{err: ErrFilesystemUndefined},
		},
		{
			in: in{config: types.Config{Storage: types.Storage{
				Filesystems: []types.Filesystem{{Name: "fs1"}},
				Files:       []types.File{{Filesystem: "fs1", Path: "/foo"}, {Filesystem: "fs1", Path: "/bar"}},
			}}},
			out: out{files: map[types.Filesystem][]types.File{types.Filesystem{Name: "fs1"}: {{Filesystem: "fs1", Path: "/foo"}, {Filesystem: "fs1", Path: "/bar"}}}},
		},
		{
			in: in{config: types.Config{Storage: types.Storage{
				Filesystems: []types.Filesystem{{Name: "fs1", Path: "/fs1"}, {Name: "fs2", Path: "/fs2"}},
				Files:       []types.File{{Filesystem: "fs1", Path: "/foo"}, {Filesystem: "fs2", Path: "/bar"}},
			}}},
			out: out{files: map[types.Filesystem][]types.File{
				types.Filesystem{Name: "fs1", Path: "/fs1"}: {{Filesystem: "fs1", Path: "/foo"}},
				types.Filesystem{Name: "fs2", Path: "/fs2"}: {{Filesystem: "fs2", Path: "/bar"}},
			}},
		},
		{
			in: in{config: types.Config{Storage: types.Storage{
				Filesystems: []types.Filesystem{{Name: "fs1"}, {Name: "fs1", Path: "/fs1"}},
				Files:       []types.File{{Filesystem: "fs1", Path: "/foo"}, {Filesystem: "fs1", Path: "/bar"}},
			}}},
			out: out{files: map[types.Filesystem][]types.File{
				types.Filesystem{Name: "fs1", Path: "/fs1"}: {{Filesystem: "fs1", Path: "/foo"}, {Filesystem: "fs1", Path: "/bar"}},
			}},
		},
	}

	for i, test := range tests {
		logger := log.New()
		files, err := stage{Util: util.Util{Logger: &logger}}.mapFilesToFilesystems(test.in.config)
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
		if !reflect.DeepEqual(test.out.files, files) {
			t.Errorf("#%d: bad map: want %#v, got %#v", i, test.out.files, files)
		}
	}
}
