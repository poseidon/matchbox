package server

import (
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/storage"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// Config configures an RPC Server.
type Config struct {
	Store storage.Store
}

// server implements the grpc GroupsServer interface.
type server struct {
	store storage.Store
}

// NewServer returns a new server.
func NewServer(config *Config) pb.GroupsServer {
	return &server{
		store: config.Store,
	}
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	group, err := s.store.GetGroup(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{Group: group}, nil
}

func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	groups, err := s.store.ListGroups()
	if err != nil {
		return nil, err
	}
	return &pb.ListResponse{Groups: groups}, nil
}