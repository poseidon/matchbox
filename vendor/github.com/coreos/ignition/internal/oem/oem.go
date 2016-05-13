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

package oem

import (
	"fmt"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/providers"
	"github.com/coreos/ignition/internal/providers/azure"
	"github.com/coreos/ignition/internal/providers/cmdline"
	"github.com/coreos/ignition/internal/providers/ec2"
	"github.com/coreos/ignition/internal/providers/gce"
	"github.com/coreos/ignition/internal/providers/noop"
	"github.com/coreos/ignition/internal/providers/vmware"
	"github.com/coreos/ignition/internal/registry"
)

// Config represents a set of command line flags that map to a particular OEM.
type Config struct {
	name     string
	flags    map[string]string
	provider providers.ProviderCreator
	config   types.Config
}

func (c Config) Name() string {
	return c.name
}

func (c Config) Flags() map[string]string {
	return c.flags
}

func (c Config) Provider() providers.ProviderCreator {
	return c.provider
}

func (c Config) Config() types.Config {
	return c.config
}

var configs = registry.Create("oem configs")

func init() {
	configs.Register(Config{
		name:     "azure",
		provider: azure.Creator{},
	})
	configs.Register(Config{
		name:     "cloudsigma",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "cloudstack",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "digitalocean",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "brightbox",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "openstack",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "ec2",
		provider: ec2.Creator{},
		flags: map[string]string{
			"online-timeout": "0",
		},
		config: types.Config{
			Systemd: types.Systemd{
				Units: []types.SystemdUnit{{
					Name:   "coreos-metadata-sshkeys@.service",
					Enable: true,
				}},
			},
		},
	})
	configs.Register(Config{
		name:     "exoscale",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "gce",
		provider: gce.Creator{},
		config: types.Config{
			Systemd: types.Systemd{
				Units: []types.SystemdUnit{{
					Name:   "coreos-metadata-sshkeys@.service",
					Enable: true,
				}},
			},
		},
	})
	configs.Register(Config{
		name:     "hyperv",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "niftycloud",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "packet",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "pxe",
		provider: cmdline.Creator{},
	})
	configs.Register(Config{
		name:     "rackspace",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "rackspace-onmetal",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "vagrant",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "vmware",
		provider: vmware.Creator{},
	})
	configs.Register(Config{
		name:     "xendom0",
		provider: noop.Creator{},
	})
	configs.Register(Config{
		name:     "interoute",
		provider: noop.Creator{},
	})
}

func Get(name string) (config Config, ok bool) {
	config, ok = configs.Get(name).(Config)
	return
}

func MustGet(name string) Config {
	if config, ok := Get(name); ok {
		return config
	} else {
		panic(fmt.Sprintf("invalid OEM name %q provided", name))
	}
}

func Names() (names []string) {
	return configs.Names()
}
