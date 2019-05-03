package server

import (
	"errors"
	"sort"

	"context"

	pb "github.com/poseidon/matchbox/matchbox/server/serverpb"
	"github.com/poseidon/matchbox/matchbox/storage"
	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

// Possible service errors
var (
	ErrNoMatchingGroup   = errors.New("matchbox: No matching Group")
	ErrNoMatchingProfile = errors.New("matchbox: No matching Profile")
)

// Server defines the matchbox server interface.
type Server interface {
	// SelectGroup returns the Group matching the given labels.
	SelectGroup(context.Context, *pb.SelectGroupRequest) (*storagepb.Group, error)
	// SelectProfile returns the Profile matching the given labels.
	SelectProfile(context.Context, *pb.SelectProfileRequest) (*storagepb.Profile, error)

	// Create or update a Group.
	GroupPut(context.Context, *pb.GroupPutRequest) (*storagepb.Group, error)
	// Get a machine Group by id.
	GroupGet(context.Context, *pb.GroupGetRequest) (*storagepb.Group, error)
	// Delete a machine Group by id.
	GroupDelete(context.Context, *pb.GroupDeleteRequest) error
	// List all machine Groups.
	GroupList(context.Context, *pb.GroupListRequest) ([]*storagepb.Group, error)

	// Create or update a Profile.
	ProfilePut(context.Context, *pb.ProfilePutRequest) (*storagepb.Profile, error)
	// Get a Profile by id.
	ProfileGet(context.Context, *pb.ProfileGetRequest) (*storagepb.Profile, error)
	// Delete a Profile by id.
	ProfileDelete(context.Context, *pb.ProfileDeleteRequest) error
	// List all Profiles.
	ProfileList(context.Context, *pb.ProfileListRequest) ([]*storagepb.Profile, error)

	// Create or update an Ignition template.
	IgnitionPut(context.Context, *pb.IgnitionPutRequest) (string, error)
	// Get an Ignition template by name.
	IgnitionGet(context.Context, *pb.IgnitionGetRequest) (string, error)
	// Delete an Ignition template by name.
	IgnitionDelete(context.Context, *pb.IgnitionDeleteRequest) error

	// Create or update an Generic template.
	GenericPut(context.Context, *pb.GenericPutRequest) (string, error)
	// Get an Generic template by name.
	GenericGet(context.Context, *pb.GenericGetRequest) (string, error)
	// Delete an Generic template by name.
	GenericDelete(context.Context, *pb.GenericDeleteRequest) error

	// Get a Cloud-Config template by name.
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

func (s *server) GroupPut(ctx context.Context, req *pb.GroupPutRequest) (*storagepb.Group, error) {
	if err := req.Group.AssertValid(); err != nil {
		return nil, err
	}
	err := s.store.GroupPut(req.Group)
	if err != nil {
		return nil, err
	}
	return req.Group, nil
}

func (s *server) GroupGet(ctx context.Context, req *pb.GroupGetRequest) (*storagepb.Group, error) {
	group, err := s.store.GroupGet(req.Id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *server) GroupDelete(ctx context.Context, req *pb.GroupDeleteRequest) error {
	return s.store.GroupDelete(req.Id)
}

func (s *server) GroupList(ctx context.Context, req *pb.GroupListRequest) ([]*storagepb.Group, error) {
	groups, err := s.store.GroupList()
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *server) ProfilePut(ctx context.Context, req *pb.ProfilePutRequest) (*storagepb.Profile, error) {
	if err := req.Profile.AssertValid(); err != nil {
		return nil, err
	}
	err := s.store.ProfilePut(req.Profile)
	if err != nil {
		return nil, err
	}
	return req.Profile, nil
}

func (s *server) ProfileGet(ctx context.Context, req *pb.ProfileGetRequest) (*storagepb.Profile, error) {
	profile, err := s.store.ProfileGet(req.Id)
	if err != nil {
		return nil, err
	}
	if err := profile.AssertValid(); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *server) ProfileDelete(ctx context.Context, req *pb.ProfileDeleteRequest) error {
	return s.store.ProfileDelete(req.Id)
}

func (s *server) ProfileList(ctx context.Context, req *pb.ProfileListRequest) ([]*storagepb.Profile, error) {
	profiles, err := s.store.ProfileList()
	if err != nil {
		return nil, err
	}
	return profiles, nil
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
	return nil, ErrNoMatchingGroup
}

func (s *server) SelectProfile(ctx context.Context, req *pb.SelectProfileRequest) (*storagepb.Profile, error) {
	group, err := s.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: req.Labels})
	if err == nil {
		// lookup the Profile by id
		profile, err := s.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err == nil {
			return profile, nil
		}
		return nil, ErrNoMatchingProfile
	}
	return nil, ErrNoMatchingGroup
}

// IgnitionPut creates or updates an Ignition template by name.
func (s *server) IgnitionPut(ctx context.Context, req *pb.IgnitionPutRequest) (string, error) {
	err := s.store.IgnitionPut(req.Name, req.Config)
	if err != nil {
		return "", err
	}
	return string(req.Config), err
}

// IgnitionGet gets an Ignition template by name.
func (s *server) IgnitionGet(ctx context.Context, req *pb.IgnitionGetRequest) (string, error) {
	return s.store.IgnitionGet(req.Name)
}

// IgnitionDelete deletes an Ignition template by name.
func (s *server) IgnitionDelete(ctx context.Context, req *pb.IgnitionDeleteRequest) error {
	return s.store.IgnitionDelete(req.Name)
}

// GenericPut creates or updates an Generic template by name.
func (s *server) GenericPut(ctx context.Context, req *pb.GenericPutRequest) (string, error) {
	err := s.store.GenericPut(req.Name, req.Config)
	if err != nil {
		return "", err
	}
	return string(req.Config), err
}

// GenericGet gets an Generic template by name.
func (s *server) GenericGet(ctx context.Context, req *pb.GenericGetRequest) (string, error) {
	return s.store.GenericGet(req.Name)
}

// GenericDelete deletes an Generic template by name.
func (s *server) GenericDelete(ctx context.Context, req *pb.GenericDeleteRequest) error {
	return s.store.GenericDelete(req.Name)
}

// CloudGet gets a Cloud-Config template by name.
func (s *server) CloudGet(ctx context.Context, name string) (string, error) {
	return s.store.CloudGet(name)
}
