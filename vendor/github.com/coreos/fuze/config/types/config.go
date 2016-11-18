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

type Config struct {
	Ignition Ignition `yaml:"ignition"`
	Storage  Storage  `yaml:"storage"`
	Systemd  Systemd  `yaml:"systemd"`
	Networkd Networkd `yaml:"networkd"`
	Passwd   Passwd   `yaml:"passwd"`
}

type Ignition struct {
	Config IgnitionConfig `yaml:"config"`
}

type IgnitionConfig struct {
	Append  []ConfigReference `yaml:"append"`
	Replace *ConfigReference  `yaml:"replace"`
}

type ConfigReference struct {
	Source       string       `yaml:"source"`
	Verification Verification `yaml:"verification"`
}
