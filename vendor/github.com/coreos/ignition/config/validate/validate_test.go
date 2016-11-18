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

package validate

import (
	"errors"
	"reflect"
	"testing"

	// Import into the same namespace to keep config definitions clean
	. "github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/config/validate/report"
)

func TestValidate(t *testing.T) {
	type in struct {
		cfg Config
	}
	type out struct {
		err error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{cfg: Config{Ignition: Ignition{Version: IgnitionVersion{Major: 2}}}},
			out: out{},
		},
		{
			in:  in{cfg: Config{}},
			out: out{err: ErrOldVersion},
		},
		{
			in: in{cfg: Config{
				Ignition: Ignition{
					Version: IgnitionVersion{Major: 2},
					Config: IgnitionConfig{
						Replace: &ConfigReference{
							Verification: Verification{
								Hash: &Hash{Function: "foobar"},
							},
						},
					},
				},
			}},
			out: out{errors.New("unrecognized hash function")},
		},
		{
			in: in{cfg: Config{
				Ignition: Ignition{Version: IgnitionVersion{Major: 2}},
				Storage: Storage{
					Filesystems: []Filesystem{
						{
							Name: "filesystem1",
							Mount: &FilesystemMount{
								Device: Path("/dev/disk/by-partlabel/ROOT"),
								Format: FilesystemFormat("btrfs"),
							},
						},
					},
				},
			}},
			out: out{},
		},
		{
			in: in{cfg: Config{
				Ignition: Ignition{Version: IgnitionVersion{Major: 2}},
				Storage: Storage{
					Filesystems: []Filesystem{
						{
							Name: "filesystem1",
							Path: func(p Path) *Path { return &p }("/sysroot"),
						},
					},
				},
			}},
			out: out{},
		},
		{
			in: in{cfg: Config{
				Ignition: Ignition{Version: IgnitionVersion{Major: 2}},
				Systemd:  Systemd{Units: []SystemdUnit{{Name: "foo.bar"}}},
			}},
			out: out{err: errors.New("invalid systemd unit extension")},
		},
	}

	for i, test := range tests {
		r := ValidateWithoutSource(reflect.ValueOf(test.in.cfg))
		expectedReport := report.ReportFromError(test.out.err, report.EntryError)
		if !reflect.DeepEqual(expectedReport, r) {
			t.Errorf("#%d: bad error: want %v, got %v", i, expectedReport, r)
		}
	}
}
