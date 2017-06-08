package testfakes

import (
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

var (
	// Group is a machine group for testing.
	Group = &storagepb.Group{
		Id:       "test-group",
		Name:     "test group",
		Profile:  "g1h2i3j4",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
		Metadata: []byte(`{"pod_network":"10.2.0.0/16","service_name":"etcd2"}`),
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
			Args: []string{
				"a=b",
				"c",
			},
		},
		CloudId:    "cloud-config.tmpl",
		IgnitionId: "ignition.tmpl",
		GenericId:  "generic.tmpl",
	}

	// IgnitionYAMLName is an Ignition template name for testing.
	IgnitionYAMLName = "ignition.tmpl"

	// IgnitionYAML is an Ignition template for testing.
	IgnitionYAML = `ignition_version: 1
systemd:
  units:
    - name: etcd2.service
      enable: true
`

	// GenericName is a Generic template name for testing.
	GenericName = "generic.tmpl"

	// Generic is a Generic template for testing.
	Generic = `
This is a generic template.
`
)
