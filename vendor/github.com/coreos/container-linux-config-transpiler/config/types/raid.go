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
	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/coreos/ignition/config/validate/astnode"
	"github.com/coreos/ignition/config/validate/report"
)

type Raid struct {
	Name    string   `yaml:"name"`
	Level   string   `yaml:"level"`
	Devices []string `yaml:"devices"`
	Spares  int      `yaml:"spares"`
	Options []string `yaml:"options"`
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		for _, array := range in.Storage.Arrays {
			newArray := ignTypes.Raid{
				Name:    array.Name,
				Level:   array.Level,
				Spares:  array.Spares,
				Devices: convertStringSliceToTypesDeviceSlice(array.Devices),
				Options: convertStringSiceToTypesRaidOptionSlice(array.Options),
			}

			out.Storage.Raid = append(out.Storage.Raid, newArray)
		}
		return out, report.Report{}, ast
	})
}

// golang--
func convertStringSliceToTypesDeviceSlice(ss []string) []ignTypes.Device {
	var res []ignTypes.Device
	for _, s := range ss {
		res = append(res, ignTypes.Device(s))
	}
	return res
}

// golang--
func convertStringSiceToTypesRaidOptionSlice(ss []string) []ignTypes.RaidOption {
	var res []ignTypes.RaidOption
	for _, s := range ss {
		res = append(res, ignTypes.RaidOption(s))
	}
	return res
}
