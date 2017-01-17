package rpc

import (
	"golang.org/x/net/context"

	"github.com/coreos/matchbox/matchbox/rpc/rpcpb"
	"github.com/coreos/matchbox/matchbox/server"
	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
)

// ignitionServer takes a matchbox Server and implements a gRPC IgnitionServer.
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
