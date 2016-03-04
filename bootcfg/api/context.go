package api

import (
	"errors"

	"golang.org/x/net/context"
)

// unexported key prevents collisions
type key int

const (
	specKey key = iota
	groupKey
)

var (
	errNoSpecFromContext  = errors.New("api: Context missing a Spec")
	errNoGroupFromContext = errors.New("api: Context missing a Group")
)

// withSpec returns a copy of ctx that stores the given Spec.
func withSpec(ctx context.Context, spec *Spec) context.Context {
	return context.WithValue(ctx, specKey, spec)
}

// specFromContext returns the Spec from the ctx.
func specFromContext(ctx context.Context) (*Spec, error) {
	spec, ok := ctx.Value(specKey).(*Spec)
	if !ok {
		return nil, errNoSpecFromContext
	}
	return spec, nil
}

// withGroup returns a copy of ctx that stores the given Group.
func withGroup(ctx context.Context, group *Group) context.Context {
	return context.WithValue(ctx, groupKey, group)
}

// groupFromContext returns the Group from the ctx.
func groupFromContext(ctx context.Context) (*Group, error) {
	group, ok := ctx.Value(groupKey).(*Group)
	if !ok {
		return nil, errNoGroupFromContext
	}
	return group, nil
}
