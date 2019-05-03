package rpc

import (
	"golang.org/x/net/context"

	"github.com/poseidon/matchbox/matchbox/rpc/rpcpb"
	"github.com/poseidon/matchbox/matchbox/server"
	pb "github.com/poseidon/matchbox/matchbox/server/serverpb"
)

// profileServer takes a matchbox Server and implements a gRPC ProfilesServer.
type profileServer struct {
	srv server.Server
}

func newProfileServer(s server.Server) rpcpb.ProfilesServer {
	return &profileServer{
		srv: s,
	}
}

func (s *profileServer) ProfilePut(ctx context.Context, req *pb.ProfilePutRequest) (*pb.ProfilePutResponse, error) {
	_, err := s.srv.ProfilePut(ctx, req)
	// TODO(dghubble): Decide on create/put and response(s).
	return &pb.ProfilePutResponse{}, grpcError(err)
}

func (s *profileServer) ProfileGet(ctx context.Context, req *pb.ProfileGetRequest) (*pb.ProfileGetResponse, error) {
	profile, err := s.srv.ProfileGet(ctx, req)
	return &pb.ProfileGetResponse{Profile: profile}, grpcError(err)
}

func (s *profileServer) ProfileDelete(ctx context.Context, req *pb.ProfileDeleteRequest) (*pb.ProfileDeleteResponse, error) {
	err := s.srv.ProfileDelete(ctx, req)
	return &pb.ProfileDeleteResponse{}, grpcError(err)
}

func (s *profileServer) ProfileList(ctx context.Context, req *pb.ProfileListRequest) (*pb.ProfileListResponse, error) {
	profiles, err := s.srv.ProfileList(ctx, req)
	return &pb.ProfileListResponse{Profiles: profiles}, grpcError(err)
}
