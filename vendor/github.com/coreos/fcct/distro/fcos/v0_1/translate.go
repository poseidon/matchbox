// Copyright 2019 Red Hat, Inc
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
// limitations under the License.)

package fcos_0_1

import (
	"github.com/coreos/ignition/v2/config/v3_0/types"
)

func (f Fcos) ToIgn3_0() (types.Config, error) {
	return types.Config{
		Ignition: types.Ignition{
			Version: "3.0.0",
		},
	}, nil
}
