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

package config

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/coreos/ignition/third_party/github.com/go-yaml/yaml"
)

func TestDevicePathUnmarshalJSON(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		device DevicePath
		err    error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: `"/path"`},
			out: out{device: DevicePath("/path")},
		},
		{
			in:  in{data: `"bad"`},
			out: out{device: DevicePath("bad"), err: errors.New("device path not absolute")},
		},
	}

	for i, test := range tests {
		var device DevicePath
		err := json.Unmarshal([]byte(test.in.data), &device)
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
		if !reflect.DeepEqual(test.out.device, device) {
			t.Errorf("#%d: bad device: want %#v, got %#v", i, test.out.device, device)
		}
	}
}

func TestDevicePathUnmarshalYAML(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		device DevicePath
		err    error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: `"/path"`},
			out: out{device: DevicePath("/path")},
		},
		{
			in:  in{data: `"bad"`},
			out: out{device: DevicePath("bad"), err: errors.New("device path not absolute")},
		},
	}

	for i, test := range tests {
		var device DevicePath
		err := yaml.Unmarshal([]byte(test.in.data), &device)
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
		if !reflect.DeepEqual(test.out.device, device) {
			t.Errorf("#%d: bad device: want %#v, got %#v", i, test.out.device, device)
		}
	}
}

func TestDevicePathAssertValid(t *testing.T) {
	type in struct {
		device DevicePath
	}
	type out struct {
		err error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{device: DevicePath("/good/path")},
			out: out{},
		},
		{
			in:  in{device: DevicePath("/name")},
			out: out{},
		},
		{
			in:  in{device: DevicePath("/this/is/a/fairly/long/path/to/a/device.")},
			out: out{},
		},
		{
			in:  in{device: DevicePath("/this one has spaces")},
			out: out{},
		},
		{
			in:  in{device: DevicePath("relative/path")},
			out: out{err: errors.New("device path not absolute")},
		},
	}

	for i, test := range tests {
		err := test.in.device.assertValid()
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
	}
}

func TestFilesystemFormatUnmarshalJSON(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		format FilesystemFormat
		err    error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: `"ext4"`},
			out: out{format: FilesystemFormat("ext4")},
		},
		{
			in:  in{data: `"bad"`},
			out: out{format: FilesystemFormat("bad"), err: errors.New("invalid filesystem format")},
		},
	}

	for i, test := range tests {
		var format FilesystemFormat
		err := json.Unmarshal([]byte(test.in.data), &format)
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
		if !reflect.DeepEqual(test.out.format, format) {
			t.Errorf("#%d: bad format: want %#v, got %#v", i, test.out.format, format)
		}
	}
}

func TestFilesystemFormatUnmarshalYAML(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		format FilesystemFormat
		err    error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: `"ext4"`},
			out: out{format: FilesystemFormat("ext4")},
		},
		{
			in:  in{data: `"bad"`},
			out: out{format: FilesystemFormat("bad"), err: errors.New("invalid filesystem format")},
		},
	}

	for i, test := range tests {
		var format FilesystemFormat
		err := yaml.Unmarshal([]byte(test.in.data), &format)
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
		if !reflect.DeepEqual(test.out.format, format) {
			t.Errorf("#%d: bad device: want %#v, got %#v", i, test.out.format, format)
		}
	}
}

func TestFilesystemFormatAssertValid(t *testing.T) {
	type in struct {
		format FilesystemFormat
	}
	type out struct {
		err error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{format: FilesystemFormat("ext4")},
			out: out{},
		},
		{
			in:  in{format: FilesystemFormat("btrfs")},
			out: out{},
		},
		{
			in:  in{format: FilesystemFormat("")},
			out: out{err: errors.New("invalid filesystem format")},
		},
	}

	for i, test := range tests {
		err := test.in.format.assertValid()
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
	}
}

func TestMkfsOptionsUnmarshalJSON(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		options MkfsOptions
		err     error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: `["--label=ROOT"]`},
			out: out{options: MkfsOptions([]string{"--label=ROOT"})},
		},
	}

	for i, test := range tests {
		var options MkfsOptions
		err := json.Unmarshal([]byte(test.in.data), &options)
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
		if !reflect.DeepEqual(test.out.options, options) {
			t.Errorf("#%d: bad format: want %#v, got %#v", i, test.out.options, options)
		}
	}
}

func TestMkfsOptionsUnmarshalYAML(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		options MkfsOptions
		err     error
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{data: `["--label=ROOT"]`},
			out: out{options: MkfsOptions([]string{"--label=ROOT"})},
		},
	}

	for i, test := range tests {
		var options MkfsOptions
		err := yaml.Unmarshal([]byte(test.in.data), &options)
		if !reflect.DeepEqual(test.out.err, err) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.err, err)
		}
		if !reflect.DeepEqual(test.out.options, options) {
			t.Errorf("#%d: bad device: want %#v, got %#v", i, test.out.options, options)
		}
	}
}
