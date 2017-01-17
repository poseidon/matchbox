package rpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"github.com/coreos/matchbox/matchbox/server"
)

func TestGRPCError(t *testing.T) {
	cases := []struct {
		input  error
		output error
	}{
		{nil, nil},
		{server.ErrNoMatchingGroup, errNoMatchingGroup},
		{server.ErrNoMatchingProfile, errNoMatchingProfile},
		{errors.New("other error"), grpcErrorf(codes.Unknown, "other error")},
	}
	for _, c := range cases {
		err := grpcError(c.input)
		assert.Equal(t, c.output, err)
	}
}
