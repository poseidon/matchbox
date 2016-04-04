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
		store  storage.Store
		labels map[string]string
		group  *storagepb.Group
		err    error
	}{
		{store, map[string]string{"uuid": "a1b2c3d4"}, fake.Group, nil},
		// no labels provided
		{store, nil, nil, ErrNoMatchingGroup},
		// empty store
		{&fake.EmptyStore{}, map[string]string{"a": "b"}, nil, ErrNoMatchingGroup},
	}
	for _, c := range cases {
		srv := NewServer(&Config{c.store})
		group, err := srv.SelectGroup(context.Background(), &pb.SelectGroupRequest{Labels: c.labels})
		if assert.Equal(t, c.err, err) {
			assert.Equal(t, c.group, group)
		}
	}
}

func TestSelectProfile(t *testing.T) {
	store := &fake.FixedStore{
		Groups:   map[string]*storagepb.Group{fake.Group.Id: fake.Group},
		Profiles: map[string]*storagepb.Profile{fake.Group.Profile: fake.Profile},
	}
	missingProfileStore := &fake.FixedStore{
		Groups: map[string]*storagepb.Group{fake.Group.Id: fake.Group},
	}
	cases := []struct {
		store   storage.Store
		labels  map[string]string
		profile *storagepb.Profile
		err     error
	}{
		{store, map[string]string{"uuid": "a1b2c3d4"}, fake.Profile, nil},
		// matching group, but missing profile
		{missingProfileStore, map[string]string{"uuid": "a1b2c3d4"}, nil, ErrNoMatchingProfile},
		// no labels provided
		{store, nil, nil, ErrNoMatchingGroup},
		// empty store
		{&fake.EmptyStore{}, map[string]string{"a": "b"}, nil, ErrNoMatchingGroup},
	}
	for _, c := range cases {
		srv := NewServer(&Config{c.store})
		profile, err := srv.SelectProfile(context.Background(), &pb.SelectProfileRequest{Labels: c.labels})
		if assert.Equal(t, c.err, err) {
			assert.Equal(t, c.profile, profile)
		}
	}
}

func TestProfilePut(t *testing.T) {
	srv := NewServer(&Config{Store: fake.NewFixedStore()})
	_, err := srv.ProfilePut(context.Background(), &pb.ProfilePutRequest{Profile: fake.Profile})
	assert.Nil(t, err)
}

func TestProfileGet(t *testing.T) {
	store := &fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{fake.Profile.Id: fake.Profile},
	}
	cases := []struct {
		id      string
		profile *storagepb.Profile
		err     error
	}{
		{fake.Profile.Id, fake.Profile, nil},
	}
	srv := NewServer(&Config{store})
	for _, c := range cases {
		profile, err := srv.ProfileGet(context.Background(), &pb.ProfileGetRequest{Id: c.id})
		assert.Equal(t, c.err, err)
		assert.Equal(t, c.profile, profile)
	}
}

func TestProfileList(t *testing.T) {
	store := &fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{fake.Profile.Id: fake.Profile},
	}
	srv := NewServer(&Config{store})
	profiles, err := srv.ProfileList(context.Background(), &pb.ProfileListRequest{})
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(profiles)) {
		assert.Equal(t, fake.Profile, profiles[0])
	}
}

func TestProfileList_Empty(t *testing.T) {
	srv := NewServer(&Config{&fake.EmptyStore{}})
	profiles, err := srv.ProfileList(context.Background(), &pb.ProfileListRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(profiles))
}
