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

package util

import (
	"testing"
	"time"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/providers"
)

type mockProvider struct {
	config  types.Config
	err     error
	online  bool
	retry   bool
	backoff time.Duration
}

func (p mockProvider) FetchConfig() (types.Config, error) { return p.config, p.err }
func (p mockProvider) IsOnline() bool                     { return p.online }
func (p mockProvider) ShouldRetry() bool                  { return p.retry }
func (p mockProvider) BackoffDuration() time.Duration     { return p.backoff }

func TestWaitUntilOnline(t *testing.T) {
	type in struct {
		provider mockProvider
		timeout  time.Duration
	}
	type out struct {
		err error
	}

	online := mockProvider{online: true}
	offline := mockProvider{online: false}
	offlineRetry := mockProvider{online: false, retry: true}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{provider: online, timeout: time.Second},
			out: out{err: nil},
		},
		{
			in:  in{provider: offline, timeout: time.Second},
			out: out{err: providers.ErrNoProvider},
		},
		{
			in:  in{provider: offlineRetry, timeout: time.Second},
			out: out{err: providers.ErrTimeout},
		},
	}

	for i, test := range tests {
		err := WaitUntilOnline(test.in.provider, test.in.timeout)
		if test.out.err != err {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
	}
}
