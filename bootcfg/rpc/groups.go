package rpc

import (
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// groupServer takes a bootcfg Server and implements a gRPC GroupsServer.
type groupServer struct {
	srv server.Server
}

func newGroupServer(s server.Server) pb.GroupsServer {
	return &groupServer{
		srv: s,
	}
}

func (s *groupServer) GroupGet(ctx context.Context, req *pb.GroupGetRequest) (*pb.GroupGetResponse, error) {	
	group, err := s.srv.GroupGet(ctx, req)
	return &pb.GroupGetResponse{Group: group}, err
}

func (s *groupServer) GroupList(ctx context.Context, req *pb.GroupListRequest) (*pb.GroupListResponse, error) {
	groups, err := s.srv.GroupList(ctx, req)
	return &pb.GroupListResponse{Groups: groups}, err
}
