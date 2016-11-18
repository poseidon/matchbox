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
	"errors"
	"reflect"
	"testing"

	fuzeTypes "github.com/coreos/fuze/config/types"
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/config/validate/report"
)

func TestParse(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		cfg fuzeTypes.Config
		r   report.Report
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in: in{data: ``},
			out: out{
				cfg: fuzeTypes.Config{},
				r: report.Report{
					Entries: []report.Entry{{
						Kind:    report.EntryWarning,
						Message: "Configuration is empty",
					}},
				},
			},
		},
		{
			in: in{data: `
networkd:
  units:
    - name: bad.blah
      contents: not valid
`},
			out: out{cfg: fuzeTypes.Config{
				Networkd: fuzeTypes.Networkd{
					Units: []fuzeTypes.NetworkdUnit{
						{Name: "bad.blah", Contents: "not valid"},
					},
				},
			}},
		},

		// Config
		{
			in: in{data: `
ignition:
  config:
    append:
      - source: http://example.com/test1
        verification:
          hash:
            function: sha512
            sum: 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
      - source: http://example.com/test2
    replace:
      source: http://example.com/test3
      verification:
        hash:
          function: sha512
          sum: 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
`},
			out: out{cfg: fuzeTypes.Config{
				Ignition: fuzeTypes.Ignition{
					Config: fuzeTypes.IgnitionConfig{
						Append: []fuzeTypes.ConfigReference{
							{
								Source: "http://example.com/test1",
								Verification: fuzeTypes.Verification{
									Hash: fuzeTypes.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							{
								Source: "http://example.com/test2",
							},
						},
						Replace: &fuzeTypes.ConfigReference{
							Source: "http://example.com/test3",
							Verification: fuzeTypes.Verification{
								Hash: fuzeTypes.Hash{
									Function: "sha512",
									Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
								},
							},
						},
					},
				},
			}},
		},

		// Storage
		{
			in: in{data: `
storage:
  disks:
    - device: /dev/sda
      wipe_table: true
      partitions:
        - label: ROOT
          number: 7
          size: 100MB
          start: 50MB
          type_guid: 11111111-1111-1111-1111-111111111111
        - label: DATA
          number: 12
          size: 1GB
          start: 300MB
          type_guid: 00000000-0000-0000-0000-000000000000
        - label: NOTHING
    - device: /dev/sdb
      wipe_table: true
  raid:
    - name: fast
      level: raid0
      devices:
        - /dev/sdc
        - /dev/sdd
    - name: durable
      level: raid1
      devices:
        - /dev/sde
        - /dev/sdf
        - /dev/sdg
      spares: 1
  filesystems:
    - name: filesystem1
      mount:
        device: /dev/disk/by-partlabel/ROOT
        format: btrfs
        create:
          force: true
          options:
            - -L
            - ROOT
    - name: filesystem2
      mount:
        device: /dev/disk/by-partlabel/DATA
        format: ext4
    - name: filesystem3
      path: /sysroot
  files:
    - path: /opt/file1
      filesystem: filesystem1
      contents:
        inline: file1
      mode: 0644
      user:
        id: 500
      group:
        id: 501
    - path: /opt/file2
      filesystem: filesystem1
      contents:
        remote:
          url: http://example.com/file2
          compression: gzip
          verification:
            hash:
              function: sha512
              sum: 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
      mode: 0644
      user:
        id: 502
      group:
        id: 503
    - path: /opt/file3
      filesystem: filesystem2
      contents:
        remote:
          url: http://example.com/file3
          compression: gzip
      mode: 0400
      user:
        id: 1000
      group:
        id: 1001
    - path: /opt/file4
      filesystem: filesystem2
`},
			out: out{cfg: fuzeTypes.Config{
				Storage: fuzeTypes.Storage{
					Disks: []fuzeTypes.Disk{
						{
							Device:    "/dev/sda",
							WipeTable: true,
							Partitions: []fuzeTypes.Partition{
								{
									Label:    "ROOT",
									Number:   7,
									Size:     "100MB",
									Start:    "50MB",
									TypeGUID: "11111111-1111-1111-1111-111111111111",
								},
								{
									Label:    "DATA",
									Number:   12,
									Size:     "1GB",
									Start:    "300MB",
									TypeGUID: "00000000-0000-0000-0000-000000000000",
								},
								{
									Label: "NOTHING",
								},
							},
						},
						{
							Device:    "/dev/sdb",
							WipeTable: true,
						},
					},
					Arrays: []fuzeTypes.Raid{
						{
							Name:    "fast",
							Level:   "raid0",
							Devices: []string{"/dev/sdc", "/dev/sdd"},
						},
						{
							Name:    "durable",
							Level:   "raid1",
							Devices: []string{"/dev/sde", "/dev/sdf", "/dev/sdg"},
							Spares:  1,
						},
					},
					Filesystems: []fuzeTypes.Filesystem{
						{
							Name: "filesystem1",
							Mount: &fuzeTypes.Mount{
								Device: "/dev/disk/by-partlabel/ROOT",
								Format: "btrfs",
								Create: &fuzeTypes.Create{
									Force:   true,
									Options: []string{"-L", "ROOT"},
								},
							},
						},
						{
							Name: "filesystem2",
							Mount: &fuzeTypes.Mount{
								Device: "/dev/disk/by-partlabel/DATA",
								Format: "ext4",
							},
						},
						{
							Name: "filesystem3",
							Path: "/sysroot",
						},
					},
					Files: []fuzeTypes.File{
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file1",
							Contents: fuzeTypes.FileContents{
								Inline: "file1",
							},
							Mode:  0644,
							User:  fuzeTypes.FileUser{Id: 500},
							Group: fuzeTypes.FileGroup{Id: 501},
						},
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file2",
							Contents: fuzeTypes.FileContents{
								Remote: fuzeTypes.Remote{
									Url:         "http://example.com/file2",
									Compression: "gzip",
									Verification: fuzeTypes.Verification{
										Hash: fuzeTypes.Hash{
											Function: "sha512",
											Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
										},
									},
								},
							},
							Mode:  0644,
							User:  fuzeTypes.FileUser{Id: 502},
							Group: fuzeTypes.FileGroup{Id: 503},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file3",
							Contents: fuzeTypes.FileContents{
								Remote: fuzeTypes.Remote{
									Url:         "http://example.com/file3",
									Compression: "gzip",
								},
							},
							Mode:  0400,
							User:  fuzeTypes.FileUser{Id: 1000},
							Group: fuzeTypes.FileGroup{Id: 1001},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file4",
							Contents: fuzeTypes.FileContents{
								Inline: "",
							},
						},
					},
				},
			}},
		},

		// systemd
		{
			in: in{data: `
systemd:
  units:
    - name: test1.service
      enable: true
      contents: test1 contents
      dropins:
        - name: conf1.conf
          contents: conf1 contents
        - name: conf2.conf
          contents: conf2 contents
    - name: test2.service
      mask: true
      contents: test2 contents
`},
			out: out{cfg: fuzeTypes.Config{
				Systemd: fuzeTypes.Systemd{
					Units: []fuzeTypes.SystemdUnit{
						{
							Name:     "test1.service",
							Enable:   true,
							Contents: "test1 contents",
							DropIns: []fuzeTypes.SystemdUnitDropIn{
								{
									Name:     "conf1.conf",
									Contents: "conf1 contents",
								},
								{
									Name:     "conf2.conf",
									Contents: "conf2 contents",
								},
							},
						},
						{
							Name:     "test2.service",
							Mask:     true,
							Contents: "test2 contents",
						},
					},
				},
			}},
		},

		// networkd
		{
			in: in{data: `
networkd:
  units:
    - name: empty.netdev
    - name: test.network
      contents: test config
`},
			out: out{cfg: fuzeTypes.Config{
				Networkd: fuzeTypes.Networkd{
					Units: []fuzeTypes.NetworkdUnit{
						{
							Name: "empty.netdev",
						},
						{
							Name:     "test.network",
							Contents: "test config",
						},
					},
				},
			}},
		},

		// passwd
		{
			in: in{data: `
passwd:
  users:
    - name: user 1
      password_hash: password 1
      ssh_authorized_keys:
        - key1
        - key2
    - name: user 2
      password_hash: password 2
      ssh_authorized_keys:
        - key3
        - key4
      create:
        uid: 123
        gecos: gecos
        home_dir: /home/user 2
        no_create_home: true
        primary_group: wheel
        groups:
          - wheel
          - plugdev
        no_user_group: true
        system: true
        no_log_init: true
        shell: /bin/zsh
    - name: user 3
      password_hash: password 3
      ssh_authorized_keys:
        - key5
        - key6
      create: {}
  groups:
    - name: group 1
      gid: 1000
      password_hash: password 1
      system: true
    - name: group 2
      password_hash: password 2
`},
			out: out{cfg: fuzeTypes.Config{
				Passwd: fuzeTypes.Passwd{
					Users: []fuzeTypes.User{
						{
							Name:              "user 1",
							PasswordHash:      "password 1",
							SSHAuthorizedKeys: []string{"key1", "key2"},
						},
						{
							Name:              "user 2",
							PasswordHash:      "password 2",
							SSHAuthorizedKeys: []string{"key3", "key4"},
							Create: &fuzeTypes.UserCreate{
								Uid:          func(i uint) *uint { return &i }(123),
								GECOS:        "gecos",
								Homedir:      "/home/user 2",
								NoCreateHome: true,
								PrimaryGroup: "wheel",
								Groups:       []string{"wheel", "plugdev"},
								NoUserGroup:  true,
								System:       true,
								NoLogInit:    true,
								Shell:        "/bin/zsh",
							},
						},
						{
							Name:              "user 3",
							PasswordHash:      "password 3",
							SSHAuthorizedKeys: []string{"key5", "key6"},
							Create:            &fuzeTypes.UserCreate{},
						},
					},
					Groups: []fuzeTypes.Group{
						{
							Name:         "group 1",
							Gid:          func(i uint) *uint { return &i }(1000),
							PasswordHash: "password 1",
							System:       true,
						},
						{
							Name:         "group 2",
							PasswordHash: "password 2",
						},
					},
				},
			}},
		},
	}

	for i, test := range tests {
		cfg, err := Parse([]byte(test.in.data))
		if !reflect.DeepEqual(err, test.out.r) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.r, err)
		}
		if !reflect.DeepEqual(cfg, test.out.cfg) {
			t.Errorf("#%d: bad config: want %#v, got %#v", i, test.out.cfg, cfg)
		}
	}
}
func TestConvertAs2_0_0(t *testing.T) {
	type in struct {
		cfg fuzeTypes.Config
	}
	type out struct {
		cfg types.Config
		r   report.Report
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{cfg: fuzeTypes.Config{}},
			out: out{cfg: types.Config{Ignition: types.Ignition{Version: types.IgnitionVersion{Major: 2}}}},
		},
		{
			in: in{cfg: fuzeTypes.Config{
				Networkd: fuzeTypes.Networkd{
					Units: []fuzeTypes.NetworkdUnit{
						{Name: "bad.blah", Contents: "not valid"},
					},
				},
			}},
			out: out{r: report.ReportFromError(errors.New("invalid networkd unit extension"), report.EntryError)},
		},

		// Config
		{
			in: in{cfg: fuzeTypes.Config{
				Ignition: fuzeTypes.Ignition{
					Config: fuzeTypes.IgnitionConfig{
						Append: []fuzeTypes.ConfigReference{
							{
								Source: "http://example.com/test1",
								Verification: fuzeTypes.Verification{
									Hash: fuzeTypes.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							{
								Source: "http://example.com/test2",
							},
						},
						Replace: &fuzeTypes.ConfigReference{
							Source: "http://example.com/test3",
							Verification: fuzeTypes.Verification{
								Hash: fuzeTypes.Hash{
									Function: "sha512",
									Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
								},
							},
						},
					},
				},
			}},
			out: out{cfg: types.Config{
				Ignition: types.Ignition{
					Version: types.IgnitionVersion{Major: 2},
					Config: types.IgnitionConfig{
						Append: []types.ConfigReference{
							{
								Source: types.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/test1",
								},
								Verification: types.Verification{
									Hash: &types.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							{
								Source: types.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/test2",
								},
							},
						},
						Replace: &types.ConfigReference{
							Source: types.Url{
								Scheme: "http",
								Host:   "example.com",
								Path:   "/test3",
							},
							Verification: types.Verification{
								Hash: &types.Hash{
									Function: "sha512",
									Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
								},
							},
						},
					},
				},
			}},
		},

		// Storage
		{
			in: in{cfg: fuzeTypes.Config{
				Storage: fuzeTypes.Storage{
					Disks: []fuzeTypes.Disk{
						{
							Device:    "/dev/sda",
							WipeTable: true,
							Partitions: []fuzeTypes.Partition{
								{
									Label:    "ROOT",
									Number:   7,
									Size:     "100MB",
									Start:    "50MB",
									TypeGUID: "11111111-1111-1111-1111-111111111111",
								},
								{
									Label:    "DATA",
									Number:   12,
									Size:     "1GB",
									Start:    "300MB",
									TypeGUID: "00000000-0000-0000-0000-000000000000",
								},
								{
									Label: "NOTHING",
								},
							},
						},
						{
							Device:    "/dev/sdb",
							WipeTable: true,
						},
					},
					Arrays: []fuzeTypes.Raid{
						{
							Name:    "fast",
							Level:   "raid0",
							Devices: []string{"/dev/sdc", "/dev/sdd"},
						},
						{
							Name:    "durable",
							Level:   "raid1",
							Devices: []string{"/dev/sde", "/dev/sdf", "/dev/sdg"},
							Spares:  1,
						},
					},
					Filesystems: []fuzeTypes.Filesystem{
						{
							Name: "filesystem1",
							Mount: &fuzeTypes.Mount{
								Device: "/dev/disk/by-partlabel/ROOT",
								Format: "btrfs",
								Create: &fuzeTypes.Create{
									Force:   true,
									Options: []string{"-L", "ROOT"},
								},
							},
						},
						{
							Name: "filesystem2",
							Mount: &fuzeTypes.Mount{
								Device: "/dev/disk/by-partlabel/DATA",
								Format: "ext4",
							},
						},
						{
							Name: "filesystem3",
							Path: "/sysroot",
						},
					},
					Files: []fuzeTypes.File{
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file1",
							Contents: fuzeTypes.FileContents{
								Inline: "file1",
							},
							Mode:  0644,
							User:  fuzeTypes.FileUser{Id: 500},
							Group: fuzeTypes.FileGroup{Id: 501},
						},
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file2",
							Contents: fuzeTypes.FileContents{
								Remote: fuzeTypes.Remote{
									Url:         "http://example.com/file2",
									Compression: "gzip",
									Verification: fuzeTypes.Verification{
										Hash: fuzeTypes.Hash{
											Function: "sha512",
											Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
										},
									},
								},
							},
							Mode:  0644,
							User:  fuzeTypes.FileUser{Id: 502},
							Group: fuzeTypes.FileGroup{Id: 503},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file3",
							Contents: fuzeTypes.FileContents{
								Remote: fuzeTypes.Remote{
									Url:         "http://example.com/file3",
									Compression: "gzip",
								},
							},
							Mode:  0400,
							User:  fuzeTypes.FileUser{Id: 1000},
							Group: fuzeTypes.FileGroup{Id: 1001},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file4",
							Contents: fuzeTypes.FileContents{
								Inline: "",
							},
						},
					},
				},
			}},
			out: out{cfg: types.Config{
				Ignition: types.Ignition{Version: types.IgnitionVersion{Major: 2}},
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device:    types.Path("/dev/sda"),
							WipeTable: true,
							Partitions: []types.Partition{
								{
									Label:    types.PartitionLabel("ROOT"),
									Number:   7,
									Size:     types.PartitionDimension(0x32000),
									Start:    types.PartitionDimension(0x19000),
									TypeGUID: "11111111-1111-1111-1111-111111111111",
								},
								{
									Label:    types.PartitionLabel("DATA"),
									Number:   12,
									Size:     types.PartitionDimension(0x200000),
									Start:    types.PartitionDimension(0x96000),
									TypeGUID: "00000000-0000-0000-0000-000000000000",
								},
								{
									Label: types.PartitionLabel("NOTHING"),
								},
							},
						},
						{
							Device:    types.Path("/dev/sdb"),
							WipeTable: true,
						},
					},
					Arrays: []types.Raid{
						{
							Name:    "fast",
							Level:   "raid0",
							Devices: []types.Path{types.Path("/dev/sdc"), types.Path("/dev/sdd")},
						},
						{
							Name:    "durable",
							Level:   "raid1",
							Devices: []types.Path{types.Path("/dev/sde"), types.Path("/dev/sdf"), types.Path("/dev/sdg")},
							Spares:  1,
						},
					},
					Filesystems: []types.Filesystem{
						{
							Name: "filesystem1",
							Mount: &types.FilesystemMount{
								Device: types.Path("/dev/disk/by-partlabel/ROOT"),
								Format: types.FilesystemFormat("btrfs"),
								Create: &types.FilesystemCreate{
									Force:   true,
									Options: types.MkfsOptions([]string{"-L", "ROOT"}),
								},
							},
						},
						{
							Name: "filesystem2",
							Mount: &types.FilesystemMount{
								Device: types.Path("/dev/disk/by-partlabel/DATA"),
								Format: types.FilesystemFormat("ext4"),
							},
						},
						{
							Name: "filesystem3",
							Path: func(p types.Path) *types.Path { return &p }("/sysroot"),
						},
					},
					Files: []types.File{
						{
							Filesystem: "filesystem1",
							Path:       types.Path("/opt/file1"),
							Contents: types.FileContents{
								Source: types.Url{
									Scheme: "data",
									Opaque: ",file1",
								},
							},
							Mode:  types.FileMode(0644),
							User:  types.FileUser{Id: 500},
							Group: types.FileGroup{Id: 501},
						},
						{
							Filesystem: "filesystem1",
							Path:       types.Path("/opt/file2"),
							Contents: types.FileContents{
								Source: types.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/file2",
								},
								Compression: "gzip",
								Verification: types.Verification{
									Hash: &types.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							Mode:  types.FileMode(0644),
							User:  types.FileUser{Id: 502},
							Group: types.FileGroup{Id: 503},
						},
						{
							Filesystem: "filesystem2",
							Path:       types.Path("/opt/file3"),
							Contents: types.FileContents{
								Source: types.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/file3",
								},
								Compression: "gzip",
							},
							Mode:  types.FileMode(0400),
							User:  types.FileUser{Id: 1000},
							Group: types.FileGroup{Id: 1001},
						},
						{
							Filesystem: "filesystem2",
							Path:       types.Path("/opt/file4"),
							Contents: types.FileContents{
								Source: types.Url{
									Scheme: "data",
									Opaque: ",",
								},
							},
						},
					},
				},
			}},
		},

		// systemd
		{
			in: in{cfg: fuzeTypes.Config{
				Systemd: fuzeTypes.Systemd{
					Units: []fuzeTypes.SystemdUnit{
						{
							Name:     "test1.service",
							Enable:   true,
							Contents: "test1 contents",
							DropIns: []fuzeTypes.SystemdUnitDropIn{
								{
									Name:     "conf1.conf",
									Contents: "conf1 contents",
								},
								{
									Name:     "conf2.conf",
									Contents: "conf2 contents",
								},
							},
						},
						{
							Name:     "test2.service",
							Mask:     true,
							Contents: "test2 contents",
						},
					},
				},
			}},
			out: out{cfg: types.Config{
				Ignition: types.Ignition{Version: types.IgnitionVersion{Major: 2}},
				Systemd: types.Systemd{
					Units: []types.SystemdUnit{
						{
							Name:     "test1.service",
							Enable:   true,
							Contents: "test1 contents",
							DropIns: []types.SystemdUnitDropIn{
								{
									Name:     "conf1.conf",
									Contents: "conf1 contents",
								},
								{
									Name:     "conf2.conf",
									Contents: "conf2 contents",
								},
							},
						},
						{
							Name:     "test2.service",
							Mask:     true,
							Contents: "test2 contents",
						},
					},
				},
			}},
		},

		// networkd
		{
			in: in{cfg: fuzeTypes.Config{
				Networkd: fuzeTypes.Networkd{
					Units: []fuzeTypes.NetworkdUnit{
						{
							Name: "empty.netdev",
						},
						{
							Name:     "test.network",
							Contents: "test config",
						},
					},
				},
			}},
			out: out{cfg: types.Config{
				Ignition: types.Ignition{Version: types.IgnitionVersion{Major: 2}},
				Networkd: types.Networkd{
					Units: []types.NetworkdUnit{
						{
							Name: "empty.netdev",
						},
						{
							Name:     "test.network",
							Contents: "test config",
						},
					},
				},
			}},
		},

		// passwd
		{
			in: in{cfg: fuzeTypes.Config{
				Passwd: fuzeTypes.Passwd{
					Users: []fuzeTypes.User{
						{
							Name:              "user 1",
							PasswordHash:      "password 1",
							SSHAuthorizedKeys: []string{"key1", "key2"},
						},
						{
							Name:              "user 2",
							PasswordHash:      "password 2",
							SSHAuthorizedKeys: []string{"key3", "key4"},
							Create: &fuzeTypes.UserCreate{
								Uid:          func(i uint) *uint { return &i }(123),
								GECOS:        "gecos",
								Homedir:      "/home/user 2",
								NoCreateHome: true,
								PrimaryGroup: "wheel",
								Groups:       []string{"wheel", "plugdev"},
								NoUserGroup:  true,
								System:       true,
								NoLogInit:    true,
								Shell:        "/bin/zsh",
							},
						},
						{
							Name:              "user 3",
							PasswordHash:      "password 3",
							SSHAuthorizedKeys: []string{"key5", "key6"},
							Create:            &fuzeTypes.UserCreate{},
						},
					},
					Groups: []fuzeTypes.Group{
						{
							Name:         "group 1",
							Gid:          func(i uint) *uint { return &i }(1000),
							PasswordHash: "password 1",
							System:       true,
						},
						{
							Name:         "group 2",
							PasswordHash: "password 2",
						},
					},
				},
			}},
			out: out{cfg: types.Config{
				Ignition: types.Ignition{Version: types.IgnitionVersion{Major: 2}},
				Passwd: types.Passwd{
					Users: []types.User{
						{
							Name:              "user 1",
							PasswordHash:      "password 1",
							SSHAuthorizedKeys: []string{"key1", "key2"},
						},
						{
							Name:              "user 2",
							PasswordHash:      "password 2",
							SSHAuthorizedKeys: []string{"key3", "key4"},
							Create: &types.UserCreate{
								Uid:          func(i uint) *uint { return &i }(123),
								GECOS:        "gecos",
								Homedir:      "/home/user 2",
								NoCreateHome: true,
								PrimaryGroup: "wheel",
								Groups:       []string{"wheel", "plugdev"},
								NoUserGroup:  true,
								System:       true,
								NoLogInit:    true,
								Shell:        "/bin/zsh",
							},
						},
						{
							Name:              "user 3",
							PasswordHash:      "password 3",
							SSHAuthorizedKeys: []string{"key5", "key6"},
							Create:            &types.UserCreate{},
						},
					},
					Groups: []types.Group{
						{
							Name:         "group 1",
							Gid:          func(i uint) *uint { return &i }(1000),
							PasswordHash: "password 1",
							System:       true,
						},
						{
							Name:         "group 2",
							PasswordHash: "password 2",
						},
					},
				},
			}},
		},
	}

	for i, test := range tests {
		cfg, r := ConvertAs2_0_0(test.in.cfg)
		if !reflect.DeepEqual(r, test.out.r) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.r, r)
		}
		if !reflect.DeepEqual(cfg, test.out.cfg) {
			t.Errorf("#%d: bad config: want %#v, got %#v", i, test.out.cfg, cfg)
		}
	}
}
