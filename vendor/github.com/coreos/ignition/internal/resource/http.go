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

package resource

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/version"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	maxAttempts    = 15
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 5 * time.Second
)

var (
	ErrAttemptsExhausted = errors.New("unable to fetch resource (no more attempts available)")
)

// HttpClient is a simple wrapper around the Go HTTP client that standardizes
// the process and logging of fetching payloads.
type HttpClient struct {
	client *http.Client
	logger *log.Logger
}

// NewHttpClient creates a new client with the given logger.
func NewHttpClient(logger *log.Logger) HttpClient {
	return HttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 10 * time.Second,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		logger: logger,
	}
}

// getReaderWithHeader performs an HTTP GET on the provided URL with the provided request header
// and returns the response body Reader, HTTP status code, and error (if any). By
// default, User-Agent is added to the header but this can be overridden.
func (c HttpClient) getReaderWithHeader(ctx context.Context, url string, header http.Header) (io.ReadCloser, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "Ignition/"+version.Raw)

	for key, values := range header {
		req.Header.Del(key)
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	duration := initialBackoff
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		c.logger.Debug("GET %s: attempt #%d", url, attempt)
		resp, err := ctxhttp.Do(ctx, c.client, req)

		if err == nil {
			c.logger.Debug("GET result: %s", http.StatusText(resp.StatusCode))
			if resp.StatusCode < 500 {
				return resp.Body, resp.StatusCode, nil
			}
			resp.Body.Close()
		} else {
			c.logger.Debug("GET error: %v", err)
		}

		duration = duration * 2
		if duration > maxBackoff {
			duration = maxBackoff
		}

		select {
		case <-time.After(duration):
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		}
	}

	return nil, 0, ErrAttemptsExhausted
}
