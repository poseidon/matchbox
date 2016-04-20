package testfakes

import (
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

var (
	// Group is a machine group for testing.
	Group = &storagepb.Group{
		Id:       "test-group",
		Name:     "test group",
		Profile:  "g1h2i3j4",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
		Metadata: []byte(`{"k8s_version":"v1.1.2","pod_network":"10.2.0.0/16","service_name":"etcd2"}`),
	}

	// GroupNoMetadata is a Group without any metadata.
	GroupNoMetadata = &storagepb.Group{
		Id:       "group-no-metadata",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
		Metadata: nil,
	}

	// Profile is a machine profile for testing.
	Profile = &storagepb.Profile{
		Id: "g1h2i3j4",
		Boot: &storagepb.NetBoot{
			Kernel: "/image/kernel",
			Initrd: []string{"/image/initrd_a", "/image/initrd_b"},
			Cmdline: map[string]string{
				"a": "b",
				"c": "",
			},
		},
		CloudId:    "cloud-config.yml",
		IgnitionId: "ignition.json",
	}

	// IgnitionYAMLName is an Ignition template name for testing.
	IgnitionYAMLName = "ignition.yaml"

	// IgnitionYAML is an Ignition template for testing.
	IgnitionYAML = `ignition_version: 1
systemd:
  units:
    - name: etcd2.service
      enable: true
`
)
