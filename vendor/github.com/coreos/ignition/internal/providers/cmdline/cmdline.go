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

// The cmdline provider fetches a remote configuration from the URL specified
// in the kernel boot option "coreos.config.url".

package cmdline

import (
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/ignition/config"
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/providers"
	putil "github.com/coreos/ignition/internal/providers/util"
	"github.com/coreos/ignition/internal/util"
)

const (
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 30 * time.Second
	cmdlinePath    = "/proc/cmdline"
	cmdlineUrlFlag = "coreos.config.url"
)

type Creator struct{}

func (Creator) Create(logger *log.Logger) providers.Provider {
	return &provider{
		logger:  logger,
		backoff: initialBackoff,
		path:    cmdlinePath,
		client:  util.NewHttpClient(logger),
	}
}

type provider struct {
	logger    *log.Logger
	backoff   time.Duration
	path      string
	client    util.HttpClient
	configUrl string
	rawConfig []byte
}

func (p provider) FetchConfig() (types.Config, error) {
	if p.rawConfig == nil {
		return types.Config{}, nil
	} else {
		return config.Parse(p.rawConfig)
	}
}

func (p *provider) IsOnline() bool {
	if p.configUrl == "" {
		args, err := ioutil.ReadFile(p.path)
		if err != nil {
			p.logger.Err("couldn't read cmdline")
			return false
		}

		p.configUrl = parseCmdline(args)
		p.logger.Debug("parsed url from cmdline: %q", p.configUrl)
		if p.configUrl == "" {
			// If the cmdline flag wasn't provided, just no-op.
			p.logger.Info("no config URL provided")
			return true
		}
	}

	return p.getRawConfig()

}

func (p provider) ShouldRetry() bool {
	return true
}

func (p *provider) BackoffDuration() time.Duration {
	return putil.ExpBackoff(&p.backoff, maxBackoff)
}

func parseCmdline(cmdline []byte) (url string) {
	for _, arg := range strings.Split(string(cmdline), " ") {
		parts := strings.SplitN(strings.TrimSpace(arg), "=", 2)
		key := parts[0]

		if key != cmdlineUrlFlag {
			continue
		}

		if len(parts) == 2 {
			url = parts[1]
		}
	}

	return
}

// getRawConfig gets the raw configuration data from p.configUrl.
// Supported URL schemes are:
// http://	remote resource accessed via http
// oem://	local file in /usr/share/oem or /mnt/oem
func (p *provider) getRawConfig() bool {
	url, err := url.Parse(p.configUrl)
	if err != nil {
		p.logger.Err("failed to parse url: %v", err)
		return false
	}

	data, err := util.FetchResource(p.logger, *url)
	if err != nil {
		p.logger.Err("failed to fetch %v: %v", url, err)
		return false
	}

	p.rawConfig = data
	return true
}
