package http

import (
	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

var (
	validMACStr = "52:da:00:89:d8:10"

	testProfileWithButane = &storagepb.Profile{
		Id:         "g1h2i3j4",
		IgnitionId: "butane.yaml",
	}
)
