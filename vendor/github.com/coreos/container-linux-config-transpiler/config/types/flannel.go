package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coreos/go-semver/semver"
	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/coreos/ignition/config/validate/astnode"
	"github.com/coreos/ignition/config/validate/report"
)

var (
	ErrFlannelTooOld                  = errors.New("invalid flannel version (too old)")
	ErrFlannelMinorTooNew             = errors.New("flannel minor version too new. Only options available in the previous minor version will be supported")
	ErrNetConfigInvalidJSON           = errors.New("flannel network config doesn't appear to be valid JSON")
	ErrNetConfigProvidedAndKubeMgrSet = errors.New("flannel network config cannot be provided if kube_subnet_mgr is set")
	OldestFlannelVersion              = *semver.New("0.5.0")
	FlannelDefaultVersion             = *semver.New("0.6.0")
)

type Flannel struct {
	Version       *FlannelVersion `yaml:"version"`
	NetworkConfig NetworkConfig   `yaml:"network_config"`
	Options
}

type flannelCommon Flannel

type FlannelVersion semver.Version

type NetworkConfig string

func (nc NetworkConfig) Validate() report.Report {
	if nc == "" {
		return report.Report{}
	}
	tmp := make(map[string]interface{})
	err := json.Unmarshal([]byte(nc), &tmp)
	if err != nil {
		return report.ReportFromError(ErrNetConfigInvalidJSON, report.EntryError)
	}
	return report.Report{}
}

func (v *FlannelVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := semver.Version(*v)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*v = FlannelVersion(t)
	return nil
}

func (fv FlannelVersion) Validate() report.Report {
	v := semver.Version(fv)
	switch {
	case v.LessThan(OldestFlannelVersion):
		return report.ReportFromError(ErrFlannelTooOld, report.EntryError)
	case v.Major == 0 && fv.Minor > 7:
		return report.ReportFromError(ErrFlannelMinorTooNew, report.EntryWarning)
	}
	return report.Report{}
}

func (fv FlannelVersion) String() string {
	return semver.Version(fv).String()
}

func (f *Flannel) Validate() report.Report {
	switch o := f.Options.(type) {
	case Flannel0_7:
		if o.KubeSubnetMgr != nil && *o.KubeSubnetMgr && f.NetworkConfig != "" {
			return report.ReportFromError(ErrNetConfigProvidedAndKubeMgrSet, report.EntryError)
		}
	}
	return report.Report{}
}

func (flannel *Flannel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := flannelCommon(*flannel)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*flannel = Flannel(t)

	var v semver.Version
	if flannel.Version == nil {
		v = FlannelDefaultVersion
	} else {
		v = semver.Version(*flannel.Version)
	}

	if v.Major == 0 && v.Minor >= 7 {
		o := Flannel0_7{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		flannel.Options = o
	} else if v.Major == 0 && v.Minor == 6 {
		o := Flannel0_6{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		flannel.Options = o
	} else if v.Major == 0 && v.Minor == 5 {
		o := Flannel0_5{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		flannel.Options = o
	}
	return nil
}

func init() {
	register(func(in Config, ast astnode.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, astnode.AstNode) {
		if in.Flannel != nil {
			contents, err := flannelContents(*in.Flannel, platform)
			if err != nil {
				return ignTypes.Config{}, report.ReportFromError(err, report.EntryError), ast
			}
			out.Systemd.Units = append(out.Systemd.Units, ignTypes.Unit{
				Name:   "flanneld.service",
				Enable: true,
				Dropins: []ignTypes.SystemdDropin{{
					Name:     "20-clct-flannel.conf",
					Contents: contents,
				}},
			})
		}
		return out, report.Report{}, ast
	})
}

// flannelContents creates the string containing the systemd drop in for flannel
func flannelContents(flannel Flannel, platform string) (string, error) {
	args := getCliArgs(flannel.Options)
	var vars []string
	if flannel.Version != nil {
		vars = []string{fmt.Sprintf("FLANNEL_IMAGE_TAG=v%s", flannel.Version)}
	}

	unit, err := assembleUnit("/usr/lib/coreos/flannel-wrapper $FLANNEL_OPTS", args, vars, platform)
	if err != nil {
		return "", err
	}

	if flannel.NetworkConfig != "" {
		pre := "ExecStartPre=/usr/bin/etcdctl"
		var endpoints *string
		var etcdCAFile *string
		var etcdCertFile *string
		var etcdKeyFile *string
		switch o := flannel.Options.(type) {
		case Flannel0_7:
			endpoints = o.EtcdEndpoints
			etcdCAFile = o.EtcdCAFile
			etcdCertFile = o.EtcdCertFile
			etcdKeyFile = o.EtcdKeyFile
		case Flannel0_6:
			endpoints = o.EtcdEndpoints
			etcdCAFile = o.EtcdCAFile
			etcdCertFile = o.EtcdCertFile
			etcdKeyFile = o.EtcdKeyFile
		case Flannel0_5:
			endpoints = o.EtcdEndpoints
			etcdCAFile = o.EtcdCAFile
			etcdCertFile = o.EtcdCertFile
			etcdKeyFile = o.EtcdKeyFile
		}
		if endpoints != nil {
			pre += fmt.Sprintf(" --endpoints=%q", *endpoints)
		}
		if etcdCAFile != nil {
			pre += fmt.Sprintf(" --ca-file=%q", *etcdCAFile)
		}
		if etcdCertFile != nil {
			pre += fmt.Sprintf(" --cert-file=%q", *etcdCertFile)
		}
		if etcdKeyFile != nil {
			pre += fmt.Sprintf(" --key-file=%q", *etcdKeyFile)
		}
		pre += fmt.Sprintf(" set /coreos.com/network/config %q", flannel.NetworkConfig)
		unit.Service.Add(pre)
	}

	return unit.String(), nil
}

// Flannel0_7 represents flannel options for version 0.7.x. Don't embed Flannel0_6 because
// the yaml parser doesn't handle embedded structs
type Flannel0_7 struct {
	EtcdUsername  *string `yaml:"etcd_username"   cli:"etcd-username"`
	EtcdPassword  *string `yaml:"etcd_password"   cli:"etcd-password"`
	EtcdEndpoints *string `yaml:"etcd_endpoints"  cli:"etcd-endpoints"`
	EtcdCAFile    *string `yaml:"etcd_cafile"     cli:"etcd-cafile"`
	EtcdCertFile  *string `yaml:"etcd_certfile"   cli:"etcd-certfile"`
	EtcdKeyFile   *string `yaml:"etcd_keyfile"    cli:"etcd-keyfile"`
	EtcdPrefix    *string `yaml:"etcd_prefix"     cli:"etcd-prefix"`
	IPMasq        *string `yaml:"ip_masq"         cli:"ip-masq"`
	SubnetFile    *string `yaml:"subnet_file"     cli:"subnet-file"`
	Iface         *string `yaml:"interface"       cli:"iface"`
	PublicIP      *string `yaml:"public_ip"       cli:"public-ip"`
	KubeSubnetMgr *bool   `yaml:"kube_subnet_mgr" cli:"kube-subnet-mgr"`
}

type Flannel0_6 struct {
	EtcdUsername  *string `yaml:"etcd_username"  cli:"etcd-username"`
	EtcdPassword  *string `yaml:"etcd_password"  cli:"etcd-password"`
	EtcdEndpoints *string `yaml:"etcd_endpoints" cli:"etcd-endpoints"`
	EtcdCAFile    *string `yaml:"etcd_cafile"    cli:"etcd-cafile"`
	EtcdCertFile  *string `yaml:"etcd_certfile"  cli:"etcd-certfile"`
	EtcdKeyFile   *string `yaml:"etcd_keyfile"   cli:"etcd-keyfile"`
	EtcdPrefix    *string `yaml:"etcd_prefix"    cli:"etcd-prefix"`
	IPMasq        *string `yaml:"ip_masq"        cli:"ip-masq"`
	SubnetFile    *string `yaml:"subnet_file"    cli:"subnet-file"`
	Iface         *string `yaml:"interface"      cli:"iface"`
	PublicIP      *string `yaml:"public_ip"      cli:"public-ip"`
}

type Flannel0_5 struct {
	EtcdEndpoints *string `yaml:"etcd_endpoints" cli:"etcd-endpoints"`
	EtcdCAFile    *string `yaml:"etcd_cafile"    cli:"etcd-cafile"`
	EtcdCertFile  *string `yaml:"etcd_certfile"  cli:"etcd-certfile"`
	EtcdKeyFile   *string `yaml:"etcd_keyfile"   cli:"etcd-keyfile"`
	EtcdPrefix    *string `yaml:"etcd_prefix"    cli:"etcd-prefix"`
	IPMasq        *string `yaml:"ip_masq"        cli:"ip-masq"`
	SubnetFile    *string `yaml:"subnet_file"    cli:"subnet-file"`
	Iface         *string `yaml:"interface"      cli:"iface"`
	PublicIP      *string `yaml:"public_ip"      cli:"public-ip"`
}
