package rpc

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/coreos/coreos-baremetal/bootcfg/rpc/rpcpb"
	"github.com/coreos/coreos-baremetal/bootcfg/server"
)

// NewServer wraps the bootcfg Server to return a new gRPC Server.
func NewServer(s server.Server, tls *tls.Config) *grpc.Server {
	var opts []grpc.ServerOption
	if tls != nil {
		// Add TLS Credentials as a ServerOption for server connections.
		opts = append(opts, grpc.Creds(credentials.NewTLS(tls)))
	}

	grpcServer := grpc.NewServer(opts...)
	rpcpb.RegisterGroupsServer(grpcServer, newGroupServer(s))
	rpcpb.RegisterProfilesServer(grpcServer, newProfileServer(s))
	rpcpb.RegisterSelectServer(grpcServer, newSelectServer(s))
	rpcpb.RegisterIgnitionServer(grpcServer, newIgnitionServer(s))
	return grpcServer
}
