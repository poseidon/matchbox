package rpc

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/coreos/matchbox/matchbox/rpc/rpcpb"
	"github.com/coreos/matchbox/matchbox/server"
)

// NewServer wraps the matchbox Server to return a new gRPC Server.
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
