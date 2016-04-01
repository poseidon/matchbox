package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
	"github.com/coreos/coreos-baremetal/bootcfg/storage"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	fake "github.com/coreos/coreos-baremetal/bootcfg/storage/testfakes"
)

func TestSelectGroup(t *testing.T) {
	store := &fake.FixedStore{
		Groups: map[string]*storagepb.Group{fake.Group.Id: fake.Group},
	}
	cases := []struct {
		store         storage.Store
		labels        map[string]string
		expectedGroup *storagepb.Group
		expectedErr   error
	}{
		{store, map[string]string{"uuid": "a1b2c3d4"}, fake.Group, nil},
		{store, nil, nil, ErrNoMatchingGroup},
		// no groups in the store
		{&fake.EmptyStore{}, map[string]string{"a": "b"}, nil, ErrNoMatchingGroup},
	}
	for _, c := range cases {
		srv := NewServer(&Config{c.store})
		group, err := srv.SelectGroup(context.Background(), &pb.SelectGroupRequest{Labels: c.labels})
		if assert.Equal(t, c.expectedErr, err) {
			assert.Equal(t, c.expectedGroup, group)
		}
	}
}
