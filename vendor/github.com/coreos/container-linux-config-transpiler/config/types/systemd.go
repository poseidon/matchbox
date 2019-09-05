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

type Systemd struct {
	Units []SystemdUnit `yaml:"units"`
}

type SystemdUnit struct {
	Name     string              `yaml:"name"`
	Enable   bool                `yaml:"enable"`
	Enabled  *bool               `yaml:"enabled"`
	Mask     bool                `yaml:"mask"`
	Contents string              `yaml:"contents"`
	Dropins  []SystemdUnitDropIn `yaml:"dropins"`
}

type SystemdUnitDropIn struct {
	Name     string `yaml:"name"`
	Contents string `yaml:"contents"`
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		for _, unit := range in.Systemd.Units {
			newUnit := ignTypes.Unit{
				Name:     unit.Name,
				Enable:   unit.Enable,
				Enabled:  unit.Enabled,
				Mask:     unit.Mask,
				Contents: unit.Contents,
			}

			for _, dropIn := range unit.Dropins {
				newUnit.Dropins = append(newUnit.Dropins, ignTypes.SystemdDropin{
					Name:     dropIn.Name,
					Contents: dropIn.Contents,
				})
			}

			out.Systemd.Units = append(out.Systemd.Units, newUnit)
		}
		return out, report.Report{}, ast
	})
}
