package rpc

import (
	"google.golang.org/grpc"

	"github.com/coreos/coreos-baremetal/bootcfg/rpc/rpcpb"
	"github.com/coreos/coreos-baremetal/bootcfg/server"
)

// NewServer wraps the bootcfg Server to return a new gRPC Server.
func NewServer(s server.Server, opts ...grpc.ServerOption) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(opts...)
	rpcpb.RegisterGroupsServer(grpcServer, newGroupServer(s))
	rpcpb.RegisterProfilesServer(grpcServer, newProfileServer(s))
	rpcpb.RegisterSelectServer(grpcServer, newSelectServer(s))
	rpcpb.RegisterIgnitionServer(grpcServer, newIgnitionServer(s))
	return grpcServer, nil
}
