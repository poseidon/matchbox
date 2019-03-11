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

type Networkd struct {
	Units []NetworkdUnit `yaml:"units"`
}

type NetworkdUnit struct {
	Name     string               `yaml:"name"`
	Contents string               `yaml:"contents"`
	Dropins  []NetworkdUnitDropIn `yaml:"dropins"`
}

type NetworkdUnitDropIn struct {
	Name     string `yaml:"name"`
	Contents string `yaml:"contents"`
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		for _, unit := range in.Networkd.Units {
			newUnit := ignTypes.Networkdunit{
				Name:     unit.Name,
				Contents: unit.Contents,
			}
			for _, dropIn := range unit.Dropins {
				newUnit.Dropins = append(newUnit.Dropins, ignTypes.NetworkdDropin{
					Name:     dropIn.Name,
					Contents: dropIn.Contents,
				})
			}
			out.Networkd.Units = append(out.Networkd.Units, newUnit)
		}
		return out, report.Report{}, ast
	})
}
