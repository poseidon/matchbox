package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_MissingEndpoints(t *testing.T) {
	cfg := &Config{
		Endpoints: []string{},
	}
	client, err := New(cfg)
	assert.Nil(t, client)
	assert.Equal(t, errNoEndpoints, err)
}
