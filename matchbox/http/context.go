package http

import (
	"errors"

	"context"

	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

// unexported key prevents collisions
type key int

const (
	profileKey key = iota
	groupKey
)

var (
	errNoProfileFromContext = errors.New("api: Context missing a Profile")
	errNoGroupFromContext   = errors.New("api: Context missing a Group")
)

// withProfile returns a copy of ctx that stores the given Profile.
func withProfile(ctx context.Context, profile *storagepb.Profile) context.Context {
	return context.WithValue(ctx, profileKey, profile)
}

// profileFromContext returns the Profile from the ctx.
func profileFromContext(ctx context.Context) (*storagepb.Profile, error) {
	profile, ok := ctx.Value(profileKey).(*storagepb.Profile)
	if !ok {
		return nil, errNoProfileFromContext
	}
	return profile, nil
}

// withGroup returns a copy of ctx that stores the given Group.
func withGroup(ctx context.Context, group *storagepb.Group) context.Context {
	return context.WithValue(ctx, groupKey, group)
}

// groupFromContext returns the Group from the ctx.
func groupFromContext(ctx context.Context) (*storagepb.Group, error) {
	group, ok := ctx.Value(groupKey).(*storagepb.Group)
	if !ok {
		return nil, errNoGroupFromContext
	}
	return group, nil
}
