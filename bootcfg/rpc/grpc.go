package rpc

import (
	"google.golang.org/grpc"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// NewServer wraps the bootcfg Server to return a new gRPC Server.
func NewServer(s server.Server, opts ...grpc.ServerOption) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGroupsServer(grpcServer, s)
	pb.RegisterProfilesServer(grpcServer, s)
	pb.RegisterSelectServer(grpcServer, newSelectServer(s))
	return grpcServer, nil
}
