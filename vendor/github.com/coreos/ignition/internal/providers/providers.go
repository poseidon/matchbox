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

package providers

import (
	"errors"
	"time"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/log"
)

var (
	ErrNoProvider = errors.New("config provider was not online")
	ErrTimeout    = errors.New("timed out while waiting for config provider to come online")
)

// Provider represents an external source of configuration. The source can be
// local to the host system or it may be remote. The provider dictates whether
// or not the source is online, if the caller should try again when the source
// is offline, and how long the caller should wait before retries.
type Provider interface {
	FetchConfig() (types.Config, error)
	IsOnline() bool
	ShouldRetry() bool
	BackoffDuration() time.Duration
}

type ProviderCreator interface {
	Create(logger *log.Logger) Provider
}
