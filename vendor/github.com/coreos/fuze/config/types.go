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

type Config struct {
	Ignition struct {
		Config struct {
			Append  []ConfigReference `yaml:"append"`
			Replace *ConfigReference  `yaml:"replace"`
		} `yaml:"config"`
	} `yaml:"ignition"`
	Storage struct {
		Disks []struct {
			Device     string `yaml:"device"`
			WipeTable  bool   `yaml:"wipe_table"`
			Partitions []struct {
				Label    string `yaml:"label"`
				Number   int    `yaml:"number"`
				Size     string `yaml:"size"`
				Start    string `yaml:"start"`
				TypeGUID string `yaml:"type_guid"`
			} `yaml:"partitions"`
		} `yaml:"disks"`
		Arrays []struct {
			Name    string   `yaml:"name"`
			Level   string   `yaml:"level"`
			Devices []string `yaml:"devices"`
			Spares  int      `yaml:"spares"`
		} `yaml:"raid"`
		Filesystems []struct {
			Name  string `yaml:"name"`
			Mount *struct {
				Device string `yaml:"device"`
				Format string `yaml:"format"`
				Create *struct {
					Force   bool     `yaml:"force"`
					Options []string `yaml:"options"`
				} `yaml:"create"`
			} `yaml:"mount"`
			Path string `yaml:"path"`
		} `yaml:"filesystems"`
		Files []struct {
			Filesystem string `yaml:"filesystem"`
			Path       string `yaml:"path"`
			Contents   struct {
				Remote struct {
					Url          string       `yaml:"url"`
					Compression  string       `yaml:"compression"`
					Verification Verification `yaml:"verification"`
				} `yaml:"remote"`
				Inline string `yaml:"inline"`
			} `yaml:"contents"`
			Mode int `yaml:"mode"`
			User struct {
				Id int `yaml:"id"`
			} `yaml:"user"`
			Group struct {
				Id int `yaml:"id"`
			} `yaml:"group"`
		} `yaml:"files"`
	} `yaml:"storage"`
	Systemd struct {
		Units []struct {
			Name     string `yaml:"name"`
			Enable   bool   `yaml:"enable"`
			Mask     bool   `yaml:"mask"`
			Contents string `yaml:"contents"`
			DropIns  []struct {
				Name     string `yaml:"name"`
				Contents string `yaml:"contents"`
			} `yaml:"dropins"`
		} `yaml:"units"`
	} `yaml:"systemd"`
	Networkd struct {
		Units []struct {
			Name     string `yaml:"name"`
			Contents string `yaml:"contents"`
		} `yaml:"units"`
	} `yaml:"networkd"`
	Passwd struct {
		Users []struct {
			Name              string   `yaml:"name"`
			PasswordHash      string   `yaml:"password_hash"`
			SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
			Create            *struct {
				Uid          *uint    `yaml:"uid"`
				GECOS        string   `yaml:"gecos"`
				Homedir      string   `yaml:"home_dir"`
				NoCreateHome bool     `yaml:"no_create_home"`
				PrimaryGroup string   `yaml:"primary_group"`
				Groups       []string `yaml:"groups"`
				NoUserGroup  bool     `yaml:"no_user_group"`
				System       bool     `yaml:"system"`
				NoLogInit    bool     `yaml:"no_log_init"`
				Shell        string   `yaml:"shell"`
			} `yaml:"create"`
		} `yaml:"users"`
		Groups []struct {
			Name         string `yaml:"name"`
			Gid          *uint  `yaml:"gid"`
			PasswordHash string `yaml:"password_hash"`
			System       bool   `yaml:"system"`
		} `yaml:"groups"`
	} `yaml:"passwd"`
}

type ConfigReference struct {
	Source       string       `yaml:"source"`
	Verification Verification `yaml:"verification"`
}

type Verification struct {
	Hash struct {
		Function string `yaml:"function"`
		Sum      string `yaml:"sum"`
	} `yaml:"hash"`
}
