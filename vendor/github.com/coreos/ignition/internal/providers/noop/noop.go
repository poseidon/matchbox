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

// The noop provider does nothing, for use by unimplemented oems.

package noop

import (
	"time"

	"github.com/coreos/ignition/config"
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/providers"
)

type Creator struct{}

func (Creator) Create(logger *log.Logger) providers.Provider {
	return &provider{
		logger: logger,
	}
}

type provider struct {
	logger *log.Logger
}

func (p provider) FetchConfig() (types.Config, error) {
	p.logger.Debug("noop provider fetching empty config")
	return types.Config{}, config.ErrEmpty
}

func (p *provider) IsOnline() bool {
	return true
}

func (p provider) ShouldRetry() bool {
	return false
}

func (p *provider) BackoffDuration() time.Duration {
	return 0
}
