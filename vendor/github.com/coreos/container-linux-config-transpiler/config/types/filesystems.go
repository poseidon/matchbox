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

type Filesystem struct {
	Name  string  `yaml:"name"`
	Mount *Mount  `yaml:"mount"`
	Path  *string `yaml:"path"`
}

type Mount struct {
	Device         string   `yaml:"device"`
	Format         string   `yaml:"format"`
	Create         *Create  `yaml:"create"`
	WipeFilesystem bool     `yaml:"wipe_filesystem"`
	Label          *string  `yaml:"label"`
	UUID           *string  `yaml:"uuid"`
	Options        []string `yaml:"options"`
}

type Create struct {
	Force   bool     `yaml:"force"`
	Options []string `yaml:"options"`
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		r := report.Report{}
		for _, filesystem := range in.Storage.Filesystems {
			newFilesystem := ignTypes.Filesystem{
				Name: filesystem.Name,
				Path: filesystem.Path,
			}

			if filesystem.Mount != nil {
				newFilesystem.Mount = &ignTypes.Mount{
					Device:         filesystem.Mount.Device,
					Format:         filesystem.Mount.Format,
					WipeFilesystem: filesystem.Mount.WipeFilesystem,
					Label:          filesystem.Mount.Label,
					UUID:           filesystem.Mount.UUID,
					Options:        convertStringSliceToTypesMountOptionSlice(filesystem.Mount.Options),
				}

				if filesystem.Mount.Create != nil {
					newFilesystem.Mount.Create = &ignTypes.Create{
						Force:   filesystem.Mount.Create.Force,
						Options: convertStringSliceToTypesCreateOptionSlice(filesystem.Mount.Create.Options),
					}
				}
			}

			out.Storage.Filesystems = append(out.Storage.Filesystems, newFilesystem)
		}
		return out, r, ast
	})
}

// golang--
func convertStringSliceToTypesCreateOptionSlice(ss []string) []ignTypes.CreateOption {
	var res []ignTypes.CreateOption
	for _, s := range ss {
		res = append(res, ignTypes.CreateOption(s))
	}
	return res
}

// golang--
func convertStringSliceToTypesMountOptionSlice(ss []string) []ignTypes.MountOption {
	var res []ignTypes.MountOption
	for _, s := range ss {
		res = append(res, ignTypes.MountOption(s))
	}
	return res
}
