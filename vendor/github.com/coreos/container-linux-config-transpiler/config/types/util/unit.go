// Copyright 2017 CoreOS, Inc.
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

package util

type SystemdUnit struct {
	Unit    *UnitSection
	Service *UnitSection
	Install *UnitSection
}

func NewSystemdUnit() SystemdUnit {
	return SystemdUnit{
		Unit:    &UnitSection{},
		Service: &UnitSection{},
		Install: &UnitSection{},
	}
}

type UnitSection []string

func (u *UnitSection) Add(line string) {
	*u = append(*u, line)
}

func (s SystemdUnit) String() string {
	res := ""

	type section struct {
		name     string
		contents []string
	}

	for _, sec := range []section{
		{"Unit", *s.Unit},
		{"Service", *s.Service},
		{"Install", *s.Install},
	} {
		if len(sec.contents) == 0 {
			continue
		}
		if res != "" {
			res += "\n\n"
		}
		res += "[" + sec.name + "]"
		for _, line := range sec.contents {
			res += "\n" + line
		}
	}
	return res
}
