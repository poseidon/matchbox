package server

import (
	"testing"

	"context"
	"github.com/stretchr/testify/assert"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
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

func TestGroupCRUD(t *testing.T) {
	srv := NewServer(&Config{Store: fake.NewFixedStore()})
	_, err := srv.GroupPut(context.Background(), &pb.GroupPutRequest{Group: fake.Group})
	// assert that:
	// - Group creation is successful
	// - Group can be retrieved by id
	// - Group can be deleted by id
	assert.Nil(t, err)

	group, err := srv.GroupGet(context.Background(), &pb.GroupGetRequest{Id: fake.Group.Id})
	assert.Nil(t, err)
	assert.Equal(t, fake.Group, group)

	err = srv.GroupDelete(context.Background(), &pb.GroupDeleteRequest{Id: fake.Group.Id})
	assert.Nil(t, err)
	_, err = srv.GroupGet(context.Background(), &pb.GroupGetRequest{Id: fake.Group.Id})
	assert.Error(t, err)
}

func TestGroupCreate_Invalid(t *testing.T) {
	srv := NewServer(&Config{Store: fake.NewFixedStore()})
	invalid := &storagepb.Group{}
	_, err := srv.GroupPut(context.Background(), &pb.GroupPutRequest{Group: invalid})
	assert.Error(t, err)
}

func TestGroupList(t *testing.T) {
	store := &fake.FixedStore{
		Groups: map[string]*storagepb.Group{fake.Group.Id: fake.Group},
	}
	srv := NewServer(&Config{store})
	groups, err := srv.GroupList(context.Background(), &pb.GroupListRequest{})
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(groups)) {
		assert.Equal(t, fake.Group, groups[0])
	}
}

func TestGroup_BrokenStore(t *testing.T) {
	srv := NewServer(&Config{&fake.BrokenStore{}})
	_, err := srv.GroupPut(context.Background(), &pb.GroupPutRequest{Group: fake.Group})
	assert.Error(t, err)
	_, err = srv.GroupGet(context.Background(), &pb.GroupGetRequest{Id: fake.Group.Id})
	assert.Error(t, err)
	err = srv.GroupDelete(context.Background(), &pb.GroupDeleteRequest{Id: fake.Group.Id})
	assert.Error(t, err)
	_, err = srv.GroupList(context.Background(), &pb.GroupListRequest{})
	assert.Error(t, err)
}

func TestProfileCRUD(t *testing.T) {
	srv := NewServer(&Config{Store: fake.NewFixedStore()})
	_, err := srv.ProfilePut(context.Background(), &pb.ProfilePutRequest{Profile: fake.Profile})
	// assert that:
	// - Profile creation is successful
	// - Profile can be retrieved by id
	// - Profile can be deleted by id
	assert.Nil(t, err)

	profile, err := srv.ProfileGet(context.Background(), &pb.ProfileGetRequest{Id: fake.Profile.Id})
	assert.Equal(t, fake.Profile, profile)
	assert.Nil(t, err)

	err = srv.ProfileDelete(context.Background(), &pb.ProfileDeleteRequest{Id: fake.Group.Id})
	assert.Nil(t, err)
	_, err = srv.ProfileGet(context.Background(), &pb.ProfileGetRequest{Id: fake.Group.Id})
	assert.Error(t, err)
}

func TestProfileCreate_Invalid(t *testing.T) {
	srv := NewServer(&Config{Store: fake.NewFixedStore()})
	invalid := &storagepb.Profile{}
	_, err := srv.ProfilePut(context.Background(), &pb.ProfilePutRequest{Profile: invalid})
	assert.Error(t, err)
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

func TestProfiles_BrokenStore(t *testing.T) {
	srv := NewServer(&Config{&fake.BrokenStore{}})
	_, err := srv.ProfilePut(context.Background(), &pb.ProfilePutRequest{Profile: fake.Profile})
	assert.Error(t, err)
	_, err = srv.ProfileGet(context.Background(), &pb.ProfileGetRequest{Id: fake.Profile.Id})
	assert.Error(t, err)
	err = srv.ProfileDelete(context.Background(), &pb.ProfileDeleteRequest{Id: fake.Profile.Id})
	assert.Error(t, err)
	_, err = srv.ProfileList(context.Background(), &pb.ProfileListRequest{})
	assert.Error(t, err)
}

func TestIgnitionCRUD(t *testing.T) {
	srv := NewServer(&Config{Store: fake.NewFixedStore()})
	req := &pb.IgnitionPutRequest{
		Name:   fake.IgnitionYAMLName,
		Config: []byte(fake.IgnitionYAML),
	}
	_, err := srv.IgnitionPut(context.Background(), req)
	// assert that:
	// - Ignition template creation is successful
	// - Ignition template can be retrieved by name
	// - Ignition template can be deleted by name
	assert.Nil(t, err)
	template, err := srv.IgnitionGet(context.Background(), &pb.IgnitionGetRequest{Name: fake.IgnitionYAMLName})
	assert.Equal(t, fake.IgnitionYAML, template)
	assert.Nil(t, err)

	err = srv.IgnitionDelete(context.Background(), &pb.IgnitionDeleteRequest{Name: fake.IgnitionYAMLName})
	assert.Nil(t, err)
	_, err = srv.IgnitionGet(context.Background(), &pb.IgnitionGetRequest{Name: fake.IgnitionYAMLName})
	assert.Error(t, err)
}

func TestIgnition_BrokenStore(t *testing.T) {
	srv := NewServer(&Config{&fake.BrokenStore{}})
	req := &pb.IgnitionPutRequest{
		Name:   fake.IgnitionYAMLName,
		Config: []byte(fake.IgnitionYAML),
	}
	_, err := srv.IgnitionPut(context.Background(), req)
	assert.Error(t, err)
	_, err = srv.IgnitionGet(context.Background(), &pb.IgnitionGetRequest{Name: fake.IgnitionYAMLName})
	assert.Error(t, err)
	err = srv.IgnitionDelete(context.Background(), &pb.IgnitionDeleteRequest{Name: fake.IgnitionYAMLName})
	assert.Error(t, err)
}

func TestGenericCRUD(t *testing.T) {
	srv := NewServer(&Config{Store: fake.NewFixedStore()})
	req := &pb.GenericPutRequest{
		Name:   fake.GenericName,
		Config: []byte(fake.Generic),
	}
	_, err := srv.GenericPut(context.Background(), req)
	// assert that:
	// - Generic template creation is successful
	// - Generic template can be retrieved by name
	// - Generic template can be deleted by name
	assert.Nil(t, err)
	template, err := srv.GenericGet(context.Background(), &pb.GenericGetRequest{Name: fake.GenericName})
	assert.Equal(t, fake.Generic, template)
	assert.Nil(t, err)

	err = srv.GenericDelete(context.Background(), &pb.GenericDeleteRequest{Name: fake.GenericName})
	assert.Nil(t, err)
	_, err = srv.GenericGet(context.Background(), &pb.GenericGetRequest{Name: fake.GenericName})
	assert.Error(t, err)
}

func TestGeneric_BrokenStore(t *testing.T) {
	srv := NewServer(&Config{&fake.BrokenStore{}})
	req := &pb.GenericPutRequest{
		Name:   fake.GenericName,
		Config: []byte(fake.Generic),
	}
	_, err := srv.GenericPut(context.Background(), req)
	assert.Error(t, err)
	_, err = srv.GenericGet(context.Background(), &pb.GenericGetRequest{Name: fake.GenericName})
	assert.Error(t, err)
	err = srv.GenericDelete(context.Background(), &pb.GenericDeleteRequest{Name: fake.GenericName})

	assert.Error(t, err)
}
