// Copyright 2019 Red Hat, Inc
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
// limitations under the License.)

package config

import (
	"errors"
	"fmt"

	"github.com/coreos/fcct/config/common"
	"github.com/coreos/fcct/config/v1_0"

	"github.com/coreos/go-semver/semver"
	"gopkg.in/yaml.v3"
)

var (
	ErrNoVariant      = errors.New("Error parsing variant. Variant must be specified")
	ErrInvalidVersion = errors.New("Error parsing version. Version must be a valid semver")

	registry = map[string]translator{
		"fcos+1.0.0": v1_0.TranslateBytes,
	}
)

func getTranslator(variant string, version semver.Version) (translator, error) {
	t, ok := registry[fmt.Sprintf("%s+%s", variant, version.String())]
	if !ok {
		return nil, fmt.Errorf("No translator exists for variant %s with version %s", variant, version.String())
	}
	return t, nil
}

type translator func([]byte, common.TranslateOptions) ([]byte, error)

// Translate wraps all of the actual translate functions in a switch that determines the correct one to call
func Translate(input []byte, options common.TranslateOptions) ([]byte, error) {
	// first determine version. This will ignore most fields, so don't use strict
	ver := common.Common{}
	if err := yaml.Unmarshal(input, &ver); err != nil {
		return nil, fmt.Errorf("Error unmarshaling yaml: %v", err)
	}

	if ver.Variant == "" {
		return nil, ErrNoVariant
	}

	tmp, err := semver.NewVersion(ver.Version)
	if err != nil {
		return nil, ErrInvalidVersion
	}
	version := *tmp

	translator, err := getTranslator(ver.Variant, version)
	if err != nil {
		return nil, err
	}

	return translator(input, options)
}
