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

package files

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/exec/stages"
	"github.com/coreos/ignition/internal/exec/util"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/resource"
)

const (
	name = "files"
)

var (
	ErrFilesystemUndefined = errors.New("the referenced filesystem was not defined")
)

func init() {
	stages.Register(creator{})
}

type creator struct{}

func (creator) Create(logger *log.Logger, client *resource.HttpClient, root string) stages.Stage {
	return &stage{
		Util: util.Util{
			DestDir: root,
			Logger:  logger,
		},
		client: client,
	}
}

func (creator) Name() string {
	return name
}

type stage struct {
	util.Util

	client *resource.HttpClient
}

func (stage) Name() string {
	return name
}

func (s stage) Run(config types.Config) bool {
	if err := s.createPasswd(config); err != nil {
		s.Logger.Crit("failed to create users/groups: %v", err)
		return false
	}

	if err := s.createFilesystemsFiles(config); err != nil {
		s.Logger.Crit("failed to create files: %v", err)
		return false
	}

	if err := s.createUnits(config); err != nil {
		s.Logger.Crit("failed to create units: %v", err)
		return false
	}

	return true
}

// createFilesystemsFiles creates the files described in config.Storage.Filesystems.
func (s stage) createFilesystemsFiles(config types.Config) error {
	if len(config.Storage.Filesystems) == 0 {
		return nil
	}
	s.Logger.PushPrefix("createFilesystemsFiles")
	defer s.Logger.PopPrefix()

	fileMap, err := s.mapFilesToFilesystems(config)
	if err != nil {
		return err
	}

	for fs, f := range fileMap {
		if err := s.createFiles(fs, f); err != nil {
			return fmt.Errorf("failed to create files: %v", err)
		}
	}

	return nil
}

// mapFilesToFilesystems builds a map of filesystems to files. If multiple
// definitions of the same filesystem are present, only the final definition is
// used.
func (s stage) mapFilesToFilesystems(config types.Config) (map[types.Filesystem][]types.File, error) {
	files := map[string][]types.File{}
	for _, f := range config.Storage.Files {
		files[f.Filesystem] = append(files[f.Filesystem], f)
	}

	filesystems := map[string]types.Filesystem{}
	for _, fs := range config.Storage.Filesystems {
		filesystems[fs.Name] = fs
	}

	fileMap := map[types.Filesystem][]types.File{}
	for fsn, f := range files {
		if fs, ok := filesystems[fsn]; ok {
			fileMap[fs] = append(fileMap[fs], f...)
		} else {
			s.Logger.Crit("the filesystem (%q), was not defined", fsn)
			return nil, ErrFilesystemUndefined
		}
	}

	return fileMap, nil
}

// createFiles creates any files listed for the filesystem in fs.Files.
func (s stage) createFiles(fs types.Filesystem, files []types.File) error {
	s.Logger.PushPrefix("createFiles")
	defer s.Logger.PopPrefix()

	var mnt string
	if fs.Path == nil {
		var err error
		mnt, err = ioutil.TempDir("", "ignition-files")
		if err != nil {
			return fmt.Errorf("failed to create temp directory: %v", err)
		}
		defer os.Remove(mnt)

		dev := string(fs.Mount.Device)
		format := string(fs.Mount.Format)

		if err := s.Logger.LogOp(
			func() error { return syscall.Mount(dev, mnt, format, 0, "") },
			"mounting %q at %q", dev, mnt,
		); err != nil {
			return fmt.Errorf("failed to mount device %q at %q: %v", dev, mnt, err)
		}
		defer s.Logger.LogOp(
			func() error { return syscall.Unmount(mnt, 0) },
			"unmounting %q at %q", dev, mnt,
		)
	} else {
		mnt = string(*fs.Path)
	}

	u := util.Util{
		Logger:  s.Logger,
		DestDir: mnt,
	}
	for _, f := range files {
		file := util.RenderFile(s.Logger, s.client, f)
		if file == nil {
			return fmt.Errorf("failed to resolve file %q", f.Path)
		}

		if err := s.Logger.LogOp(
			func() error { return u.WriteFile(file) },
			"writing file %q", string(f.Path),
		); err != nil {
			return fmt.Errorf("failed to create file %q: %v", file.Path, err)
		}
	}

	return nil
}

// createUnits creates the units listed under systemd.units and networkd.units.
func (s stage) createUnits(config types.Config) error {
	for _, unit := range config.Systemd.Units {
		if err := s.writeSystemdUnit(unit); err != nil {
			return err
		}
		if unit.Enable {
			if err := s.Logger.LogOp(
				func() error { return s.EnableUnit(unit) },
				"enabling unit %q", unit.Name,
			); err != nil {
				return err
			}
		}
		if unit.Mask {
			if err := s.Logger.LogOp(
				func() error { return s.MaskUnit(unit) },
				"masking unit %q", unit.Name,
			); err != nil {
				return err
			}
		}
	}
	for _, unit := range config.Networkd.Units {
		if err := s.writeNetworkdUnit(unit); err != nil {
			return err
		}
	}
	return nil
}

// writeSystemdUnit creates the specified unit and any dropins for that unit.
// If the contents of the unit or are empty, the unit is not created. The same
// applies to the unit's dropins.
func (s stage) writeSystemdUnit(unit types.SystemdUnit) error {
	return s.Logger.LogOp(func() error {
		for _, dropin := range unit.DropIns {
			if dropin.Contents == "" {
				continue
			}

			f := util.FileFromUnitDropin(unit, dropin)
			if err := s.Logger.LogOp(
				func() error { return s.WriteFile(f) },
				"writing drop-in %q at %q", dropin.Name, f.Path,
			); err != nil {
				return err
			}
		}

		if unit.Contents == "" {
			return nil
		}

		f := util.FileFromSystemdUnit(unit)
		if err := s.Logger.LogOp(
			func() error { return s.WriteFile(f) },
			"writing unit %q at %q", unit.Name, f.Path,
		); err != nil {
			return err
		}

		return nil
	}, "processing unit %q", unit.Name)
}

// writeNetworkdUnit creates the specified unit. If the contents of the unit or
// are empty, the unit is not created.
func (s stage) writeNetworkdUnit(unit types.NetworkdUnit) error {
	return s.Logger.LogOp(func() error {
		if unit.Contents == "" {
			return nil
		}

		f := util.FileFromNetworkdUnit(unit)
		if err := s.Logger.LogOp(
			func() error { return s.WriteFile(f) },
			"writing unit %q at %q", unit.Name, f.Path,
		); err != nil {
			return err
		}

		return nil
	}, "processing unit %q", unit.Name)
}

// createPasswd creates the users and groups as described in config.Passwd.
func (s stage) createPasswd(config types.Config) error {
	if err := s.createGroups(config); err != nil {
		return fmt.Errorf("failed to create groups: %v", err)
	}

	if err := s.createUsers(config); err != nil {
		return fmt.Errorf("failed to create users: %v", err)
	}

	return nil
}

// createUsers creates the users as described in config.Passwd.Users.
func (s stage) createUsers(config types.Config) error {
	if len(config.Passwd.Users) == 0 {
		return nil
	}
	s.Logger.PushPrefix("createUsers")
	defer s.Logger.PopPrefix()

	for _, u := range config.Passwd.Users {
		if err := s.CreateUser(u); err != nil {
			return fmt.Errorf("failed to create user %q: %v",
				u.Name, err)
		}

		if err := s.SetPasswordHash(u); err != nil {
			return fmt.Errorf("failed to set password for %q: %v",
				u.Name, err)
		}

		if err := s.AuthorizeSSHKeys(u); err != nil {
			return fmt.Errorf("failed to add keys to user %q: %v",
				u.Name, err)
		}
	}

	return nil
}

// createGroups creates the users as described in config.Passwd.Groups.
func (s stage) createGroups(config types.Config) error {
	if len(config.Passwd.Groups) == 0 {
		return nil
	}
	s.Logger.PushPrefix("createGroups")
	defer s.Logger.PopPrefix()

	for _, g := range config.Passwd.Groups {
		if err := s.CreateGroup(g); err != nil {
			return fmt.Errorf("failed to create group %q: %v",
				g.Name, err)
		}
	}

	return nil
}
