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
	"time"

	"github.com/coreos/ignition/internal/providers"
)

// WaitUntilOnline waits for the provider to come online. If the provider will
// never be online, or if the timeout elapses before it is online, this returns
// an appropriate error.
func WaitUntilOnline(provider providers.Provider, timeout time.Duration) error {
	online := make(chan bool, 1)
	stop := make(chan struct{})
	defer close(stop)

	go func() {
		for {
			if provider.IsOnline() {
				online <- true
				return
			} else if !provider.ShouldRetry() {
				online <- false
				return
			}

			select {
			case <-time.After(provider.BackoffDuration()):
			case <-stop:
				return
			}
		}
	}()

	expired := make(chan struct{})
	if timeout > 0 {
		go func() {
			<-time.After(timeout)
			close(expired)
		}()
	}

	select {
	case on := <-online:
		if !on {
			return providers.ErrNoProvider
		}
	case <-expired:
		return providers.ErrTimeout
	}

	return nil
}
