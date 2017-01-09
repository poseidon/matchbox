package rpc

import (
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/rpc/rpcpb"
	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// ignitionServer takes a bootcfg Server and implements a gRPC IgnitionServer.
type ignitionServer struct {
	srv server.Server
}

func newIgnitionServer(s server.Server) rpcpb.IgnitionServer {
	return &ignitionServer{
		srv: s,
	}
}

func (s *ignitionServer) IgnitionPut(ctx context.Context, req *pb.IgnitionPutRequest) (*pb.IgnitionPutResponse, error) {
	_, err := s.srv.IgnitionPut(ctx, req)
	return &pb.IgnitionPutResponse{}, grpcError(err)
}
