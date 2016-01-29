package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContextSpec(t *testing.T) {
	expectedSpec := &Spec{ID: "g1h2i3j4"}
	ctx := withSpec(context.Background(), expectedSpec)
	spec, err := specFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedSpec, spec)
}

func TestContextSpec_Error(t *testing.T) {
	spec, err := specFromContext(context.Background())
	assert.Nil(t, spec)
	if assert.NotNil(t, err) {
		assert.Equal(t, errNoSpecFromContext, err)
	}
}

func TestGroupSpec(t *testing.T) {
	expectedGroup := &Group{Name: "test group"}
	ctx := withGroup(context.Background(), expectedGroup)
	group, err := groupFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedGroup, group)
}

func TestGroupSpec_Error(t *testing.T) {
	group, err := groupFromContext(context.Background())
	assert.Nil(t, group)
	if assert.NotNil(t, err) {
		assert.Equal(t, errNoGroupFromContext, err)
	}
}
