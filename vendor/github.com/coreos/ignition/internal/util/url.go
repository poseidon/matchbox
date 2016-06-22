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

package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"

	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/systemd"

	"github.com/vincent-petithory/dataurl"
)

var (
	ErrSchemeUnsupported = errors.New("unsupported source scheme")
	ErrPathNotAbsolute   = errors.New("path is not absolute")
	ErrNotFound          = errors.New("resource not found")
	ErrFailed            = errors.New("failed to fetch resource")
)

const (
	oemDevicePath = "/dev/disk/by-label/OEM" // Device link where oem partition is found.
	oemDirPath    = "/usr/share/oem"         // OEM dir within root fs to consider for pxe scenarios.
	oemMountPath  = "/mnt/oem"               // Mountpoint where oem partition is mounted when present.
)

// FetchResource fetches a resource given a URL. The supported schemes are http, data, and oem.
func FetchResource(l *log.Logger, u url.URL) ([]byte, error) {
	switch u.Scheme {
	case "http", "https":
		client := NewHttpClient(l)
		data, status, err := client.Get(u.String())
		if err != nil {
			return nil, err
		}

		l.Debug("GET result: %s", http.StatusText(status))
		switch status {
		case http.StatusOK, http.StatusNoContent:
			return data, nil
		case http.StatusNotFound:
			return nil, ErrNotFound
		default:
			return nil, ErrFailed
		}

	case "data":
		url, err := dataurl.DecodeString(u.String())
		if err != nil {
			return nil, err
		}

		return url.Data, nil

	case "oem":
		path := filepath.Clean(u.Path)
		if !filepath.IsAbs(path) {
			l.Err("oem path is not absolute: %q", u.Path)
			return nil, ErrPathNotAbsolute
		}

		// check if present under oemDirPath, if so use it.
		absPath := filepath.Join(oemDirPath, path)
		data, err := ioutil.ReadFile(absPath)
		if os.IsNotExist(err) {
			l.Info("oem config not found in %q, trying %q",
				oemDirPath, oemMountPath)

			// try oemMountPath, requires mounting it.
			err := mountOEM(l)
			if err != nil {
				l.Err("failed to mount oem partition: %v", err)
				return nil, ErrFailed
			}

			absPath := filepath.Join(oemMountPath, path)
			data, err = ioutil.ReadFile(absPath)
			umountOEM(l)
		} else if err != nil {
			l.Err("failed to read oem config: %v", err)
			return nil, ErrFailed
		}

		return data, nil

	default:
		return nil, ErrSchemeUnsupported
	}
}

// mountOEM waits for the presence of and mounts the oem partition at oemMountPath.
func mountOEM(l *log.Logger) error {
	dev := []string{oemDevicePath}
	if err := systemd.WaitOnDevices(dev, "oem-cmdline"); err != nil {
		l.Err("failed to wait for oem device: %v", err)
		return err
	}

	if err := os.MkdirAll(oemMountPath, 0700); err != nil {
		l.Err("failed to create oem mount point: %v", err)
		return err
	}

	if err := l.LogOp(
		func() error {
			return syscall.Mount(dev[0], oemMountPath, "ext4", 0, "")
		},
		"mounting %q at %q", oemDevicePath, oemMountPath,
	); err != nil {
		return fmt.Errorf("failed to mount device %q at %q: %v",
			oemDevicePath, oemMountPath, err)
	}

	return nil
}

// umountOEM unmounts the oem partition at oemMountPath.
func umountOEM(l *log.Logger) {
	l.LogOp(
		func() error { return syscall.Unmount(oemMountPath, 0) },
		"unmounting %q", oemMountPath,
	)
}
