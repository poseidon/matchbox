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
	"fmt"

	"github.com/alecthomas/units"
	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate"
	"github.com/coreos/ignition/config/validate/report"
)

const (
	BYTES_PER_SECTOR = 512
)

type Disk struct {
	Device     string      `yaml:"device"`
	WipeTable  bool        `yaml:"wipe_table"`
	Partitions []Partition `yaml:"partitions"`
}

type Partition struct {
	Label    string `yaml:"label"`
	Number   int    `yaml:"number"`
	Size     string `yaml:"size"`
	Start    string `yaml:"start"`
	TypeGUID string `yaml:"type_guid"`
}

func init() {
	register2_0(func(in Config, ast validate.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, validate.AstNode) {
		r := report.Report{}
		for disk_idx, disk := range in.Storage.Disks {
			newDisk := ignTypes.Disk{
				Device:    ignTypes.Path(disk.Device),
				WipeTable: disk.WipeTable,
			}

			for part_idx, partition := range disk.Partitions {
				size, err := convertPartitionDimension(partition.Size)
				if err != nil {
					convertReport := report.ReportFromError(err, report.EntryError)
					if sub_node, err := getNodeChildPath(ast, "storage", "disks", disk_idx, "partitions", part_idx, "size"); err == nil {
						convertReport.AddPosition(sub_node.ValueLineCol(nil))
					}
					r.Merge(convertReport)
					// dont add invalid partitions
					continue
				}
				start, err := convertPartitionDimension(partition.Start)
				if err != nil {
					convertReport := report.ReportFromError(err, report.EntryError)
					if sub_node, err := getNodeChildPath(ast, "storage", "disks", disk_idx, "partitions", part_idx, "start"); err == nil {
						convertReport.AddPosition(sub_node.ValueLineCol(nil))
					}
					r.Merge(convertReport)
					// dont add invalid partitions
					continue
				}

				newPart := ignTypes.Partition{
					Label:    ignTypes.PartitionLabel(partition.Label),
					Number:   partition.Number,
					Size:     size,
					Start:    start,
					TypeGUID: ignTypes.PartitionTypeGUID(partition.TypeGUID),
				}
				newDisk.Partitions = append(newDisk.Partitions, newPart)
			}

			out.Storage.Disks = append(out.Storage.Disks, newDisk)
		}
		return out, r, ast
	})
}

func convertPartitionDimension(in string) (ignTypes.PartitionDimension, error) {
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
	return ignTypes.PartitionDimension(uint64(sectors)), nil
}
