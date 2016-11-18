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

type Passwd struct {
	Users  []User  `yaml:"users"`
	Groups []Group `yaml:"groups"`
}

type User struct {
	Name              string      `yaml:"name"`
	PasswordHash      string      `yaml:"password_hash"`
	SSHAuthorizedKeys []string    `yaml:"ssh_authorized_keys"`
	Create            *UserCreate `yaml:"create"`
}

type UserCreate struct {
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
}

type Group struct {
	Name         string `yaml:"name"`
	Gid          *uint  `yaml:"gid"`
	PasswordHash string `yaml:"password_hash"`
	System       bool   `yaml:"system"`
}
