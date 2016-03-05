package server

import (
	"golang.org/x/net/context"

	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
	"github.com/coreos/coreos-baremetal/bootcfg/storage"
)

// Server defines a bootcfg Server.
type Server interface {
	pb.GroupsServer
	pb.ProfilesServer
}

// Config configures a server implementation.
type Config struct {
	Store storage.Store
}

// server implements the Server interface.
type server struct {
	store storage.Store
}

// NewServer returns a new Server.
func NewServer(config *Config) Server {
	return &server{
		store: config.Store,
	}
}

func (s *server) GroupGet(ctx context.Context, req *pb.GroupGetRequest) (*pb.GroupGetResponse, error) {
	group, err := s.store.GroupGet(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GroupGetResponse{Group: group}, nil
}

func (s *server) GroupList(ctx context.Context, req *pb.GroupListRequest) (*pb.GroupListResponse, error) {
	groups, err := s.store.GroupList()
	if err != nil {
		return nil, err
	}
	return &pb.GroupListResponse{Groups: groups}, nil
}

func (s *server) ProfileGet(ctx context.Context, req *pb.ProfileGetRequest) (*pb.ProfileGetResponse, error) {
	profile, err := s.store.ProfileGet(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProfileGetResponse{Profile: profile}, nil
}

func (s *server) ProfileList(ctx context.Context, req *pb.ProfileListRequest) (*pb.ProfileListResponse, error) {
	profiles, err := s.store.ProfileList()
	if err != nil {
		return nil, err
	}
	return &pb.ProfileListResponse{Profiles: profiles}, nil
}
