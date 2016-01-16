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
		assert.Equal(t, "api: Context missing a Spec", err.Error())
	}
}
