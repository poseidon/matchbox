// Copyright 2015 CoreOS, Inc.
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

// The vmware provider fetches a configuration from the VMware Guest Info
// interface.

package vmware

import (
	"github.com/coreos/ignition/config"
	"github.com/coreos/ignition/config/types"

	"github.com/sigma/vmw-guestinfo/rpcvmx"
	"github.com/sigma/vmw-guestinfo/vmcheck"
)

func (p provider) FetchConfig() (types.Config, error) {
	info := rpcvmx.NewConfig()
	data, err := info.String("coreos.config.data", "")
	if err != nil {
		p.logger.Debug("failed to fetch config: %v", err)
		return types.Config{}, err
	}

	encoding, err := info.String("coreos.config.data.encoding", "")
	if err != nil {
		p.logger.Debug("failed to fetch config encoding: %v", err)
		return types.Config{}, err
	}

	decodedData, err := decodeData(data, encoding)
	if err != nil {
		p.logger.Debug("failed to decode config: %v", err)
		return types.Config{}, err
	}

	p.logger.Debug("config successfully fetched")
	return config.Parse(decodedData)
}

func (p *provider) IsOnline() bool {
	return vmcheck.IsVirtualWorld()
}
