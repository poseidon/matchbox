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

	"github.com/coreos/ignition/config/types"
	"github.com/go-yaml/yaml"
)

func ParseAsV2_0_0(data []byte) (types.Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return types.Config{}, err
	}

	var keyMap map[interface{}]interface{}
	if err := yaml.Unmarshal(data, &keyMap); err != nil {
		return types.Config{}, err
	}

	if err := assertKeysValid(keyMap, reflect.TypeOf(Config{})); err != nil {
		return types.Config{}, err
	}

	return ConvertAs2_0_0(cfg)
}
