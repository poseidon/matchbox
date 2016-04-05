package api

import (
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

var (
	validMACStr = "52:da:00:89:d8:10"

	testProfileIgnitionYAML = &storagepb.Profile{
		Id:         "g1h2i3j4",
		IgnitionId: "ignition.yaml",
	}

	testGroupWithMAC = &storagepb.Group{
		Id:           "test-group",
		Name:         "test group",
		Profile:      "g1h2i3j4",
		Requirements: map[string]string{"mac": validMACStr},
	}
)
