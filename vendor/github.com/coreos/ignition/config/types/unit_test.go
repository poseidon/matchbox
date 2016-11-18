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
	"errors"
	"reflect"
	"testing"

	"github.com/coreos/ignition/config/validate/report"
)

func TestSystemdUnitNameValidate(t *testing.T) {
	type in struct {
		unit SystemdUnitName
	}
	type out struct {
		err error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{unit: SystemdUnitName("test.service")},
			out: out{err: nil},
		},
		{
			in:  in{unit: SystemdUnitName("test.socket")},
			out: out{err: nil},
		},
		{
			in:  in{unit: SystemdUnitName("test.blah")},
			out: out{err: errors.New("invalid systemd unit extension")},
		},
	}

	for i, test := range tests {
		err := test.in.unit.Validate()
		if !reflect.DeepEqual(report.ReportFromError(test.out.err, report.EntryError), err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
	}
}

func TestNetworkdUnitNameValidate(t *testing.T) {
	type in struct {
		unit NetworkdUnitName
	}
	type out struct {
		err error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{unit: NetworkdUnitName("test.network")},
			out: out{err: nil},
		},
		{
			in:  in{unit: NetworkdUnitName("test.link")},
			out: out{err: nil},
		},
		{
			in:  in{unit: NetworkdUnitName("test.netdev")},
			out: out{err: nil},
		},
		{
			in:  in{unit: NetworkdUnitName("test.blah")},
			out: out{err: errors.New("invalid networkd unit extension")},
		},
	}

	for i, test := range tests {
		err := test.in.unit.Validate()
		if !reflect.DeepEqual(report.ReportFromError(test.out.err, report.EntryError), err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
	}
}
