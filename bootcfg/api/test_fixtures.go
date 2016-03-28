package api

import (
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

var (
	validMACStr = "52:da:00:89:d8:10"

	testProfile = &storagepb.Profile{
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

	testProfileIgnitionYAML = &storagepb.Profile{
		Id:         "g1h2i3j4",
		IgnitionId: "ignition.yaml",
	}

	testGroup = &storagepb.Group{
		Id:           "test-group",
		Name:         "test group",
		Profile:      "g1h2i3j4",
		Requirements: map[string]string{"uuid": "a1b2c3d4"},
		Metadata:     []byte(`{"k8s_version":"v1.1.2","service_name":"etcd2","pod_network": "10.2.0.0/16"}`),
	}

	testGroupWithMAC = &storagepb.Group{
		Id:           "test-group",
		Name:         "test group",
		Profile:      "g1h2i3j4",
		Requirements: map[string]string{"mac": validMACStr},
	}
)
