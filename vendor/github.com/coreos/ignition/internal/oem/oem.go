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
	"github.com/coreos/ignition/internal/providers/digitalocean"
	"github.com/coreos/ignition/internal/providers/ec2"
	"github.com/coreos/ignition/internal/providers/file"
	"github.com/coreos/ignition/internal/providers/gce"
	"github.com/coreos/ignition/internal/providers/noop"
	"github.com/coreos/ignition/internal/providers/openstack"
	"github.com/coreos/ignition/internal/providers/packet"
	"github.com/coreos/ignition/internal/providers/qemu"
	"github.com/coreos/ignition/internal/providers/vmware"
	"github.com/coreos/ignition/internal/registry"

	"github.com/vincent-petithory/dataurl"
)

// Config represents a set of options that map to a particular OEM.
type Config struct {
	name              string
	fetch             providers.FuncFetchConfig
	baseConfig        types.Config
	defaultUserConfig types.Config
}

func (c Config) Name() string {
	return c.name
}

func (c Config) FetchFunc() providers.FuncFetchConfig {
	return c.fetch
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
		name:  "azure",
		fetch: azure.FetchConfig,
	})
	configs.Register(Config{
		name:  "cloudsigma",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "cloudstack",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "digitalocean",
		fetch: digitalocean.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.SystemdUnit{{Enable: true, Name: "coreos-metadata-sshkeys@.service"}},
			},
		},
		defaultUserConfig: types.Config{Systemd: types.Systemd{Units: []types.SystemdUnit{userCloudInit("DigitalOcean", "digitalocean")}}},
	})
	configs.Register(Config{
		name:  "brightbox",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "openstack",
		fetch: openstack.FetchConfig,
	})
	configs.Register(Config{
		name:  "ec2",
		fetch: ec2.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.SystemdUnit{{Enable: true, Name: "coreos-metadata-sshkeys@.service"}},
			},
		},
	})
	configs.Register(Config{
		name:  "exoscale",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "gce",
		fetch: gce.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.SystemdUnit{
					{Enable: true, Name: "coreos-metadata-sshkeys@.service"},
					{Enable: true, Name: "oem-gce.service"},
				},
			},
			Storage: types.Storage{
				Files: []types.File{
					serviceFromOem("oem-gce.service"),
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
		name:  "hyperv",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "niftycloud",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "packet",
		fetch: packet.FetchConfig,
	})
	configs.Register(Config{
		name:  "pxe",
		fetch: cmdline.FetchConfig,
	})
	configs.Register(Config{
		name:  "rackspace",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "rackspace-onmetal",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "vagrant",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "vmware",
		fetch: vmware.FetchConfig,
	})
	configs.Register(Config{
		name:  "xendom0",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "interoute",
		fetch: noop.FetchConfig,
	})
	configs.Register(Config{
		name:  "qemu",
		fetch: qemu.FetchConfig,
	})
	configs.Register(Config{
		name:  "file",
		fetch: file.FetchConfig,
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
