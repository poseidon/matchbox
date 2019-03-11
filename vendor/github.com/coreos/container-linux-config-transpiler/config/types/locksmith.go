// Copyright 2017 CoreOS, Inc.
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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/ignition/config/validate/report"
)

var (
	ErrMissingStartOrLength = errors.New("window-start and window-length must both be specified")
	ErrUnknownStrategy      = errors.New("unknown reboot strategy")
	ErrParsingWindowStart   = errors.New("couldn't parse window start")
	ErrUnknownDay           = errors.New("unknown day in window start")
	ErrParsingWindow        = errors.New("couldn't parse window start")
	ErrParsingLength        = errors.New("couldn't parse window length")
)

type Locksmith struct {
	RebootStrategy *string `yaml:"reboot_strategy" locksmith:"REBOOT_STRATEGY"`
	WindowStart    *string `yaml:"window_start"    locksmith:"LOCKSMITHD_REBOOT_WINDOW_START"`
	WindowLength   *string `yaml:"window_length"   locksmith:"LOCKSMITHD_REBOOT_WINDOW_LENGTH"`
	Group          *string `yaml:"group"           locksmith:"LOCKSMITHD_GROUP"`
	EtcdEndpoints  *string `yaml:"etcd_endpoints"  locksmith:"LOCKSMITHD_ENDPOINT"`
	EtcdCAFile     *string `yaml:"etcd_cafile"     locksmith:"LOCKSMITHD_ETCD_CAFILE"`
	EtcdCertFile   *string `yaml:"etcd_certfile"   locksmith:"LOCKSMITHD_ETCD_CERTFILE"`
	EtcdKeyFile    *string `yaml:"etcd_keyfile"    locksmith:"LOCKSMITHD_ETCD_KEYFILE"`
}

func (l Locksmith) configLines() []string {
	return getArgs("%s=%v", "locksmith", l)
}

func nilOrEmpty(s *string) bool {
	return s == nil || *s == ""
}

func (l Locksmith) Validate() report.Report {
	if (!nilOrEmpty(l.WindowStart) && nilOrEmpty(l.WindowLength)) || (nilOrEmpty(l.WindowStart) && !nilOrEmpty(l.WindowLength)) {
		return report.ReportFromError(ErrMissingStartOrLength, report.EntryError)
	}
	return report.Report{}
}

func (l Locksmith) ValidateRebootStrategy() report.Report {
	if nilOrEmpty(l.RebootStrategy) {
		return report.Report{}
	}
	switch strings.ToLower(*l.RebootStrategy) {
	case "reboot", "etcd-lock", "off":
		return report.Report{}
	default:
		return report.ReportFromError(ErrUnknownStrategy, report.EntryError)
	}
}

func (l Locksmith) ValidateWindowStart() report.Report {
	if nilOrEmpty(l.WindowStart) {
		return report.Report{}
	}
	var day string
	var t string

	_, err := fmt.Sscanf(*l.WindowStart, "%s %s", &day, &t)
	if err != nil {
		day = "not-present"
		t = *l.WindowStart
	}

	switch strings.ToLower(day) {
	case "sun", "mon", "tue", "wed", "thu", "fri", "sat", "not-present":
		break
	default:
		return report.ReportFromError(ErrUnknownDay, report.EntryError)
	}

	_, err = time.Parse("15:04", t)
	if err != nil {
		return report.ReportFromError(ErrParsingWindow, report.EntryError)
	}

	return report.Report{}
}

func (l Locksmith) ValidateWindowLength() report.Report {
	if nilOrEmpty(l.WindowLength) {
		return report.Report{}
	}
	_, err := time.ParseDuration(*l.WindowLength)
	if err != nil {
		return report.ReportFromError(ErrParsingLength, report.EntryError)
	}
	return report.Report{}
}
