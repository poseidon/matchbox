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
	"fmt"
	"net/url"
	"reflect"

	"github.com/alecthomas/units"
	fuzeTypes "github.com/coreos/fuze/config/types"
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/config/validate"
	"github.com/coreos/ignition/config/validate/report"
	"github.com/vincent-petithory/dataurl"
)

const (
	BYTES_PER_SECTOR = 512
)

func ConvertAs2_0_0(in fuzeTypes.Config) (types.Config, report.Report) {
	out := types.Config{
		Ignition: types.Ignition{
			Version: types.IgnitionVersion{Major: 2, Minor: 0},
		},
	}

	for _, ref := range in.Ignition.Config.Append {
		newRef, err := convertConfigReference(ref)
		if err != nil {
			return types.Config{}, report.ReportFromError(err, report.EntryError)
		}
		out.Ignition.Config.Append = append(out.Ignition.Config.Append, newRef)
	}

	if in.Ignition.Config.Replace != nil {
		newRef, err := convertConfigReference(*in.Ignition.Config.Replace)
		if err != nil {
			return types.Config{}, report.ReportFromError(err, report.EntryError)
		}
		out.Ignition.Config.Replace = &newRef
	}

	for _, disk := range in.Storage.Disks {
		newDisk := types.Disk{
			Device:    types.Path(disk.Device),
			WipeTable: disk.WipeTable,
		}

		for _, partition := range disk.Partitions {
			size, err := convertPartitionDimension(partition.Size)
			if err != nil {
				return types.Config{}, report.ReportFromError(err, report.EntryError)
			}
			start, err := convertPartitionDimension(partition.Start)
			if err != nil {
				return types.Config{}, report.ReportFromError(err, report.EntryError)
			}

			newDisk.Partitions = append(newDisk.Partitions, types.Partition{
				Label:    types.PartitionLabel(partition.Label),
				Number:   partition.Number,
				Size:     size,
				Start:    start,
				TypeGUID: types.PartitionTypeGUID(partition.TypeGUID),
			})
		}

		out.Storage.Disks = append(out.Storage.Disks, newDisk)
	}

	for _, array := range in.Storage.Arrays {
		newArray := types.Raid{
			Name:   array.Name,
			Level:  array.Level,
			Spares: array.Spares,
		}

		for _, device := range array.Devices {
			newArray.Devices = append(newArray.Devices, types.Path(device))
		}

		out.Storage.Arrays = append(out.Storage.Arrays, newArray)
	}

	for _, filesystem := range in.Storage.Filesystems {
		newFilesystem := types.Filesystem{
			Name: filesystem.Name,
			Path: func(p types.Path) *types.Path {
				if p == "" {
					return nil
				}

				return &p
			}(types.Path(filesystem.Path)),
		}

		if filesystem.Mount != nil {
			newFilesystem.Mount = &types.FilesystemMount{
				Device: types.Path(filesystem.Mount.Device),
				Format: types.FilesystemFormat(filesystem.Mount.Format),
			}

			if filesystem.Mount.Create != nil {
				newFilesystem.Mount.Create = &types.FilesystemCreate{
					Force:   filesystem.Mount.Create.Force,
					Options: types.MkfsOptions(filesystem.Mount.Create.Options),
				}
			}
		}

		out.Storage.Filesystems = append(out.Storage.Filesystems, newFilesystem)
	}

	for _, file := range in.Storage.Files {
		newFile := types.File{
			Filesystem: file.Filesystem,
			Path:       types.Path(file.Path),
			Mode:       types.FileMode(file.Mode),
			User:       types.FileUser{Id: file.User.Id},
			Group:      types.FileGroup{Id: file.Group.Id},
		}

		if file.Contents.Inline != "" {
			newFile.Contents = types.FileContents{
				Source: types.Url{
					Scheme: "data",
					Opaque: "," + dataurl.EscapeString(file.Contents.Inline),
				},
			}
		}

		if file.Contents.Remote.Url != "" {
			source, err := url.Parse(file.Contents.Remote.Url)
			if err != nil {
				return types.Config{}, report.ReportFromError(err, report.EntryError)
			}

			newFile.Contents = types.FileContents{Source: types.Url(*source)}
		}

		if newFile.Contents == (types.FileContents{}) {
			newFile.Contents = types.FileContents{
				Source: types.Url{
					Scheme: "data",
					Opaque: ",",
				},
			}
		}

		newFile.Contents.Compression = types.Compression(file.Contents.Remote.Compression)
		newFile.Contents.Verification = convertVerification(file.Contents.Remote.Verification)

		out.Storage.Files = append(out.Storage.Files, newFile)
	}

	for _, unit := range in.Systemd.Units {
		newUnit := types.SystemdUnit{
			Name:     types.SystemdUnitName(unit.Name),
			Enable:   unit.Enable,
			Mask:     unit.Mask,
			Contents: unit.Contents,
		}

		for _, dropIn := range unit.DropIns {
			newUnit.DropIns = append(newUnit.DropIns, types.SystemdUnitDropIn{
				Name:     types.SystemdUnitDropInName(dropIn.Name),
				Contents: dropIn.Contents,
			})
		}

		out.Systemd.Units = append(out.Systemd.Units, newUnit)
	}

	for _, unit := range in.Networkd.Units {
		out.Networkd.Units = append(out.Networkd.Units, types.NetworkdUnit{
			Name:     types.NetworkdUnitName(unit.Name),
			Contents: unit.Contents,
		})
	}

	for _, user := range in.Passwd.Users {
		newUser := types.User{
			Name:              user.Name,
			PasswordHash:      user.PasswordHash,
			SSHAuthorizedKeys: user.SSHAuthorizedKeys,
		}

		if user.Create != nil {
			newUser.Create = &types.UserCreate{
				Uid:          user.Create.Uid,
				GECOS:        user.Create.GECOS,
				Homedir:      user.Create.Homedir,
				NoCreateHome: user.Create.NoCreateHome,
				PrimaryGroup: user.Create.PrimaryGroup,
				Groups:       user.Create.Groups,
				NoUserGroup:  user.Create.NoUserGroup,
				System:       user.Create.System,
				NoLogInit:    user.Create.NoLogInit,
				Shell:        user.Create.Shell,
			}
		}

		out.Passwd.Users = append(out.Passwd.Users, newUser)
	}

	for _, group := range in.Passwd.Groups {
		out.Passwd.Groups = append(out.Passwd.Groups, types.Group{
			Name:         group.Name,
			Gid:          group.Gid,
			PasswordHash: group.PasswordHash,
			System:       group.System,
		})
	}

	r := validate.ValidateWithoutSource(reflect.ValueOf(out))
	if r.IsFatal() {
		return types.Config{}, r
	}

	return out, r
}

func convertConfigReference(in fuzeTypes.ConfigReference) (types.ConfigReference, error) {
	source, err := url.Parse(in.Source)
	if err != nil {
		return types.ConfigReference{}, err
	}

	return types.ConfigReference{
		Source:       types.Url(*source),
		Verification: convertVerification(in.Verification),
	}, nil
}

func convertVerification(in fuzeTypes.Verification) types.Verification {
	if in.Hash.Function == "" || in.Hash.Sum == "" {
		return types.Verification{}
	}

	return types.Verification{
		&types.Hash{
			Function: in.Hash.Function,
			Sum:      in.Hash.Sum,
		},
	}
}

func convertPartitionDimension(in string) (types.PartitionDimension, error) {
	if in == "" {
		return 0, nil
	}

	b, err := units.ParseBase2Bytes(in)
	if err != nil {
		return 0, err
	}
	if b < 0 {
		return 0, fmt.Errorf("invalid dimension (negative): %q", in)
	}

	// Translate bytes into sectors
	sectors := (b / BYTES_PER_SECTOR)
	if b%BYTES_PER_SECTOR != 0 {
		sectors++
	}
	return types.PartitionDimension(uint64(sectors)), nil
}
