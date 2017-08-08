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
	RebootStrategy RebootStrategy `yaml:"reboot_strategy" locksmith:"REBOOT_STRATEGY"`
	WindowStart    WindowStart    `yaml:"window_start"    locksmith:"LOCKSMITHD_REBOOT_WINDOW_START"`
	WindowLength   WindowLength   `yaml:"window_length"   locksmith:"LOCKSMITHD_REBOOT_WINDOW_LENGTH"`
	Group          string         `yaml:"group"           locksmith:"LOCKSMITHD_GROUP"`
	EtcdEndpoints  string         `yaml:"etcd_endpoints"  locksmith:"LOCKSMITHD_ENDPOINT"`
	EtcdCAFile     string         `yaml:"etcd_cafile"     locksmith:"LOCKSMITHD_ETCD_CAFILE"`
	EtcdCertFile   string         `yaml:"etcd_certfile"   locksmith:"LOCKSMITHD_ETCD_CERTFILE"`
	EtcdKeyFile    string         `yaml:"etcd_keyfile"    locksmith:"LOCKSMITHD_ETCD_KEYFILE"`
}

func (l Locksmith) configLines() []string {
	return getArgs("%s=%q", "locksmith", l)
}

type RebootStrategy string
type WindowStart string
type WindowLength string

func (l Locksmith) Validate() report.Report {
	if (l.WindowStart != "" && l.WindowLength == "") || (l.WindowStart == "" && l.WindowLength != "") {
		return report.ReportFromError(ErrMissingStartOrLength, report.EntryError)
	}
	return report.Report{}
}

func (r RebootStrategy) Validate() report.Report {
	switch strings.ToLower(string(r)) {
	case "reboot", "etcd-lock", "off":
		return report.Report{}
	default:
		return report.ReportFromError(ErrUnknownStrategy, report.EntryError)
	}
}

func (s WindowStart) Validate() report.Report {
	if s == "" {
		return report.Report{}
	}
	var day string
	var t string

	_, err := fmt.Sscanf(string(s), "%s %s", &day, &t)
	if err != nil {
		day = "not-present"
		t = string(s)
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

func (l WindowLength) Validate() report.Report {
	if l == "" {
		return report.Report{}
	}
	_, err := time.ParseDuration(string(l))
	if err != nil {
		return report.ReportFromError(ErrParsingLength, report.EntryError)
	}
	return report.Report{}
}
