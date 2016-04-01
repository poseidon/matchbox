package server

import (
	"errors"
	"sort"

	"golang.org/x/net/context"

	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
	"github.com/coreos/coreos-baremetal/bootcfg/storage"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

var (
	errNoMatchingGroup = errors.New("bootcfg: No matching Group")
	errNoProfileFound  = errors.New("bootcfg: No Profile found")
)

// Server defines a bootcfg Server.
type Server interface {
	SelectGroup(ctx context.Context, req *pb.SelectGroupRequest) (*storagepb.Group, error)
	SelectProfile(ctx context.Context, req *pb.SelectProfileRequest) (*storagepb.Profile, error)
	pb.GroupsServer
	pb.ProfilesServer
	IgnitionGet(ctx context.Context, name string) (string, error)
	CloudGet(ctx context.Context, name string) (string, error)
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

func (s *server) ProfilePut(ctx context.Context, req *pb.ProfilePutRequest) (*pb.ProfilePutResponse, error) {
	if err := req.Profile.AssertValid(); err != nil {
		return nil, err
	}
	err := s.store.ProfilePut(req.Profile)
	if err != nil {
		return nil, err
	}
	return &pb.ProfilePutResponse{}, nil
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

// SelectGroup selects the Group whose selector matches the given labels.
// Groups are evaluated in sorted order from most selectors to least, using
// alphabetical order as a deterministic tie-breaker.
func (s *server) SelectGroup(ctx context.Context, req *pb.SelectGroupRequest) (*storagepb.Group, error) {
	groups, err := s.store.GroupList()
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(storagepb.ByReqs(groups)))
	for _, group := range groups {
		if group.Matches(req.Labels) {
			return group, nil
		}
	}
	return nil, errNoMatchingGroup
}

func (s *server) SelectProfile(ctx context.Context, req *pb.SelectProfileRequest) (*storagepb.Profile, error) {
	group, err := s.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: req.Labels})
	if err == nil {
		// lookup the Profile by id
		resp, err := s.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err == nil {
			return resp.Profile, nil
		}
		return nil, errNoProfileFound
	}
	return nil, errNoMatchingGroup
}

// IgnitionGet gets an Ignition Config template by name.
func (s *server) IgnitionGet(ctx context.Context, name string) (string, error) {
	return s.store.IgnitionGet(name)
}

// CloudGet gets a Cloud-Config template by name.
func (s *server) CloudGet(ctx context.Context, name string) (string, error) {
	return s.store.CloudGet(name)
}
