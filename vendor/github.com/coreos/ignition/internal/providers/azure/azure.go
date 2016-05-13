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

// The azure provider fetches a configuration from the Azure OVF DVD.

package azure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/coreos/ignition/config"
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/providers"
	"github.com/coreos/ignition/internal/providers/util"
)

const (
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 30 * time.Second
	configDevice   = "/dev/disk/by-id/ata-Virtual_CD"
	configPath     = "/CustomData.bin"
)

// These constants come from <cdrom.h>.
const (
	CDROM_DRIVE_STATUS = 0x5326
)

// These constants come from <cdrom.h>.
const (
	CDS_NO_INFO = iota
	CDS_NO_DISC
	CDS_TRAY_OPEN
	CDS_DRIVE_NOT_READY
	CDS_DISC_OK
)

type Creator struct{}

func (Creator) Create(logger *log.Logger) providers.Provider {
	return &provider{
		logger:  logger,
		backoff: initialBackoff,
	}
}

type provider struct {
	logger  *log.Logger
	backoff time.Duration
}

func (p provider) FetchConfig() (types.Config, error) {
	p.logger.Debug("creating temporary mount point")
	mnt, err := ioutil.TempDir("", "ignition-azure")
	if err != nil {
		return types.Config{}, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.Remove(mnt)

	p.logger.Debug("mounting config device")
	if err := p.logger.LogOp(
		func() error { return syscall.Mount(configDevice, mnt, "udf", syscall.MS_RDONLY, "") },
		"mounting %q at %q", configDevice, mnt,
	); err != nil {
		return types.Config{}, fmt.Errorf("failed to mount device %q at %q: %v", configDevice, mnt, err)
	}
	defer p.logger.LogOp(
		func() error { return syscall.Unmount(mnt, 0) },
		"unmounting %q at %q", configDevice, mnt,
	)

	p.logger.Debug("reading config")
	rawConfig, err := ioutil.ReadFile(filepath.Join(mnt, configPath))
	if err != nil && !os.IsNotExist(err) {
		return types.Config{}, fmt.Errorf("failed to read config: %v", err)
	}

	return config.Parse(rawConfig)
}

func (p provider) IsOnline() bool {
	p.logger.Debug("opening config device")
	device, err := os.Open(configDevice)
	if err != nil {
		p.logger.Info("failed to open config device: %v", err)
		return false
	}
	defer device.Close()

	p.logger.Debug("getting drive status")
	status, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(device.Fd()),
		uintptr(CDROM_DRIVE_STATUS),
		uintptr(0),
	)

	switch status {
	case CDS_NO_INFO:
		p.logger.Info("drive status: no info")
	case CDS_NO_DISC:
		p.logger.Info("drive status: no disc")
	case CDS_TRAY_OPEN:
		p.logger.Info("drive status: open")
	case CDS_DRIVE_NOT_READY:
		p.logger.Info("drive status: not ready")
	case CDS_DISC_OK:
		p.logger.Info("drive status: OK")
	default:
		p.logger.Err("failed to get drive status: %s", errno.Error())
	}

	return (status == CDS_DISC_OK)
}

func (p provider) ShouldRetry() bool {
	return true
}

func (p *provider) BackoffDuration() time.Duration {
	return util.ExpBackoff(&p.backoff, maxBackoff)
}
