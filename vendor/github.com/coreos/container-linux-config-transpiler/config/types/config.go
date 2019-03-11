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

package types

import (
	"net/url"

	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/coreos/ignition/config/validate/astnode"
	"github.com/coreos/ignition/config/validate/report"
)

type Config struct {
	Ignition  Ignition   `yaml:"ignition"`
	Storage   Storage    `yaml:"storage"`
	Systemd   Systemd    `yaml:"systemd"`
	Networkd  Networkd   `yaml:"networkd"`
	Passwd    Passwd     `yaml:"passwd"`
	Etcd      *Etcd      `yaml:"etcd"`
	Flannel   *Flannel   `yaml:"flannel"`
	Update    *Update    `yaml:"update"`
	Docker    *Docker    `yaml:"docker"`
	Locksmith *Locksmith `yaml:"locksmith"`
}

type Ignition struct {
	Config   IgnitionConfig `yaml:"config"`
	Timeouts Timeouts       `yaml:"timeouts"`
	Security Security       `yaml:"security"`
}

type IgnitionConfig struct {
	Append  []ConfigReference `yaml:"append"`
	Replace *ConfigReference  `yaml:"replace"`
}

type ConfigReference struct {
	Source       string       `yaml:"source"`
	Verification Verification `yaml:"verification"`
}

type Timeouts struct {
	HTTPResponseHeaders *int `yaml:"http_response_headers"`
	HTTPTotal           *int `yaml:"http_total"`
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		r := report.Report{}
		out.Ignition.Timeouts.HTTPResponseHeaders = in.Ignition.Timeouts.HTTPResponseHeaders
		out.Ignition.Timeouts.HTTPTotal = in.Ignition.Timeouts.HTTPTotal
		cfgNode, _ := getNodeChildPath(ast, "ignition", "config", "append")
		for i, ref := range in.Ignition.Config.Append {
			tmp, _ := getNodeChild(cfgNode, i)
			newRef, convertReport := convertConfigReference(ref, tmp)
			r.Merge(convertReport)
			if convertReport.IsFatal() {
				// don't add to the output if invalid
				continue
			}
			out.Ignition.Config.Append = append(out.Ignition.Config.Append, newRef)
		}

		cfgNode, _ = getNodeChildPath(ast, "ignition", "config", "replace")
		if in.Ignition.Config.Replace != nil {
			newRef, convertReport := convertConfigReference(*in.Ignition.Config.Replace, cfgNode)
			r.Merge(convertReport)
			if convertReport.IsFatal() {
				// don't add to the output if invalid
				return out, r, ast
			}
			out.Ignition.Config.Replace = &newRef
		}
		return out, r, ast
	})
}

func convertConfigReference(in ConfigReference, ast astnode.AstNode) (ignTypes.ConfigReference, report.Report) {
	_, err := url.Parse(in.Source)
	if err != nil {
		r := report.ReportFromError(err, report.EntryError)
		if n, err := getNodeChild(ast, "source"); err == nil {
			r.AddPosition(n.ValueLineCol(nil))
		}
		return ignTypes.ConfigReference{}, r
	}

	return ignTypes.ConfigReference{
		Source:       in.Source,
		Verification: convertVerification(in.Verification),
	}, report.Report{}
}

func convertVerification(in Verification) ignTypes.Verification {
	if in.Hash.Function == "" || in.Hash.Sum == "" {
		return ignTypes.Verification{}
	}
	s := in.Hash.String()

	return ignTypes.Verification{
		Hash: &s,
	}
}
