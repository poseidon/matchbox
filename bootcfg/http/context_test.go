package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/mikeynap/coreos-baremetal/bootcfg/storage/storagepb"
)

func TestContextProfile(t *testing.T) {
	expectedProfile := &storagepb.Profile{Id: "g1h2i3j4"}
	ctx := withProfile(context.Background(), expectedProfile)
	profile, err := profileFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedProfile, profile)
}

func TestContextProfile_Error(t *testing.T) {
	profile, err := profileFromContext(context.Background())
	assert.Nil(t, profile)
	if assert.NotNil(t, err) {
		assert.Equal(t, errNoProfileFromContext, err)
	}
}

func TestContextGroup(t *testing.T) {
	expectedGroup := &storagepb.Group{Name: "test group"}
	ctx := withGroup(context.Background(), expectedGroup)
	group, err := groupFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedGroup, group)
}

func TestContextGroup_Error(t *testing.T) {
	group, err := groupFromContext(context.Background())
	assert.Nil(t, group)
	if assert.NotNil(t, err) {
		assert.Equal(t, errNoGroupFromContext, err)
	}
}
