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
	"os"

	"github.com/coreos/ignition/config/validate/report"
)

var (
	ErrNoFilesystem    = errors.New("no filesystem specified")
	ErrFileIllegalMode = errors.New("illegal file mode")
)

// Node represents all common info for files (special types, e.g. directories, included).
type Node struct {
	Filesystem string    `json:"filesystem,omitempty"`
	Path       Path      `json:"path,omitempty"`
	Mode       NodeMode  `json:"mode,omitempty"`
	User       NodeUser  `json:"user,omitempty"`
	Group      NodeGroup `json:"group,omitempty"`
}

type NodeUser struct {
	Id int `json:"id,omitempty"`
}

type NodeGroup struct {
	Id int `json:"id,omitempty"`
}

func (n Node) Validate() report.Report {
	if n.Filesystem == "" {
		return report.ReportFromError(ErrNoFilesystem, report.EntryError)
	}
	return report.Report{}
}

type NodeMode os.FileMode

func (m NodeMode) Validate() report.Report {
	if (m &^ 07777) != 0 {
		return report.ReportFromError(ErrFileIllegalMode, report.EntryError)
	}
	return report.Report{}
}
