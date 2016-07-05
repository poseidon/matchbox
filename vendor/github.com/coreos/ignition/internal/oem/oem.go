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

	"github.com/vincent-petithory/dataurl"
)

// Config represents a set of command line flags that map to a particular OEM.
type Config struct {
	name              string
	flags             map[string]string
	provider          providers.ProviderCreator
	baseConfig        types.Config
	defaultUserConfig types.Config
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

func (c Config) BaseConfig() types.Config {
	return c.baseConfig
}

func (c Config) DefaultUserConfig() types.Config {
	return c.defaultUserConfig
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
		baseConfig: types.Config{
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
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.SystemdUnit{
					{Enable: true, Name: "coreos-metadata-sshkeys@.service"},
					{Enable: true, Name: "google-accounts-manager.service"},
					{Enable: true, Name: "google-address-manager.service"},
					{Enable: true, Name: "google-clock-sync-manager.service"},
					{Enable: true, Name: "google-startup-scripts-onboot.service"},
					{Enable: true, Name: "google-startup-scripts.service"},
				},
			},
			Storage: types.Storage{
				Files: []types.File{
					serviceFromOem("google-accounts-manager.service"),
					serviceFromOem("google-address-manager.service"),
					serviceFromOem("google-clock-sync-manager.service"),
					serviceFromOem("google-startup-scripts-onboot.service"),
					serviceFromOem("google-startup-scripts.service"),
					{
						Filesystem: "root",
						Path:       "/etc/hosts",
						Mode:       0444,
						Contents:   contentsFromString("169.254.169.254 metadata\n127.0.0.1 localhost\n"),
					},
					{
						Filesystem: "root",
						Path:       "/etc/profile.d/google-cloud-sdk.sh",
						Mode:       0444,
						Contents: contentsFromString(`#!/bin/sh
alias gcloud="(docker images google/cloud-sdk || docker pull google/cloud-sdk) > /dev/null;docker run -t -i --net="host" -v $HOME/.config:/.config -v /var/run/docker.sock:/var/run/doker.sock google/cloud-sdk gcloud"
alias gcutil="(docker images google/cloud-sdk || docker pull google/cloud-sdk) > /dev/null;docker run -t -i --net="host" -v $HOME/.config:/.config google/cloud-sdk gcutil"
alias gsutil="(docker images google/cloud-sdk || docker pull google/cloud-sdk) > /dev/null;docker run -t -i --net="host" -v $HOME/.config:/.config google/cloud-sdk gsutil"
`),
					},
				},
			},
		},
		defaultUserConfig: types.Config{Systemd: types.Systemd{Units: []types.SystemdUnit{userCloudInit("GCE", "gce")}}},
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

func contentsFromString(data string) types.FileContents {
	return types.FileContents{
		Source: types.Url{
			Scheme: "data",
			Opaque: "," + dataurl.EscapeString(data),
		},
	}
}

func contentsFromOem(path string) types.FileContents {
	return types.FileContents{
		Source: types.Url{
			Scheme: "oem",
			Path:   path,
		},
	}
}

func userCloudInit(name string, oem string) types.SystemdUnit {
	contents := `[Unit]
Description=Cloudinit from %s metadata

[Service]
Type=oneshot
ExecStart=/usr/bin/coreos-cloudinit --oem=%s

[Install]
WantedBy=multi-user.target
`

	return types.SystemdUnit{
		Name:     "oem-cloudinit.service",
		Enable:   true,
		Contents: fmt.Sprintf(contents, name, oem),
	}
}

func serviceFromOem(unit string) types.File {
	return types.File{
		Filesystem: "root",
		Path:       types.Path("/etc/systemd/system/" + unit),
		Mode:       0444,
		Contents:   contentsFromOem("/units/" + unit),
	}
}
