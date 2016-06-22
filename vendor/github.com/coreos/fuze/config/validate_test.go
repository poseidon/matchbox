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

package config

import (
	"reflect"
	"testing"

	"github.com/go-yaml/yaml"
)

func TestAssertKeysValid(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		err ErrKeysUnrecognized
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: "ignition:\n  config:"},
			out: out{},
		},
		{
			in:  in{data: "passwd:\n  groups:\n    - name: example"},
			out: out{},
		},
		{
			in:  in{data: "password:\n  groups:"},
			out: out{err: ErrKeysUnrecognized{"password"}},
		},
		{
			in:  in{data: "passwd:\n  groups:\n    - naem: example"},
			out: out{err: ErrKeysUnrecognized{"naem"}},
		},
	}

	for i, test := range tests {
		var cfg interface{}
		if err := yaml.Unmarshal([]byte(test.in.data), &cfg); err != nil {
			t.Errorf("#%d: unmarshal failed: %v", i, err)
			continue
		}
		if err := assertKeysValid(cfg, reflect.TypeOf(Config{})); !reflect.DeepEqual(err, test.out.err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
	}
}
