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
	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/coreos/ignition/config/validate/astnode"
	"github.com/coreos/ignition/config/validate/report"
)

type Passwd struct {
	Users  []User  `yaml:"users"`
	Groups []Group `yaml:"groups"`
}

type User struct {
	Name              string      `yaml:"name"`
	PasswordHash      *string     `yaml:"password_hash"`
	SSHAuthorizedKeys []string    `yaml:"ssh_authorized_keys"`
	Create            *UserCreate `yaml:"create"`
	UID               *int        `yaml:"uid"`
	Gecos             string      `yaml:"gecos"`
	HomeDir           string      `yaml:"home_dir"`
	NoCreateHome      bool        `yaml:"no_create_home"`
	PrimaryGroup      string      `yaml:"primary_group"`
	Groups            []string    `yaml:"groups"`
	NoUserGroup       bool        `yaml:"no_user_group"`
	System            bool        `yaml:"system"`
	NoLogInit         bool        `yaml:"no_log_init"`
	Shell             string      `yaml:"shell"`
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

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		for _, user := range in.Passwd.Users {
			newUser := ignTypes.PasswdUser{
				Name:              user.Name,
				PasswordHash:      user.PasswordHash,
				SSHAuthorizedKeys: convertStringSliceIntoTypesSSHAuthorizedKeySlice(user.SSHAuthorizedKeys),
				UID:               user.UID,
				Gecos:             user.Gecos,
				HomeDir:           user.HomeDir,
				NoCreateHome:      user.NoCreateHome,
				PrimaryGroup:      user.PrimaryGroup,
				Groups:            convertStringSliceIntoTypesGroupSlice(user.Groups),
				NoUserGroup:       user.NoUserGroup,
				System:            user.System,
				NoLogInit:         user.NoLogInit,
				Shell:             user.Shell,
			}

			if user.Create != nil {
				newUser.Create = &ignTypes.Usercreate{
					UID:          convertUintPointerToIntPointer(user.Create.Uid),
					Gecos:        user.Create.GECOS,
					HomeDir:      user.Create.Homedir,
					NoCreateHome: user.Create.NoCreateHome,
					PrimaryGroup: user.Create.PrimaryGroup,
					Groups:       convertStringSliceIntoTypesUsercreateGroupSlice(user.Create.Groups),
					NoUserGroup:  user.Create.NoUserGroup,
					System:       user.Create.System,
					NoLogInit:    user.Create.NoLogInit,
					Shell:        user.Create.Shell,
				}
			}

			out.Passwd.Users = append(out.Passwd.Users, newUser)
		}

		for _, group := range in.Passwd.Groups {
			out.Passwd.Groups = append(out.Passwd.Groups, ignTypes.PasswdGroup{
				Name:         group.Name,
				Gid:          convertUintPointerToIntPointer(group.Gid),
				PasswordHash: group.PasswordHash,
				System:       group.System,
			})
		}
		return out, report.Report{}, ast
	})
}

// golang--
func convertStringSliceIntoTypesSSHAuthorizedKeySlice(ss []string) []ignTypes.SSHAuthorizedKey {
	var res []ignTypes.SSHAuthorizedKey
	for _, s := range ss {
		res = append(res, ignTypes.SSHAuthorizedKey(s))
	}
	return res
}

// golang--
func convertStringSliceIntoTypesUsercreateGroupSlice(ss []string) []ignTypes.UsercreateGroup {
	var res []ignTypes.UsercreateGroup
	for _, s := range ss {
		res = append(res, ignTypes.UsercreateGroup(s))
	}
	return res
}

// golang--
func convertStringSliceIntoTypesGroupSlice(ss []string) []ignTypes.Group {
	var res []ignTypes.Group
	for _, s := range ss {
		res = append(res, ignTypes.Group(s))
	}
	return res
}

// golang--
func convertUintPointerToIntPointer(u *uint) *int {
	if u == nil {
		return nil
	}
	x := int(*u)
	return &x
}
