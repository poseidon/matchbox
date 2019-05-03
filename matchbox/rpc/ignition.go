package rpc

import (
	"golang.org/x/net/context"

	"github.com/poseidon/matchbox/matchbox/rpc/rpcpb"
	"github.com/poseidon/matchbox/matchbox/server"
	pb "github.com/poseidon/matchbox/matchbox/server/serverpb"
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

func (s *ignitionServer) IgnitionGet(ctx context.Context, req *pb.IgnitionGetRequest) (*pb.IgnitionGetResponse, error) {
	template, err := s.srv.IgnitionGet(ctx, req)
	return &pb.IgnitionGetResponse{Config: []byte(template)}, grpcError(err)
}

func (s *ignitionServer) IgnitionDelete(ctx context.Context, req *pb.IgnitionDeleteRequest) (*pb.IgnitionDeleteResponse, error) {
	err := s.srv.IgnitionDelete(ctx, req)
	return &pb.IgnitionDeleteResponse{}, grpcError(err)
}
