// Copyright 2018 CoreOS, Inc.
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
	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/coreos/ignition/config/validate/astnode"
	"github.com/coreos/ignition/config/validate/report"
)

type Security struct {
	TLS TLS `yaml:"tls"`
}

type TLS struct {
	CertificateAuthorities []CaReference `yaml:"certificate_authorities"`
}

type CaReference struct {
	Source       string       `yaml:"source"`
	Verification Verification `yaml:"verification"`
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		for _, ca := range in.Ignition.Security.TLS.CertificateAuthorities {
			out.Ignition.Security.TLS.CertificateAuthorities = append(out.Ignition.Security.TLS.CertificateAuthorities, ignTypes.CaReference{
				Source:       ca.Source,
				Verification: convertVerification(ca.Verification),
			})
		}
		return out, report.Report{}, ast
	})
}
