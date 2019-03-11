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
	"reflect"

	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/coreos/ignition/config/validate/report"

	"github.com/coreos/container-linux-config-transpiler/config/platform"
	"github.com/coreos/container-linux-config-transpiler/config/templating"
	"github.com/coreos/container-linux-config-transpiler/config/types/util"
	iutil "github.com/coreos/container-linux-config-transpiler/internal/util"
	"github.com/coreos/ignition/config/validate/astnode"
)

var (
	ErrPlatformUnspecified = fmt.Errorf("platform must be specified to use templating")
	ErrInvalidKey          = errors.New("Key is invalid (wrong type or not found")
	ErrNilNode             = errors.New("Ast node is nil")
	ErrKeyNotFound         = errors.New("Key not found")
)

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, p string) (ignTypes.Config, report.Report, astnode.AstNode) {
		if p == platform.OpenStackMetadata || p == platform.CloudStackConfigDrive {
			out.Systemd.Units = append(out.Systemd.Units, ignTypes.Unit{
				Name: "coreos-metadata.service",
				Dropins: []ignTypes.SystemdDropin{{
					Name:     "20-clct-provider-override.conf",
					Contents: fmt.Sprintf("[Service]\nEnvironment=COREOS_METADATA_OPT_PROVIDER=--provider=%s", p),
				}},
			})
			out.Systemd.Units = append(out.Systemd.Units, ignTypes.Unit{
				Name:    "coreos-metadata-sshkeys@.service",
				Enabled: iutil.BoolToPtr(true),
				Dropins: []ignTypes.SystemdDropin{{
					Name:     "20-clct-provider-override.conf",
					Contents: fmt.Sprintf("[Service]\nEnvironment=COREOS_METADATA_OPT_PROVIDER=--provider=%s", p),
				}},
			})
		}
		return out, report.Report{}, ast
	})
}

func isZero(v interface{}) bool {
	if v == nil {
		return true
	}
	zv := reflect.Zero(reflect.TypeOf(v))
	return reflect.DeepEqual(v, zv.Interface())
}

// assembleUnit will assemble the contents of a systemd unit dropin that will
// have the given environment variables, and call the given exec line with the
// provided args prepended to it
func assembleUnit(exec string, args, vars []string, p string) (util.SystemdUnit, error) {
	hasTemplating := templating.HasTemplating(args)

	out := util.NewSystemdUnit()
	if hasTemplating {
		if p == "" {
			return util.SystemdUnit{}, ErrPlatformUnspecified
		}
		out.Unit.Add("Requires=coreos-metadata.service")
		out.Unit.Add("After=coreos-metadata.service")
		out.Service.Add("EnvironmentFile=/run/metadata/coreos")
		var err error
		args, err = templating.PerformTemplating(p, args)
		if err != nil {
			return util.SystemdUnit{}, err
		}
	}

	for _, v := range vars {
		out.Service.Add(fmt.Sprintf("Environment=\"%s\"", v))
	}
	for _, a := range args {
		exec += fmt.Sprintf(" \\\n  %s", a)
	}
	out.Service.Add("ExecStart=")
	out.Service.Add(fmt.Sprintf("ExecStart=%s", exec))
	return out, nil
}

func getArgs(format, tagName string, e interface{}) []string {
	if e == nil {
		return nil
	}
	et := reflect.TypeOf(e)
	ev := reflect.ValueOf(e)

	vars := []string{}
	for i := 0; i < et.NumField(); i++ {
		if val := ev.Field(i).Interface(); !isZero(val) {
			if et.Field(i).Anonymous {
				vars = append(vars, getCliArgs(val)...)
			} else {
				key := et.Field(i).Tag.Get(tagName)
				if ev.Field(i).Kind() == reflect.Ptr {
					val = reflect.Indirect(ev.Field(i)).Interface()
				}
				if _, ok := val.(string); ok {
					// to handle whitespace characters
					val = fmt.Sprintf("%q", val)
				}
				vars = append(vars, fmt.Sprintf(format, key, val))
			}
		}
	}

	return vars
}

// getCliArgs builds a list of --ARG=VAL from a struct with cli: tags on its members.
func getCliArgs(e interface{}) []string {
	return getArgs("--%s=%v", "cli", e)
}

// Get returns the value for key, where key is an int or string and n is a
// sequence node or mapping node, respectively.
func getNodeChild(n astnode.AstNode, key interface{}) (astnode.AstNode, error) {
	if n == nil {
		return nil, ErrNilNode
	}
	switch k := key.(type) {
	case int:
		if child, ok := n.SliceChild(k); ok {
			return child, nil
		}
	case string:
		kvmap, ok := n.KeyValueMap()
		if !ok {
			return nil, ErrInvalidKey
		}
		if v, ok := kvmap[k]; ok {
			return v, nil
		}
	default:
		return nil, ErrInvalidKey
	}
	// not found
	return nil, ErrKeyNotFound
}

// GetPath works like calling Get() repeatly with each argument.
func getNodeChildPath(n astnode.AstNode, key ...interface{}) (astnode.AstNode, error) {
	if len(key) == 0 {
		return n, nil
	}
	next, err := getNodeChild(n, key[0])
	if err != nil {
		return nil, err
	}
	return getNodeChildPath(next, key[1:]...)
}
