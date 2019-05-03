package http

import (
	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

var (
	validMACStr = "52:da:00:89:d8:10"

	testProfileIgnitionYAML = &storagepb.Profile{
		Id:         "g1h2i3j4",
		IgnitionId: "ignition.yaml",
	}

	testProfileGeneric = &storagepb.Profile{
		Id:         "g1h2i3j4",
		IgnitionId: "generic.tmpl",
	}

	testGroupWithMAC = &storagepb.Group{
		Id:       "test-group",
		Name:     "test group",
		Profile:  "g1h2i3j4",
		Selector: map[string]string{"mac": validMACStr},
	}
)
