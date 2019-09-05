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

package v1_0

import (
	"errors"
	"fmt"
	"reflect"

	base_0_1 "github.com/coreos/fcct/base/v0_1"
	"github.com/coreos/fcct/config/common"
	fcos_0_1 "github.com/coreos/fcct/distro/fcos/v0_1"

	"github.com/coreos/ignition/v2/config/v3_0"
	"github.com/coreos/ignition/v2/config/v3_0/types"
	ignvalidate "github.com/coreos/ignition/v2/config/validate"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/coreos/vcontext/validate"
)

var (
	ErrInvalidConfig = errors.New("config generated was invalid")
)

type Config struct {
	common.Common   `yaml:",inline"`
	base_0_1.Config `yaml:",inline"`
	fcos_0_1.Fcos   `yaml:",inline"`
}

func (c Config) Translate() (types.Config, error) {
	base, err := c.Config.ToIgn3_0()
	if err != nil {
		return types.Config{}, err
	}

	distro, err := c.Fcos.ToIgn3_0()
	if err != nil {
		return types.Config{}, err
	}

	return v3_0.Merge(distro, base), nil
}

func TranslateBytes(input []byte, options common.TranslateOptions) ([]byte, error) {
	cfg := Config{}

	contextTree, err := common.Unmarshal(input, &cfg, options.Strict)
	if err != nil {
		return nil, err
	}

	r := validate.Validate(cfg, "yaml")
	unusedKeyCheck := func(v reflect.Value, c path.ContextPath) report.Report {
		return ignvalidate.ValidateUnusedKeys(v, c, contextTree)
	}
	r.Merge(validate.ValidateCustom(cfg, "yaml", unusedKeyCheck))
	r.Correlate(contextTree)
	if r.IsFatal() {
		fmt.Println(r.String())
		return nil, ErrInvalidConfig
	}

	final, err := cfg.Translate()
	if err != nil {
		return nil, err
	}

	translatedTree := common.ToCamelCase(contextTree)
	second := validate.Validate(final, "json")
	second.Correlate(translatedTree)
	r.Merge(second)
	fmt.Println(r.String())

	if r.IsFatal() {
		return nil, ErrInvalidConfig
	}

	return common.Marshal(final, options.Pretty)
}
