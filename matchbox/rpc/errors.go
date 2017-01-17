package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/coreos/matchbox/matchbox/server"
)

var (
	// work around go vet false positive https://github.com/grpc/grpc-go/issues/90
	grpcErrorf           = grpc.Errorf
	errNoMatchingGroup   = grpcErrorf(codes.NotFound, "matchbox: No matching Group")
	errNoMatchingProfile = grpcErrorf(codes.NotFound, "matchbox: No matching Profile")
)

// grpcError transforms an error into a gRPC errors with canonical error codes.
func grpcError(err error) error {
	if err == nil {
		return err
	}
	switch err {
	case server.ErrNoMatchingGroup:
		return errNoMatchingGroup
	case server.ErrNoMatchingProfile:
		return errNoMatchingProfile
	default:
		return grpcErrorf(codes.Unknown, err.Error())
	}
}
